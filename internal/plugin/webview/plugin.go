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

type FlattenPath struct {
	Name       string
	Path       *jen.Statement
	AssignPath *jen.Statement
	Paths      []PathName
	Children   []FlattenPath
	IsArray    bool
	Method     string
	Type       any
}

type PathName struct {
	name     string
	jsonName string
}

func (n PathName) Value() string {
	return n.name
}

func (n PathName) JSON() string {
	if n.jsonName == "" {
		return n.name
	}
	return n.jsonName
}

type FlattenProcessor struct {
	allPaths    []FlattenPath
	currentPath []PathName
}

func (p *FlattenProcessor) varsByType(t any) (vars types.Vars) {
	switch t := t.(type) {
	case *types.Struct:
		for _, f := range t.Fields {
			vars = append(vars, f)
		}
		return vars
	case *types.Named:
		return p.varsByType(t.Type)
	}
	return nil
}

func (p *FlattenProcessor) flattenVar(v *types.Var) {
	var jsonName string
	if tag, err := v.SysTags.Get("json"); err == nil {
		jsonName = tag.Name
	}
	p.currentPath = append(p.currentPath, PathName{name: v.Name, jsonName: jsonName})
	vars := p.varsByType(v.Type)
	if len(vars) == 0 {

		var (
			children []FlattenPath
			isArray  bool
			method   string
		)

		if s, ok := v.Type.(*types.Slice); ok {
			isArray = true
			for _, v := range p.varsByType(s.Value) {
				children = append(children, new(FlattenProcessor).Flatten(v)...)
			}
			method = methodByType(s.Value)
		} else {
			method = methodByType(v.Type)
		}

		currentPath := make([]PathName, len(p.currentPath))
		copy(currentPath, p.currentPath)

		assignPath := jen.Id("Get").Call(jen.Lit(currentPath[0].JSON())).Do(func(s *jen.Statement) {
			for i := 1; i < len(currentPath); i++ {
				s.Dot("Get").Call(jen.Lit(currentPath[i].JSON()))
			}
		})

		path := jen.Id(currentPath[0].Value()).Do(func(s *jen.Statement) {
			for i := 1; i < len(currentPath); i++ {
				s.Dot(currentPath[i].Value())
			}
		})

		p.allPaths = append(p.allPaths, FlattenPath{
			Name:       v.Name,
			Paths:      currentPath,
			Children:   children,
			Path:       path,
			AssignPath: assignPath,
			Method:     method,
			IsArray:    isArray,
			Type:       v.Type,
		})
	} else {
		for _, v := range vars {
			p.flattenVar(v)
		}
	}

	p.currentPath = p.currentPath[:len(p.currentPath)-1]
}

func (p *FlattenProcessor) Flatten(v *types.Var) []FlattenPath {
	p.flattenVar(v)
	return p.allPaths
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
						g.Lit(f.Name).Op(":").Add(makeParamValue(parentName+"."+f.Name, f.Type))
					}
				})
			}
		}
		return jen.Nil()
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

				var assignFunc func(nestedIndex int, g *jen.Group, p FlattenPath, parentPath, parentAssignPath *jen.Statement)
				assignFunc = func(nestedIndex int, g *jen.Group, p FlattenPath, parentPath *jen.Statement, parentAssignPath *jen.Statement) {

					if p.IsArray {
						indexName := makeIndexName(nestedIndex)

						g.If(jen.Op("!").Id("args").Index(jen.Lit(0)).Op(".").Add(parentAssignPath.Clone().Dot("IsNull").Call())).BlockFunc(func(g *jen.Group) {

							g.Add(parentPath).Op("=").Make(types.Convert(p.Type, clientFile.Import), jen.Id("args").Index(jen.Lit(0)).Op(".").Add(parentAssignPath.Clone().Dot("Length").Call()))

							g.For(
								jen.Id(indexName).Op(":=").Lit(0),
								jen.Id(indexName).Op("<").Len(parentPath),
								jen.Id(indexName).Op("++"),
							).BlockFunc(func(g *jen.Group) {
								for _, child := range p.Children {
									path := jen.Add(parentPath).Index(jen.Id(indexName)).Op(".").Add(child.Path)
									assignPath := jen.Add(parentAssignPath).Op(".").Id("Index").Call(jen.Id(indexName)).Op(".").Add(child.AssignPath)
									if child.IsArray {
										assignFunc(nestedIndex+1, g, child, path, assignPath)
									} else {
										g.Add(path).Op("=").Id("args").Index(jen.Lit(0)).Op(".").Add(assignPath).Dot(child.Method).Call()
									}
								}
							})
						})
					} else {
						g.Add(p.Path).Op("=").Add(p.AssignPath).Dot(p.Method).Call()
					}

				}

				g.Id("thenFunc").Op(":=").Qual(jsPkg, "FuncOf").Call(
					jen.Func().Params(
						jen.Id("this").Qual(jsPkg, "Value"),
						jen.Id("args").Index().Qual(jsPkg, "Value"),
					).Any().BlockFunc(func(g *jen.Group) {
						for _, r := range method.Sig.Results {
							if r.IsError {
								continue
							}
							g.Add(types.NewConstruct(clientFile.Import).Convert(r))
						}

						for _, r := range method.Sig.Results {
							if r.IsError {
								continue
							}
							allPaths := new(FlattenProcessor).Flatten(r)
							for _, p := range allPaths {
								assignFunc(0, g, p, p.Path, p.AssignPath)
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
		serverFile.Func().Id("Setup"+iface.Named.Name).Params(
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
	return filepath.Join(p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("output", "internal/transport/handler.go"))
}

func (p *Plugin) ClientOutput() string {
	return filepath.Join(p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("output", "internal/client/client.go"))
}

func (p *Plugin) Dependencies() []string { return nil }
