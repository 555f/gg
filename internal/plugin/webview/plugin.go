package webview

import (
	"path/filepath"

	"github.com/555f/gg/internal/plugin/webview/handlermux"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gen"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

const (
	jsPkg      = "syscall/js"
	contextPkg = "context"
	errorPkg   = "errors"
	syncPkg    = "sync"
	jsonPkg    = "encoding/json"
	base64Pkg  = "encoding/base64"
)

type Plugin struct {
	ctx *gg.Context
}

func (p *Plugin) Name() string { return "webview" }

func (p *Plugin) Exec() (files []file.File, errs error) {
	serverFile := file.NewGoFile(p.ctx.Module, p.ServerOutput())
	clientFile := file.NewGoFile(p.ctx.Module, p.ClientOutput())

	var (
		jsClientInterfaces  gg.Interfaces
		asmClientInterfaces gg.Interfaces
	)

	for _, iface := range p.ctx.Interfaces {
		for _, tag := range iface.Named.Tags.GetSlice("webview-client") {
			switch tag.Value {
			case "asm":
				asmClientInterfaces = append(asmClientInterfaces, iface)
			case "js":
				jsClientInterfaces = append(jsClientInterfaces, iface)
				break
			}
		}
	}

	if len(jsClientInterfaces) > 0 {
		jsClientFile := file.NewTxtFile(filepath.Join(p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("output", "internal/transport/client.js")))

		jsClientFile.WriteText("/*global executeHandler*/\n\n")

		for _, iface := range jsClientInterfaces {
			className := iface.Named.Name
			if tag, ok := iface.Named.Tags.Get("webview-name"); ok {
				className = tag.Value
			}
			jsClientFile.WriteText("export class %s {}\n\n", className)
			for _, m := range iface.Named.Interface().Methods {
				jsClientFile.WriteText("%s.%s = async (meta", className, strcase.ToLowerCamel(m.Name))
				for _, p := range m.Sig.Params {
					if p.IsContext {
						continue
					}
					jsClientFile.WriteText(", %s", p.Name)
				}
				jsClientFile.WriteText(") => {\n")

				jsClientFile.WriteText("  const payload = {\n")
				jsClientFile.WriteText("    Meta: meta,\n")
				jsClientFile.WriteText("    Method: \"%s.%s\",\n", iface.Named.Name, m.Name)
				jsClientFile.WriteText("    Params: {\n")
				for _, p := range m.Sig.Params {
					if p.IsContext {
						continue
					}
					jsClientFile.WriteText("      %[1]s: %[1]s,\n", strcase.ToLowerCamel(p.Name))
				}
				jsClientFile.WriteText("    }\n")
				jsClientFile.WriteText("  }\n\n")
				jsClientFile.WriteText("  const result = await executeHandler(JSON.stringify(payload));\n")
				jsClientFile.WriteText("  return JSON.parse(atob(result));\n")
				jsClientFile.WriteText("}\n\n")
			}
		}

		files = append(files, jsClientFile)
	}

	serverFile.Add(handlermux.New().Generate())

	serverFile.Func().Id("RegisterHandlers").ParamsFunc(func(g *jen.Group) {
		for _, iface := range p.ctx.Interfaces {
			g.Id(strcase.ToLowerCamel(iface.Named.Name)).Do(serverFile.Import(iface.Named.Pkg.Path, iface.Named.Name))
		}
	}).BlockFunc(func(g *jen.Group) {
		for _, iface := range p.ctx.Interfaces {
			for _, m := range iface.Type.Methods {
				g.Id("Register").Call(
					jen.Lit(iface.Named.Name+"."+m.Name),
					jen.Func().Params(
						jen.Id("ctx").Qual("context", "Context"),
						jen.Id("r").Op("*").Id("Request"),
					).Params(jen.String(), jen.Error()).BlockFunc(func(g *jen.Group) {
						g.Var().Id("reqParams").StructFunc(func(g *jen.Group) {
							for _, p := range m.Sig.Params {
								if p.IsContext {
									continue
								}
								tags := map[string]string{
									"json": strcase.ToLowerCamel(p.Name),
								}
								g.Id(strcase.ToCamel(p.Name)).Add(types.Convert(p.Type, serverFile.Import)).Tag(tags)
							}
						})
						g.Err().Op(":=").Qual(jsonPkg, "Unmarshal").Call(jen.Index().Byte().Call(jen.Id("r").Dot("Params")), jen.Op("&").Id("reqParams"))
						g.Do(gen.CheckErr(
							jen.Return(jen.Lit(""), jen.Err()),
						))

						resultLen := m.Sig.Results.Len()
						if m.Sig.Results.HasError() {
							resultLen = resultLen - 1
						}

						if resultLen > 0 {
							g.Var().Id("resp").StructFunc(func(g *jen.Group) {
								for _, r := range m.Sig.Results {
									if r.IsError {
										continue
									}
									tags := map[string]string{
										"json": strcase.ToLowerCamel(r.Name),
									}
									g.Id(strcase.ToCamel(r.Name)).Add(types.Convert(r.Type, serverFile.Import)).Tag(tags)
								}
							})
						}
						g.Do(func(s *jen.Statement) {
							if m.Sig.Results.Len() > 0 {
								s.ListFunc(func(g *jen.Group) {
									for _, r := range m.Sig.Results {
										if r.IsError {
											g.Err()
											continue
										}
										g.Id("resp").Dot(strcase.ToCamel(r.Name))
									}
								})
								s.Op("=")
							}
							s.Id(strcase.ToLowerCamel(iface.Named.Name)).Dot(m.Name).CallFunc(func(g *jen.Group) {
								for _, p := range m.Sig.Params {
									if p.IsContext {
										g.Id("ctx")
										continue
									}
									g.Id("reqParams").Dot(strcase.ToCamel(p.Name))

								}
							})
						})
						if m.Sig.Results.HasError() {
							g.Do(gen.CheckErr(
								jen.Return(jen.Lit(""), jen.Err()),
							))
						}

						if resultLen > 0 {
							g.List(jen.Id("data"), jen.Err()).Op(":=").Qual(jsonPkg, "Marshal").Call(jen.Id("resp"))
							g.Do(gen.CheckErr(
								jen.Return(jen.Lit(""), jen.Err()),
							))
							g.Return(jen.Qual(base64Pkg, "StdEncoding").Dot("EncodeToString").Call(jen.Id("data")), jen.Nil())
						} else {
							g.Return(jen.Lit(""), jen.Nil())
						}
					}),
				)
			}
		}
	})

	files = append(files, serverFile)

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

	if len(asmClientInterfaces) > 0 {
		for _, iface := range asmClientInterfaces {
			clientStructName := strcase.ToCamel(iface.Named.Name) + "Client"

			clientFile.Type().Id(clientStructName).Struct()

			for _, method := range iface.Type.Methods {
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
					g.Var().Id("params").StructFunc(func(g *jen.Group) {
						for _, p := range method.Sig.Params {
							if p.IsContext {
								continue
							}
							g.Id(strcase.ToCamel(p.Name)).Add(types.Convert(p.Type, clientFile.Import))
						}
					})

					for _, p := range method.Sig.Params {
						if p.IsContext {
							continue
						}
						g.Id("params").Dot(strcase.ToCamel(p.Name)).Op("=").Id("p").Dot(p.Name)
					}

					g.List(jen.Id("paramsBytes"), jen.Id("_")).Op(":=").Qual(jsonPkg, "Marshal").Call(jen.Id("params"))

					g.Var().Id("req").Struct(
						jen.Id("Method").String(),
						jen.Id("Params").Qual(jsonPkg, "RawMessage"),
					)

					g.Id("req").Dot("Method").Op("=").Lit(iface.Named.Name + "." + method.Name)
					g.Id("req").Dot("Params").Op("=").Id("paramsBytes")

					g.List(jen.Id("data"), jen.Id("_")).Op(":=").Qual(jsonPkg, "Marshal").Call(jen.Id("req"))

					g.Id("v").Op(":=").Do(clientFile.Import(jsPkg, "Global")).Call().Dot("Call").CallFunc(func(g *jen.Group) {
						g.Lit("execute")
						g.Qual(jsPkg, "ValueOf").Call(jen.String().Call(jen.Id("data")))
					})
					jen.List(jen.Id("data"), jen.Id("_")).Op(":=").Qual(jsonPkg, "Marshal").Call(jen.Id("req"))

					g.Id("thenFunc").Op(":=").Qual(jsPkg, "FuncOf").Call(
						jen.Func().Params(
							jen.Id("this").Qual(jsPkg, "Value"),
							jen.Id("args").Index().Qual(jsPkg, "Value"),
						).Any().BlockFunc(func(g *jen.Group) {
							g.Var().Id("resp").StructFunc(func(g *jen.Group) {
								for _, r := range method.Sig.Results {
									if r.IsError {
										continue
									}
									g.Id(strcase.ToCamel(r.Name)).Add(types.Convert(r.Type, clientFile.Import))
								}
							})
							g.List(jen.Id("data"), jen.Id("_")).Op(":=").Qual(base64Pkg, "StdEncoding").Dot("DecodeString").Call(
								jen.Id("args").Index(jen.Lit(0)).Dot("String").Call(),
							)
							g.Qual(jsonPkg, "Unmarshal").Call(
								jen.Id("data"),
								jen.Op("&").Id("resp"),
							)

							g.Id("cb").CallFunc(func(g *jen.Group) {
								for _, r := range method.Sig.Results {
									if r.IsError {
										g.Nil()
										continue
									}
									g.Id("resp").Dot(strcase.ToCamel(r.Name))
								}
							})
							g.Return(jen.Nil())
						}),
					)

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
		files = append(files, clientFile)
	}

	return files, nil
}

func (p *Plugin) ServerOutput() string {
	return filepath.Join(p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("output", "internal/transport/handler.go"))
}

func (p *Plugin) ClientOutput() string {
	return filepath.Join(p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("output", "internal/client/client.go"))
}

func (p *Plugin) Dependencies() []string { return nil }
