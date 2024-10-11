package webview

import (
	"path/filepath"

	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
	"github.com/hashicorp/go-multierror"
)

const (
	jsPkg      = "syscall/js"
	contextPkg = "context"
	errorPkg   = "errors"
)

type resultValPath struct {
	path       []string
	typeMethod string
}

type Plugin struct {
	ctx *gg.Context
}

func (p *Plugin) Name() string { return "webview" }

func (p *Plugin) Exec() (files []file.File, errs error) {
	serverFile := file.NewGoFile(p.ctx.Module, p.ServerOutput())
	clientFile := file.NewGoFile(p.ctx.Module, p.ClientOutput())

	serverFile.Type().Id("Context").Map(jen.String()).Any()

	serverFile.Func().Params(jen.Id("c").Id("Context")).Id("Get").Params(jen.Id("key").String()).Any().Block(
		jen.List(jen.Id("v"), jen.Id("_")).Op(":=").Id("c").Index(jen.Id("key")),
		jen.Return(jen.Id("v")),
	)

	for _, iface := range p.ctx.Interfaces {
		optionsName := iface.Named.Name + "Options"
		optionName := iface.Named.Name + "Option"
		serverFile.Type().Id(optionName).Func().Params(jen.Op("*").Id(optionsName))
		serverFile.Type().Id(optionsName).StructFunc(func(g *jen.Group) {})
	}

	serverFile.Type().Id("binder").Interface(
		jen.Id("Bind").Params(jen.String(), jen.Any()).Error(),
	)

	serverFile.Type().Id("bindCallbackResult").Struct(
		jen.Id("value").Any(),
		jen.Err().Error(),
	)

	serverFile.Func().Params(jen.Id("r").Op("*").Id("bindCallbackResult")).Id("Error").Params().Error().Block(
		jen.Return(jen.Id("r").Dot("err")),
	)

	serverFile.Func().Params(jen.Id("r").Op("*").Id("bindCallbackResult")).Id("Value").Params().Any().Block(
		jen.Return(jen.Id("r").Dot("value")),
	)

	var makeParamValue func(parentName string, t any) jen.Code
	makeParamValue = func(parentName string, t any) jen.Code {
		switch t := t.(type) {
		default:
			return jen.Id(parentName)
		case *types.Named:
			if st := t.Struct(); st != nil {
				return jen.Map(jen.String()).Interface().ValuesFunc(func(g *jen.Group) {
					for _, f := range st.Fields {
						g.Lit(f.Var.Name).Op(":").Add(makeParamValue(parentName+"."+f.Var.Name, f.Var.Type))
					}
				})
			}
		}
		return jen.Nil()
	}

	var makePathRecursive func(v *types.Var, path []string, cb func([]string, string))
	makePathRecursive = func(v *types.Var, path []string, cb func([]string, string)) {
		if named, ok := v.Type.(*types.Named); ok {
			if st := named.Struct(); st != nil {
				for _, f := range st.Fields {
					makePathRecursive(f.Var, append(path, f.Var.Name), cb)
				}
			}
		} else if t, ok := v.Type.(*types.Basic); ok {
			var typeMethod string
			switch {
			case t.IsBool():
				typeMethod = "Bool"
			case t.IsInteger():
				typeMethod = "Int"
			case t.IsFloat():
				typeMethod = "Float"
			case t.IsString():
				typeMethod = "String"
			}
			cb(path, typeMethod)
		}
	}

	makePath := func(v *types.Var) (paths []resultValPath) {
		makePathRecursive(v, nil, func(path []string, typeMethod string) {
			paths = append(paths, resultValPath{
				path:       path,
				typeMethod: typeMethod,
			})
		})
		return
	}

	for _, iface := range p.ctx.Interfaces {
		clientStructName := strcase.ToCamel(iface.Named.Name) + "Client"

		clientFile.Type().Id(clientStructName).Struct()

		for _, method := range iface.Type.Methods {
			bindName := makeBindName(iface, method)
			promiseName := strcase.ToCamel(iface.Named.Name) + method.Name + "Promise"

			clientFile.Type().Id(promiseName).StructFunc(func(g *jen.Group) {
				for _, p := range method.Sig.Params {
					if p.IsContext {
						continue
					}
					g.Add(types.Convert(p, clientFile.Import))
				}
			})

			clientFile.Func().Params(
				jen.Id("p").Op("*").Id(promiseName),
			).Id("Execute").ParamsFunc(func(g *jen.Group) {
				g.Id("ctx").Qual(contextPkg, "Context")
				g.Id("cb").Func().ParamsFunc(func(g *jen.Group) {
					for _, r := range method.Sig.Results {
						g.Add(types.Convert(r, clientFile.Import))
					}
				})
			}).BlockFunc(func(g *jen.Group) {
				g.Id("wvCtx").Op(":=").Map(jen.String()).Any().ValuesFunc(func(g *jen.Group) {})
				for _, p := range method.Sig.Params {
					if p.IsContext {
						continue
					}
					g.Id(p.Name + "Req").Op(":=").Add(makeParamValue("p."+p.Name, p.Type))
				}
				g.Id("v").Op(":=").Do(clientFile.Import(jsPkg, "Global")).Call().Dot("Call").CallFunc(func(g *jen.Group) {
					g.Lit(bindName)
					for _, p := range method.Sig.Params {
						if p.IsContext {
							g.Id("wvCtx")
							continue
						}
						g.Id(p.Name + "Req")
					}
				})
				g.Id("thenFunc").Op(":=").Qual(jsPkg, "FuncOf").Call(
					jen.Func().Params(
						jen.Id("this").Qual(jsPkg, "Value"),
						jen.Id("args").Index().Qual(jsPkg, "Value"),
					).Any().BlockFunc(func(g *jen.Group) {
						for _, r := range method.Sig.Results {
							if r.IsError {
								continue
							}
							g.Var().Add(types.Convert(r, clientFile.Import))
						}
						for _, r := range method.Sig.Results {
							if r.IsError {
								continue
							}
							for _, p := range makePath(r) {
								g.Id(r.Name).Do(func(s *jen.Statement) {
									for _, fldName := range p.path {
										s.Dot(fldName)
									}
								}).Op("=").Id("args").Index(jen.Lit(0)).Dot("Get").Call(jen.Lit(r.Name)).Do(func(s *jen.Statement) {
									for _, fldName := range p.path {
										s.Dot("Get").Call(jen.Lit(fldName))
									}
									s.Dot(p.typeMethod).Call()
								})
							}
						}
						g.Id("cb").CallFunc(func(g *jen.Group) {
							for _, r := range method.Sig.Results {
								if r.IsError {
									g.Nil()
									continue
								}
								g.Id(r.Name)
							}
						})
						g.Return(jen.Nil())
					}),
				)
				// g.Defer().Id("thenFunc").Dot("Release").Call()

				g.Id("catchFunc").Op(":=").Qual(jsPkg, "FuncOf").Call(
					jen.Func().Params(
						jen.Id("this").Qual(jsPkg, "Value"),
						jen.Id("args").Index().Qual(jsPkg, "Value"),
					).Any().BlockFunc(func(g *jen.Group) {
						for _, r := range method.Sig.Results {
							if r.IsError {
							}
						}
						g.Id("cb").CallFunc(func(g *jen.Group) {
							for _, r := range method.Sig.Results {
								if r.IsError {
									g.Qual(errorPkg, "New").Call(jen.Id("args").Index(jen.Lit(0)).Dot("String").Call())
									continue
								}
								g.Add(gg.ZeroValue(r.Type, clientFile.Import))
							}
						})
						g.Return(jen.Nil())
					}),
				)
				// g.Defer().Id("catchFunc").Dot("Release").Call()

				g.Id("v").
					Dot("Call").Call(jen.Lit("then"), jen.Id("thenFunc")).
					Dot("Call").Call(jen.Lit("catch"), jen.Id("catchFunc"))
			})

			clientFile.Func().
				Params(
					jen.Id("m").Op("*").Id(clientStructName),
				).
				Id(method.Name).ParamsFunc(func(g *jen.Group) {
				for _, p := range method.Sig.Params {
					if p.IsContext {
						continue
					}
					g.Add(types.Convert(p, clientFile.Import))
				}
			}).Op("*").Id(promiseName).
				BlockFunc(func(g *jen.Group) {

					// // v.Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
					// // 	if len(args)>0 {
					// // 		fmt.Println(len(args))
					// // 		fmt.Println(args[0].Get("testVar").String())
					// // 	}
					// // 	return nil
					// // }))

					// g.Id("waitCh").Op(":=").Make(jen.Chan().Struct())

					// g.Id("v").Dot("Call").Call(jen.Lit("then"), jen.Qual(jsPkg, "FuncOf").Call(
					// 	jen.Func().Params(
					// 		jen.Id("this").Qual(jsPkg, "Value"),
					// 		jen.Id("args").Index().Qual(jsPkg, "Value"),
					// 	).Any().BlockFunc(func(g *jen.Group) {
					// 		g.Defer().Close(jen.Id("waitCh"))
					// 		for _, r := range method.Sig.Results {
					// 			if r.IsError {
					// 				continue
					// 			}
					// 			for _, p := range makePath(r) {
					// 				g.Id(r.Name).Do(func(s *jen.Statement) {
					// 					for _, fldName := range p.path {
					// 						s.Dot(fldName)
					// 					}
					// 				}).Op("=").Id("args").Index(jen.Lit(0)).Dot("Get").Call(jen.Lit(r.Name)).Do(func(s *jen.Statement) {
					// 					for _, fldName := range p.path {
					// 						s.Dot("Get").Call(jen.Lit(fldName))
					// 					}
					// 					s.Dot(p.typeMethod).Call()
					// 				})
					// 			}
					// 		}
					// 		g.Return(jen.Nil())
					// 	}),
					// ))

					// g.Op("<-").Id("waitCh")

					g.Return(
						jen.Op("&").Id(promiseName).ValuesFunc(func(g *jen.Group) {
							for _, p := range method.Sig.Params {
								if p.IsContext {
									continue
								}
								g.Id(p.Name).Op(":").Id(p.Name)
							}
						}),
					)
				})
		}
	}

	for _, iface := range p.ctx.Interfaces {
		serverFile.Func().Id("SetupRoutes"+iface.Named.Name).Params(
			jen.Id("svc").Do(serverFile.Qual(iface.Named.Pkg.Path, iface.Named.Name)),
			jen.Id("w").Id("binder"),
			jen.Id("opts").Op("...").Id(iface.Named.Name+"Option"),
		).BlockFunc(func(g *jen.Group) {
			for _, m := range iface.Type.Methods {
				bindName := makeBindName(iface, m)
				g.Id("w").Dot("Bind").Call(
					jen.Lit(bindName), jen.Func().ParamsFunc(func(g *jen.Group) {
						g.Id("wvCtx").Id("Context")
						for _, p := range m.Sig.Params {
							if p.IsContext {
								continue
							}
							g.Id(p.Name).Add(types.Convert(p.Type, serverFile.Qual))
						}
					}).Params(jen.Chan().Op("*").Id("bindCallbackResult")).BlockFunc(func(g *jen.Group) {
						g.Id("ch").Op(":=").Make(jen.Chan().Op("*").Id("bindCallbackResult"))

						g.Go().Func().Params().BlockFunc(func(g *jen.Group) {
							g.Id("result").Op(":=").StructFunc(func(g *jen.Group) {
								for _, r := range m.Sig.Results {
									if r.IsError {
										continue
									}
									g.Id(strcase.ToCamel(r.Name)).Add(types.Convert(r.Type, serverFile.Qual)).Tag(map[string]string{"json": strcase.ToLowerCamel(r.Name)})
								}
							}).Values()

							g.Id("ctx").Op(":=").Qual("context", "TODO").Call()

							tags := m.Tags.GetSlice("webview-context")
							for _, t := range tags {
								if t.Value == "" {
									errs = multierror.Append(errs, errors.Error("the path to the context key is required", t.Position))
									return
								}
								pkgPath, name, err := p.ctx.Module.ParseImportPath(t.Value)
								if err != nil {
									errs = multierror.Append(errs, err)
									return
								}
								g.Id("ctx").Op("=").Qual("context", "WithValue").Call(
									jen.Id("ctx"),
									jen.Qual(pkgPath, name),
									jen.Id("wvCtx").Dot("Get").Call(jen.Lit(name)),
								)
							}
							g.Var().Err().Error()
							g.ListFunc(func(g *jen.Group) {
								for _, r := range m.Sig.Results {
									if r.IsError {
										g.Err()
										continue
									}
									g.Id("result").Dot(strcase.ToCamel(r.Name))
								}
							}).Op("=").Id("svc").Dot(m.Name).CallFunc(func(g *jen.Group) {
								for _, p := range m.Sig.Params {
									if p.IsContext {
										g.Id("ctx")
										continue
									}
									g.Id(p.Name)
								}
							})
							g.Id("ch").Op("<-").Op("&").Id("bindCallbackResult").Values(
								jen.Id("value").Op(":").Id("result"),
								jen.Id("err").Op(":").Id("err"),
							)
						}).Call()
						g.Return(jen.Id("ch"))
					}),
				)
			}
		})
	}

	return []file.File{serverFile, clientFile}, nil
}

func (p *Plugin) ServerOutput() string {
	return filepath.Join(p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("output", "internal/handler/handler.go"))
}

func (p *Plugin) ClientOutput() string {
	return filepath.Join(p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("output", "internal/client/client.go"))
}

func (p *Plugin) Dependencies() []string { return nil }
