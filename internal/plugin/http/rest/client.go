package rest

import (
	"strings"

	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gen"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"

	. "github.com/dave/jennifer/jen"
)

func GenClient(s options.Iface, errorWrapper *options.ErrorWrapper) func(f *file.GoFile) {
	return func(f *file.GoFile) {
		clientName := s.Name + "Client"

		f.Type().Id(clientName).StructFunc(func(g *Group) {
			g.Id("client").Op("*").Qual("net/http", "Client")
			g.Id("target").String()

		})
		for _, endpoint := range s.Endpoints {
			methodRequestName := s.Name + endpoint.MethodName + "Request"
			recvName := strcase.ToLowerCamel(endpoint.MethodName)

			f.Type().Id(methodRequestName).StructFunc(func(g *Group) {
				g.Id("c").Op("*").Id(clientName)
				g.Id("client").Op("*").Qual("net/http", "Client")
				g.Id("methodOpts").Op("*").Id("clientMethodOptions")
				g.Id("params").StructFunc(func(g *Group) {
					for _, param := range endpoint.Params {
						if len(param.Params) > 0 {
							for _, childParam := range param.Params {
								g.Add(makeRequestStructParam(param, childParam, f.Import))
							}
							continue
						}
						g.Add(makeRequestStructParam(nil, param, f.Import))
					}
				})
			})

			for _, param := range endpoint.Params {
				if len(param.Params) > 0 {
					for _, childParam := range param.Params {
						f.Add(makeSetFunc(recvName, methodRequestName, param, childParam, f.Import))
					}
				} else {
					if !param.Required {
						f.Add(makeSetFunc(recvName, methodRequestName, nil, param, f.Import))
					}
				}
			}

			f.Func().Params(Id(recvName).Op("*").Id(methodRequestName)).Id("Execute").Params(Id("opts").Op("...").Id("ClientMethodOption")).
				ParamsFunc(func(g *Group) {
					for _, result := range endpoint.Sig.Results {
						g.Id(result.Name).Add(types.Convert(result.Type, f.Import))
					}
				}).
				Block(
					For(List(Id("_"), Id("o")).Op(":=").Range().Id("opts")).Block(
						Id("o").Call(Id(recvName).Dot("methodOpts")),
					),
					Do(func(s *Statement) {
						if len(endpoint.BodyParams) > 0 {
							s.Var().Id("body").StructFunc(func(g *Group) {
								for _, param := range endpoint.BodyParams {
									jsonTag := param.Name
									fld := g.Id(param.FldName)
									if !param.Required {
										jsonTag += ",omitempty"
										fld.Op("*")
									}
									fld.Add(types.Convert(param.Type, f.Import)).Tag(map[string]string{"json": jsonTag})
								}
							})
						}
					}),
					List(Id("ctx"), Id("cancel")).Op(":=").Qual("context", "WithCancel").Call(Id(recvName).Dot("methodOpts").Dot("ctx")),
					Id("path").Op(":=").Qual("fmt", "Sprintf").Do(func(s *Statement) {
						var paramsCall []Code

						parts := strings.Split(endpoint.Path, "/")
						pathParamsMap := make(map[string]*options.EndpointParam, len(endpoint.PathParams))
						for _, param := range endpoint.PathParams {
							pathParamsMap[param.Name] = param
						}

						for i, part := range parts {
							startIndex := strings.Index(part, "{")
							endIndex := strings.Index(part, "}")

							if startIndex != -1 && endIndex != -1 {
								paramName := part[startIndex+1 : endIndex]
								if param, ok := pathParamsMap[paramName]; ok {
									parts[i] = "%s"
									if tp, ok := param.Type.(*types.Basic); ok {
										if tp.IsSigned() || tp.IsUnsigned() {
											parts[i] = "%d"
										} else if tp.IsFloat() {
											parts[i] = "%f"
										}
									}
								}
							}
						}
						paramsCall = append(paramsCall, Lit(strings.Join(parts, "/")))
						for _, name := range endpoint.ParamsNameIdx {
							paramsCall = append(paramsCall, Id(recvName).Dot("params").Dot(name))
						}
						s.Call(paramsCall...)
					}),
					List(Id("req"), Err()).Op(":=").Qual("net/http", "NewRequest").
						Call(Lit(endpoint.HTTPMethod), Id(recvName).Dot("c").Dot("target").Op("+").Id("path"), Nil()),
					Do(gen.CheckErr(
						Id("cancel").Call(),
						Return(),
					)),
					CustomFunc(Options{Multi: true}, func(g *Group) {
						if len(endpoint.BodyParams) > 0 {
							for _, param := range endpoint.BodyParams {
								g.Id("body").Dot(param.FldName).Op("=").Id(recvName).Dot("params").Dot(param.FldNameUnExport)
							}

							g.Var().Id("reqData").Qual("bytes", "Buffer")
							g.Err().Op("=").Qual("encoding/json", "NewEncoder").Call(Op("&").Id("reqData")).Dot("Encode").Call(Id("body"))
							g.If(Err().Op("!=").Nil()).Block(
								Id("cancel").Call(),
								Return(),
							)
							g.Id("req").Dot("Body").Op("=").Qual("io", "NopCloser").Call(Op("&").Id("reqData"))
						}
					}),

					CustomFunc(Options{Multi: true}, func(g *Group) {
						if len(endpoint.QueryParams) > 0 || len(endpoint.QueryValues) > 0 {
							g.Id("q").Op(":=").Id("req").Dot("URL").Dot("Query").Call()

							for _, param := range endpoint.QueryParams {
								g.Add(makeAddQueryParam(recvName, param.Parent, param, f.Import, endpoint.TimeFormat))
							}
							for _, param := range endpoint.QueryValues {
								g.Id("q").Dot("Add").Call(Lit(param.Name), Lit(param.Value))
							}
							g.Id("req").Dot("URL").Dot("RawQuery").Op("=").Id("q").Dot("Encode").Call()
						}

						g.Id("req").Dot("Header").Dot("Add").Call(Lit("Content-Type"), Lit("application/json"))

						for _, param := range endpoint.HeaderParams {
							g.Id("req").Dot("Header").Dot("Add").Call(Lit(param.Name), gen.FormatValue(Id(recvName).Dot("params").Dot(param.FldNameUnExport), param.Type, f.Import, endpoint.TimeFormat))
						}
						for _, param := range endpoint.CookieParams {
							g.Id("req").Dot("AddCookie").Call(Op("&").Qual("net/http", "Cookie").Values(
								Id("Name").Op(":").Lit(param.Name),
								Id("Value").Op(":").Add(gen.FormatValue(Id(recvName).Dot("params").Dot(param.FldNameUnExport), param.Type, f.Import, endpoint.TimeFormat)),
							))
						}
					}),
					For(List(Id("_"), Id("before")).Op(":=").Range().Id(recvName).Dot("methodOpts").Dot("before")).Block(
						Id("ctx").Op("=").Id("before").Call(Id("ctx"), Id("req")),
					),
					List(Id("resp"), Err()).Op(":=").Id(recvName).Dot("client").Dot("Do").Call(Id("req")),
					Do(gen.CheckErr(
						Id("cancel").Call(),
						Return(),
					)),
					For(List(Id("_"), Id("after")).Op(":=").Range().Id(recvName).Dot("methodOpts").Dot("after")).Block(
						Id("ctx").Op("=").Id("after").Call(Id("ctx"), Id("resp")),
					),
					Defer().Id("resp").Dot("Body").Dot("Close").Call(),
					Defer().Id("cancel").Call(),

					If(Id("resp").Dot("StatusCode").Op(">").Lit(399)).BlockFunc(func(g *Group) {
						if errorWrapper != nil {
							g.Var().Id("errorWrapper").Do(f.Import(errorWrapper.Struct.Named.Pkg.Path, errorWrapper.Struct.Named.Name))
							g.Id("err").Op("=").Qual("encoding/json", "NewDecoder").
								Call(Id("resp").Dot("Body")).Dot("Decode").
								Call(Op("&").Id("errorWrapper"))
							g.Do(gen.CheckErr(
								Return(),
							))

							g.Id("err").Op("=").Op("&").Do(f.Import(errorWrapper.Default.Named.Pkg.Path, errorWrapper.Default.Named.Name)).ValuesFunc(func(g *Group) {
								for _, field := range errorWrapper.Fields {
									g.Id(strcase.ToCamel(field.FldName)).Op(":").Id("errorWrapper").Dot(field.FldName)
								}
								var statusCodeFound bool
								for _, field := range errorWrapper.Default.Type.Fields {
									if t, ok := field.Var.Type.(*types.Basic); ok && field.Var.Name == "StatusCode" && t.IsInt() {
										statusCodeFound = true
										break
									}
								}
								if statusCodeFound {
									g.Id("StatusCode").Op(":").Id("resp").Dot("StatusCode")
								}
							})
						}
						g.Return()
					}),
					Do(func(s *Statement) {
						if !endpoint.DisabledWrapResponse {
							s.Var().Id("respBody").StructFunc(gen.WrapResponse(endpoint.WrapResponse, endpoint.BodyResults, f.Import))
						} else {
						}
					}),
					Id("err").Op("=").Qual("encoding/json", "NewDecoder").
						Call(Id("resp").Dot("Body")).Dot("Decode").
						Call(Op("&").Id("respBody")),
					Do(gen.CheckErr(
						Return(),
					)),
					ReturnFunc(func(g *Group) {
						var ids []Code
						for _, name := range endpoint.WrapResponse {
							ids = append(ids, Dot(strcase.ToCamel(name)))
						}

						for _, result := range endpoint.Sig.Results {
							if result.IsError {
								g.Id(result.Name)
								continue
							}
							g.Id("respBody").Add(ids...).Dot(strcase.ToCamel(result.Name))
						}
					}),
				)
		}

		for _, endpoint := range s.Endpoints {
			methodRequestName := s.Name + endpoint.MethodName + "Request"
			recvName := strcase.ToLowerCamel(s.Name)
			f.Func().Params(Id(recvName).Op("*").Id(clientName)).Id(endpoint.MethodName).
				ParamsFunc(func(g *Group) {
					for _, param := range endpoint.Params {
						if param.Required {
							g.Id(param.FldNameUnExport).Add(types.Convert(param.Type, f.Import))
						}
					}
				}).
				Op("*").Id(methodRequestName).BlockFunc(func(g *Group) {
				g.Id("m").Op(":=").Op("&").Id(methodRequestName).Values(
					Id("client").Op(":").Id(recvName).Dot("client"),
					Id("methodOpts").Op(":").Op("&").Id("clientMethodOptions").Values(
						Id("ctx").Op(":").Qual("context", "TODO").Call(),
					),
					Id("c").Op(":").Id(recvName),
				)
				for _, param := range endpoint.Params {
					if param.Required {
						g.Id("m").Dot("params").Dot(param.FldNameUnExport).Op("=").Id(param.FldNameUnExport)
					}
				}
				g.Return(Id("m"))
			})
		}
		f.Func().Id("New" + s.Name + "Client").Params(Id("target").String()).Op("*").Id(clientName).BlockFunc(
			func(g *Group) {
				g.Id("c").Op(":=").Op("&").Id(clientName).Values(
					Id("target").Op(":").Id("target"),
					Id("client").Op(":").Qual("net/http", "DefaultClient"),
				)
				g.Return(Id("c"))
			},
		)
	}
}
