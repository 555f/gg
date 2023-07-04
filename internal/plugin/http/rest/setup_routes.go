package rest

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gen"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	. "github.com/dave/jennifer/jen"
)

func GenStruct(s options.Iface) func(f *file.GoFile) {
	return func(f *file.GoFile) {

		for _, ep := range s.Endpoints {
			var (
				stRequests  []Code
				stResponses []Code
			)

			reqDecName := strcase.ToLowerCamel(s.Name+ep.MethodName) + "ReqDec"
			respEncName := strcase.ToLowerCamel(s.Name+ep.MethodName) + "RespEnc"
			reqName := strcase.ToLowerCamel(s.Name+ep.MethodName) + "Req"
			respName := strcase.ToLowerCamel(s.Name+ep.MethodName) + "Resp"
			epName := strcase.ToLowerCamel(s.Name+ep.MethodName) + "Endpoint"

			if ep.ReqRootXMLName != "" {
				stRequests = append(stRequests, Id("XMLName").Qual("encoding/xml", "Name").Tag(map[string]string{"xml": ep.ReqRootXMLName}))
			}

			for _, p := range ep.Params {
				st := Id(p.FldName).Add(types.Convert(p.Type, f.Import))
				if p.Name != "" && p.HTTPType == "body" {
					st.Tag(map[string]string{"json": p.Name})
				} else {
					st.Tag(map[string]string{"json": "-"})
				}
				stRequests = append(stRequests, st)
			}

			for _, p := range ep.Results {
				st := Id(p.FldNameExport).Add(types.Convert(p.Type, f.Import))
				if p.Name != "" && p.HTTPType == "body" {
					st.Tag(map[string]string{"json": p.Name})
				} else {
					st.Tag(map[string]string{"json": "-"})
				}
				stResponses = append(stResponses, st)
			}

			if len(ep.Params) > 0 {
				f.Type().Id(reqName).Struct(stRequests...)
			}
			if len(ep.Results) > 0 {
				f.Type().Id(respName).Struct(stResponses...)
			}

			f.Func().Id(epName).Params(
				Id("svc").Qual(s.PkgPath, s.Name),
			).Add(epFunc).Block(
				Return(Add(epFunc).BlockFunc(func(g *Group) {
					if len(ep.Params) > 0 {
						g.Id("r").Op(":=").Id("request").Assert(Op("*").Id(reqName))
					}
					g.Do(func(s *Statement) {
						s.ListFunc(func(g *Group) {
							for _, r := range ep.Results {
								g.Id(r.FldName)
							}
							if ep.Error != nil {
								g.Id(ep.Error.Name)
							}
						})
						s.Op(":=")
					}).Id("svc").Dot(ep.MethodName).CallFunc(func(g *Group) {
						if ep.Context != nil {
							g.Id("ctx")
						}
						for _, p := range ep.Params {
							g.Id("r").Dot(p.FldName)
						}
					})
					if ep.Error != nil {
						g.Do(gen.CheckErr(Return(Nil(), Err())))
					}
					g.ReturnFunc(func(g *Group) {
						if len(ep.Results) > 0 {
							g.Op("&").Id(respName).ValuesFunc(func(g *Group) {
								for _, p := range ep.Results {
									g.Id(p.FldNameExport).Op(":").Id(p.FldName)
								}
							})
							g.Nil()
						} else {
							g.Nil()
							g.Nil()
						}
					})
				})),
			)

			if len(ep.Results) > 0 {
				f.Func().Id(respEncName).Params(
					Id("result").Any(),
				).Params(
					Any(),
					Error(),
				).BlockFunc(func(g *Group) {
					if len(ep.BodyResults) > 0 {
						if !ep.NoWrapResponse {
							g.Var().Id("wrapResult").StructFunc(gen.WrapResponse(ep.WrapResponse, ep.Results, f.Import))
							for _, r := range ep.Results {
								g.Id("wrapResult").Do(func(s *Statement) {
									for _, name := range ep.WrapResponse {
										s.Dot(strcase.ToCamel(name))
									}
								}).Dot(r.FldNameExport).Op("=").Id("result").Assert(Op("*").Id(respName)).Dot(r.FldNameExport)
							}

							g.Id("result").Op("=").Id("wrapResult")

						} else if len(ep.BodyResults) == 1 {
							g.Id("result").Op("=").Id("result").Assert(Op("*").Id(respName)).Dot(ep.BodyResults[0].FldNameExport)
						}
						g.Return(Id("result"), Nil())
					} else {
						g.Return(Nil(), Nil())
					}
				})
			}

			if len(ep.Params) > 0 {
				f.Func().Id(reqDecName).ParamsFunc(func(g *Group) {
					g.Id("pathParams").Id("pathParams")
					g.Id("r").Op("*").Qual("net/http", "Request")
					g.Do(func(s *Statement) {
						if len(ep.BodyParams) > 0 {
							s.Id("params")
						} else {
							s.Id("_")
						}
					}).Qual("encoding/json", "RawMessage")
				}).Params(
					Id("result").Any(),
					Err().Error(),
				).BlockFunc(func(g *Group) {
					g.Var().Id("param").Id(reqName)

					if len(ep.BodyParams) > 0 {
						switch ep.HTTPMethod {
						case "POST", "PUT", "PATCH", "DELETE":
							g.Id("contentType").Op(":=").Id("r").Dot("Header").Dot("Get").Call(Lit("content-type"))

							g.Id("parts").Op(":=").Qual("strings", "Split").Call(Id("contentType"), Lit(";"))
							g.If(Len(Id("parts")).Op("==").Lit(0)).Block(
								Return(Nil(), Err()),
							)
							g.Id("contentType").Op("=").Id("parts").Index(Lit(0))

							g.Switch(Id("contentType")).BlockFunc(func(g *Group) {
								g.Default().Block(
									Return(Nil(), Op("&").Id("contentTypeInvalidError").Values()),
								)
								for _, contentType := range ep.ContentTypes {
									switch contentType {
									case "xml":
										g.Case(Lit("application/xml")).BlockFunc(func(g *Group) {
											g.Err().Op("=").Qual("encoding/xml", "Unmarshal").Call(Id("params"), Op("&").Id("param"))
											g.Do(gen.CheckErr(
												Return(Nil(), Err()),
											))
										})
									case "json":
										g.Case(Lit("application/json")).BlockFunc(func(g *Group) {
											g.Err().Op("=").Qual("encoding/json", "Unmarshal").Call(Id("params"), Op("&").Id("param"))
											g.Do(gen.CheckErr(
												Return(Nil(), Err()),
											))
										})
									case "urlencoded":
										g.Case(Lit("application/x-www-form-urlencoded")).BlockFunc(func(g *Group) {
											g.Err().Op("=").Id("r").Dot("ParseForm").Call()
											g.Do(gen.CheckErr(
												Return(Nil(), Err()),
											))
											for _, p := range ep.BodyParams {
												g.Add(gen.ParseValue(Id("r").Dot("Form").Dot("Get").Call(Lit(p.Name)), Id("param").Dot(p.FldName), "=", p.Type, f.Import))
												if b, ok := p.Type.(*types.Basic); (ok && !b.IsString()) || !ok {
													g.Do(gen.CheckErr(
														Return(Nil(), Err()),
													))
												}
											}
										})
									case "multipart":
										g.Case(Lit("multipart/form-data")).BlockFunc(func(g *Group) {
											g.Err().Op("=").Id("r").Dot("ParseMultipartForm").Call(Lit(ep.MultipartMaxMemory))
											g.Do(gen.CheckErr(
												Return(Nil(), Err()),
											))
											for _, p := range ep.BodyParams {
												g.Add(gen.ParseValue(Id("r").Dot("FormValue").Call(Lit(p.Name)), Id("param").Dot(p.FldName), "=", p.Type, f.Import))
												if b, ok := p.Type.(*types.Basic); (ok && !b.IsString()) || !ok {
													g.Do(gen.CheckErr(
														Return(Nil(), Err()),
													))
												}
											}
										})
									}
								}
							})
						}
					}

					if len(ep.PathParams) > 0 {
						for _, p := range ep.PathParams {
							g.If(Id("s").Op(":=").Id("pathParams").Dot("Param").Call(Lit(p.Name)), Id("s").Op("!=").Lit("")).Block(
								Add(gen.ParseValue(Id("s"), Id("param").Dot(p.FldName), "=", p.Type, f.Import)),
								Do(gen.CheckErr(Return())),
							)
						}
					}

					if len(ep.QueryParams) > 0 {
						g.Id("q").Op(":=").Id("r").Dot("URL").Dot("Query").Call()
						for _, param := range ep.QueryParams {
							paramID := Id("param").Dot(param.FldName)
							if param.Parent != nil {
								paramID = Id("param").Dot(param.Parent.FldName).Dot(param.FldName)
							}
							g.If(Id("s").Op(":=").Id("q").Dot("Get").Call(Lit(param.Name)), Id("s").Op("!=").Lit("")).Block(
								Add(gen.ParseValue(Id("s"), paramID, "=", param.Type, f.Import)),
								Do(gen.CheckErr(Return())),
							)
						}
					}

					if len(ep.HeaderParams) > 0 {
						for _, p := range ep.HeaderParams {
							g.If(Id("s").Op(":=").Id("r").Dot("Header").Dot("Get").Call(Lit(p.Name)), Id("s").Op("!=").Lit("")).Block(
								Add(gen.ParseValue(Id("s"), Id("param").Dot(p.FldName), "=", p.Type, f.Import)),
								Do(gen.CheckErr(Return())),
							)
						}
					}
					if len(ep.Params) > 0 {
						g.Return(Op("&").Id("param"), Nil())
					} else {
						g.Return()
					}
				})
			}
		}
		f.Commentf("// SetupRoutes%s route init for service", s.Name)
		f.Func().Id("SetupRoutes" + s.Name).ParamsFunc(func(g *Group) {
			g.Id("svc").Do(f.Import(s.PkgPath, s.Name))
			switch s.Type {
			case "rest":
				switch s.Lib {
				case "echo":
					g.Id("s").Op("*").Qual("github.com/labstack/echo/v4", "Echo")
				case "http":
					g.Id("s").Op("*").Qual("net/http", "ServeMux")
				}
			case "jsonrpc":
				switch s.Lib {
				case "std":
					g.Id("s").Op("*").Qual("github.com/555f/jsonrpc", "Server")
				}
			}
			g.Id("opts").Op("...").Id(s.Name + "Option")
		}).BlockFunc(func(g *Group) {
			g.Id("o").Op(":=").Op("&").Id(s.Name + "Options").Values()
			g.For(List(Id("_"), Id("opt")).Op(":=").Range().Id("opts")).Block(
				Id("opt").Call(Id("o")),
			)

			for _, ep := range s.Endpoints {
				epName := strcase.ToLowerCamel(s.Name+ep.MethodName) + "Endpoint"
				reqDecName := strcase.ToLowerCamel(s.Name+ep.MethodName) + "ReqDec"
				respEncName := strcase.ToLowerCamel(s.Name+ep.MethodName) + "RespEnc"
				switch s.Type {
				case "rest":
					switch s.Lib {
					case "http":
						g.Id("s").Dot("Handle").Call(
							Lit(ep.Path),

							Qual("net/http", "HandlerFunc").Call(
								Func().Params(
									Id("w").Qual("net/http", "ResponseWriter"),
									Id("r").Op("*").Qual("net/http", "Request"),
								).Block(

									If(Id("r").Dot("Method").Op("==").Lit(ep.HTTPMethod)).Block(
										Id("httpHandler").CallFunc(func(g *Group) {
											g.Id(epName).Call(Id("svc"))
											if len(ep.Params) > 0 {
												g.Id(reqDecName)
											} else {
												g.Nil()
											}
											if len(ep.Results) > 0 {
												g.Id(respEncName)
											} else {
												g.Nil()
											}
											g.Op("&").Id("pathParamsNoop").Values()
										}).Dot("ServeHTTP").Call(Id("w"), Id("r")),
									),
								),
							),
						)
					case "echo":
						g.Id("s").Dot("Add").Call(
							Lit(ep.HTTPMethod),
							Lit(ep.Path),
							Func().Params(
								Id("ctx").Qual("github.com/labstack/echo/v4", "Context"),
							).Error().Block(
								Id("httpHandler").CallFunc(func(g *Group) {
									g.Id(epName).Call(Id("svc"))
									if len(ep.Params) > 0 {
										g.Id(reqDecName)
									} else {
										g.Nil()
									}
									if len(ep.Results) > 0 {
										g.Id(respEncName)
									} else {
										g.Nil()
									}
									g.Id("ctx")
								}).Dot("ServeHTTP").Call(Id("ctx").Dot("Response").Call().Dot("Writer"), Id("ctx").Dot("Request").Call()),
								Return(Nil()),
							),
							Append(Id("o").Dot("middleware"), Id("o").Dot("middleware"+ep.MethodName).Op("...")).Op("..."),
						)
					}
				case "jsonrpc":
					switch s.Lib {
					case "std":
						g.Id("s").Dot("Register").Call(
							Lit(ep.Path),
							Id(epName).Call(Id("svc")),
							Id(reqDecName),
						)
					}
				}
			}
		}).Line()
	}
}
