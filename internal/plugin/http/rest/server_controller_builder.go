package rest

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/gen"
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

type serverControllerBuilder struct {
	*BaseServerBuilder
	handlerStrategy HandlerStrategy
	iface           options.Iface
}

func (b *serverControllerBuilder) BuildHandlers() ServerControllerBuilder {
	middlewareType := b.handlerStrategy.MiddlewareType()

	errorEncoderType := jen.Func().Params(
		jen.Id(b.handlerStrategy.RespArgName()).Add(b.handlerStrategy.RespType()),
		jen.Err().Error(),
	)

	optionsName := b.iface.Name + "Options"
	optionName := b.iface.Name + "Option"
	reqName := "req"
	respName := "resp"

	b.codes = append(b.codes,
		jen.Type().Id(optionName).Func().Params(jen.Op("*").Id(optionsName)),
		jen.Type().Id(optionsName).StructFunc(func(g *jen.Group) {
			g.Id("errorEncoder").Add(errorEncoderType)
			g.Id("middleware").Index().Add(middlewareType)
			for _, ep := range b.iface.Endpoints {
				g.Id("middleware" + ep.MethodName).Index().Add(middlewareType)
			}
		}),
		jen.Func().Id(b.iface.Name+"Middleware").Params(jen.Id("middleware").Op("...").Add(middlewareType)).Id(optionName).Block(
			jen.Return(
				jen.Func().Params(jen.Id("o").Op("*").Id(optionsName)).Block(
					jen.Id("o").Dot("middleware").Op("=").Append(jen.Id("o").Dot("middleware"), jen.Id("middleware").Op("...")),
				),
			),
		),

		jen.Func().Id(b.iface.Name+"WithErrorEncoder").Params(jen.Id("errorEncoder").Add(errorEncoderType)).Id(optionName).Block(
			jen.Return(
				jen.Func().Params(jen.Id("o").Op("*").Id(optionsName)).Block(
					jen.Id("o").Dot("errorEncoder").Op("=").Id("errorEncoder"),
				),
			),
		),
	)

	for _, ep := range b.iface.Endpoints {
		b.codes = append(b.codes,
			jen.Func().Id(b.iface.Name+ep.MethodName+"Middleware").Params(jen.Id("middleware").Op("...").Add(middlewareType)).Id(optionName).Block(
				jen.Return(
					jen.Func().Params(jen.Id("o").Op("*").Id(optionsName)).Block(
						jen.Id("o").Dot("middleware"+ep.MethodName).Op("=").Append(jen.Id("o").Dot("middleware"+ep.MethodName), jen.Id("middleware").Op("...")),
					),
				),
			),
		)
	}

	b.codes = append(b.codes,
		jen.Func().Id("SetupRoutes"+b.iface.Name).Params(
			jen.Id("svc").Do(b.qualifier.Qual(b.iface.PkgPath, b.iface.Name)),
			jen.Id(b.handlerStrategy.LibArgName()).Add(b.handlerStrategy.LibType()),
			jen.Id("opts").Op("...").Id(b.iface.Name+"Option"),
		).BlockFunc(func(g *jen.Group) {
			g.Id("o").Op(":=").Op("&").Id(b.iface.Name + "Options").Values(
				jen.Id("errorEncoder").Op(":").Id(b.handlerStrategy.ID() + "DefaultErrorEncoder"),
			)
			g.For(jen.List(jen.Id("_"), jen.Id("opt")).Op(":=").Range().Id("opts")).Block(
				jen.Id("opt").Call(jen.Id("o")),
			)
			for _, ep := range b.iface.Endpoints {
				handlerFunc := func(g *jen.Group) {
					if len(ep.Results) > 0 {
						g.Var().Id("err").Error()
					}
					if len(ep.Params) > 0 {
						if len(ep.BodyParams) > 0 {
							var stRequests []jen.Code
							if ep.ReqRootXMLName != "" {
								stRequests = append(stRequests, jen.Id("XMLName").Qual("encoding/xml", "Name").Tag(map[string]string{"xml": ep.ReqRootXMLName}))
							}
							for _, p := range ep.BodyParams {
								st := jen.Id(p.FldName).Add(types.Convert(p.Type, b.qualifier.Qual))
								if p.Name != "" && p.HTTPType == "body" {
									st.Tag(map[string]string{"json": p.Name})
								} else {
									st.Tag(map[string]string{"json": "-"})
								}
								stRequests = append(stRequests, st)
							}

							g.Var().Id(reqName).Struct(stRequests...)

							switch ep.HTTPMethod {
							case "POST", "PUT", "PATCH", "DELETE":
								nameContentType, typ := b.handlerStrategy.HeaderParam("content-type")

								g.Add(typ)

								g.Id("parts").Op(":=").Qual("strings", "Split").Call(jen.Id(nameContentType), jen.Lit(";"))
								g.If(jen.Len(jen.Id("parts")).Op(">").Lit(0)).Block(
									jen.Id(nameContentType).Op("=").Id("parts").Index(jen.Lit(0)),
								)

								g.Var().Id("bodyData").Op("=").Make(jen.Index().Byte(), jen.Lit(0), jen.Lit(10485760)) // 10MB
								g.Id("buf").Op(":=").Qual("bytes", "NewBuffer").Call(jen.Id("bodyData"))
								g.List(jen.Id("written"), jen.Id("err")).Op(":=").Qual("io", "Copy").Call(jen.Id("buf"), b.handlerStrategy.BodyPathParam())
								g.Do(gen.CheckErr(
									jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.ReqArgName()), jen.Err()),
									jen.Return(),
								))
								g.Switch(jen.Id(nameContentType)).BlockFunc(func(g *jen.Group) {
									g.Default().Block(
										jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.ReqArgName()), jen.Op("&").Id("contentTypeInvalidError").Values()),
										jen.Return(),
									)
									for _, contentType := range ep.ContentTypes {
										bodyParams := jen.Id(reqName)
										if len(ep.BodyParams) == 1 && ep.NoWrapRequest {
											bodyParams = bodyParams.Dot(ep.BodyParams[0].FldName)
										}

										switch contentType {
										case "xml":
											g.Case(jen.Lit("application/xml")).BlockFunc(func(g *jen.Group) {
												g.Err().Op("=").Qual("encoding/xml", "Unmarshal").Call(jen.Id("bodyData").Index(jen.Op(":").Id("written")), jen.Op("&").Add(bodyParams))
												g.Do(gen.CheckErr(
													jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.ReqArgName()), jen.Err()),
													jen.Return(),
												))
											})
										case "json":
											g.Case(jen.Lit("application/json")).BlockFunc(func(g *jen.Group) {
												g.Err().Op("=").Qual("encoding/json", "Unmarshal").Call(jen.Id("bodyData").Index(jen.Op(":").Id("written")), jen.Op("&").Add(bodyParams))
												g.Do(gen.CheckErr(
													jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.ReqArgName()), jen.Err()),
													jen.Return(),
												))
											})
										case "urlencoded":
											g.Case(jen.Lit("application/x-www-form-urlencoded")).BlockFunc(func(g *jen.Group) {
												typ, hasErr := b.handlerStrategy.FormParams()

												g.Add(typ)

												if hasErr {
													g.Do(gen.CheckErr(
														jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.ReqArgName()), jen.Err()),
														jen.Return(),
													))
												}

												for _, p := range ep.BodyParams {
													_, typ := b.handlerStrategy.FormParam(p.Name)

													g.Add(gen.ParseValue(typ, jen.Add(bodyParams).Dot(p.FldName), "=", p.Type, b.qualifier.Qual, func() jen.Code {
														return jen.Do(gen.CheckErr(
															jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.ReqArgName()), jen.Err()),
															jen.Return(),
														))
													}))
												}
											})
										case "multipart":
											g.Case(jen.Lit("multipart/form-data")).BlockFunc(func(g *jen.Group) {
												typ, hasErr := b.handlerStrategy.MultipartFormParams(ep.MultipartMaxMemory)
												g.Add(typ)
												if hasErr {
													g.Do(gen.CheckErr(
														jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.ReqArgName()), jen.Err()),
														jen.Return(),
													))
												}
												for _, p := range ep.BodyParams {
													_, typ := b.handlerStrategy.MultipartFormParam(p.Name)

													g.Add(gen.ParseValue(typ, jen.Add(bodyParams).Dot(p.FldName), "=", p.Type, b.qualifier.Qual, func() jen.Code {
														return jen.Do(gen.CheckErr(
															jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.ReqArgName()), jen.Err()),
															jen.Return(),
														))
													}))
												}
											})
										}
									}
								})
							}
						}

						buildParams := func(g *jen.Group, params options.EndpointParams, f func(pathName string) (name string, typ jen.Code)) {
							for _, p := range ep.PathParams {
								paramName, typ := f(p.Name)
								paramVarName := "param" + p.FldName

								g.Add(typ)

								g.Var().Id(paramVarName).Add(types.Convert(p.Type, b.qualifier.Qual))

								g.If(jen.Id(paramName).Op("!=").Lit("")).Block(
									jen.Add(gen.ParseValue(jen.Id(paramName), jen.Id(paramVarName), "=", p.Type, b.qualifier.Qual, func() jen.Code {
										return jen.Do(gen.CheckErr(
											jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.ReqArgName()), jen.Err()),
											jen.Return(),
										))
									})),
								)
							}
						}

						if b.handlerStrategy.UsePathParams() && len(ep.PathParams) > 0 {
							buildParams(g, ep.PathParams, b.handlerStrategy.PathParam)
						}
						if len(ep.HeaderParams) > 0 {
							buildParams(g, ep.HeaderParams, b.handlerStrategy.HeaderParam)
						}
						if len(ep.QueryParams) > 0 {
							buildParams(g, ep.QueryParams, b.handlerStrategy.QueryParam)
						}
					}

					g.Do(func(s *jen.Statement) {
						s.ListFunc(func(g *jen.Group) {
							for _, r := range ep.Results {
								g.Id(r.FldName)
							}
							if ep.Error != nil {
								g.Id(ep.Error.Name)
							}
						})

						if len(ep.Results) > 0 {
							s.Op(":=")
						} else {
							s.Op("=")
						}
					}).Id("svc").Dot(ep.MethodName).CallFunc(func(g *jen.Group) {
						if ep.Context != nil {
							g.Id("ctx")
						}
						for _, p := range ep.Params {
							switch p.HTTPType {
							default:
								g.Id("req").Dot(p.FldName)
							case options.PathHTTPType, options.CookieHTTPType, options.QueryHTTPType:
								g.Id("param" + p.FldName)
							}
						}
					})
					if ep.Error != nil {
						g.Do(gen.CheckErr(
							jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.ReqArgName()), jen.Err()),
							jen.Return(),
						))
					}

					if len(ep.BodyResults) > 0 {
						var (
							responseFields []jen.Code
							assignFields   []jen.Code
						)

						for _, p := range ep.BodyResults {
							fld := jen.Id(p.FldNameExport).Add(types.Convert(p.Type, b.qualifier.Qual))
							if p.Name != "" && p.HTTPType == "body" {
								fld.Tag(map[string]string{"json": p.Name})
							} else {
								fld.Tag(map[string]string{"json": "-"})
							}
							responseFields = append(responseFields, fld)
							assignFields = append(assignFields, jen.Id("resp").Dot(p.FldNameExport).Op("=").Id(p.FldName))
						}
						g.Var().Id(respName).Struct(responseFields...)

						g.Add(assignFields...)

						b.handlerStrategy.WriteBody(jen.Id(respName))

						// if !ep.NoWrapResponse {
						// 	g.Var().Id("wrapResult").StructFunc(gen.WrapResponse(ep.WrapResponse, ep.BodyResults, b.qualifier.Qual))
						// 	for _, r := range ep.BodyResults {
						// 		g.Id("wrapResult").Do(func(s *jen.Statement) {
						// 			for _, name := range ep.WrapResponse {
						// 				s.Dot(strcase.ToCamel(name))
						// 			}
						// 		}).Dot(r.FldNameExport).Op("=").Id("resp").Dot(r.FldNameExport)
						// 	}

						// 	g.Id("result").Op("=").Id("wrapResult")

						// } else if len(ep.BodyResults) == 1 {
						// 	// g.Id("result").Op("=").Id("result").Assert(jen.Op("*").Id(respName)).Dot(ep.BodyResults[0].FldNameExport)
						// }

					}
				}

				g.Add(b.handlerStrategy.HandlerFunc(ep.HTTPMethod, ep.Path, handlerFunc))

				// g.Add(b.handlerStrategy.HandlerFunc(
				// 	ep.HTTPMethod,
				// 	ep.Path,
				// 	resultID,
				// 	jen.Append(jen.Id("o").Dot("middleware"), jen.Id("o").Dot("middleware"+ep.MethodName).Op("...")),
				// 	handlerFuncBody,
				// ))

				// epName := endpointFuncName(b.iface, ep)
				// reqDecName := reqDecFuncName(b.iface, ep)
				// respEncName := respEncFuncName(b.iface, ep)

				// hasResults := len(ep.Results) > 0
				// hasParams := len(ep.Params) > 0
				// epResultOp := ":="
				// epResultName := "resp"

				// if hasParams && !hasResults {
				// epResultOp = "="
				// }
				// if !hasResults {
				// epResultName = "_"
				// }

				// var resultID jen.Code

				// errorEncoder := jen.Id("o").Dot("errorEncoder").Call(
				// jen.Id(b.handlerStrategy.RespArgName()),
				// jen.Err(),
				// )

				// handlerFuncBody := jen.CustomFunc(jen.Options{Multi: true}, func(g *jen.Group) {
				// 	g.Id("reqCtx").Op(":=").Qual(contextPkg, "TODO").Call()

				// 	if hasParams {
				// 		g.List(jen.Id("req"), jen.Err()).Op(":=").Id(reqDecName).Call(jen.Id(b.handlerStrategy.ReqArgName()))
				// 		g.Do(gen.CheckErr(jen.Return()))
				// 	}
				// 	// g.List(jen.Id(epResultName), jen.Err()).Op(epResultOp).Id(epName).Call(jen.Id("svc")).CallFunc(func(g *jen.Group) {
				// 	// 	g.Id("reqCtx")
				// 	// 	if hasParams {
				// 	// 		g.Id("req")
				// 	// 	} else {
				// 	// 		g.Nil()
				// 	// 	}
				// 	// })
				// 	// g.Do(gen.CheckErr(
				// 	// errorEncoder,
				// 	// jen.Return(),
				// 	// ))
				// 	// if hasResults {
				// 	// resultID = jen.Id("result")
				// 	// g.List(resultID, jen.Err()).Op(":=").Id(respEncName).Call(jen.Id("resp"))
				// 	// g.Do(gen.CheckErr(
				// 	// errorEncoder,
				// 	// jen.Return(),
				// 	// ))
				// 	// }
				// })

				// g.Add(b.handlerStrategy.HandlerFunc(
				// 	ep.HTTPMethod,
				// 	ep.Path,
				// 	resultID,
				// 	jen.Append(jen.Id("o").Dot("middleware"), jen.Id("o").Dot("middleware"+ep.MethodName).Op("...")),
				// 	handlerFuncBody,
				// ))
			}
		}))
	return b
}

func (b *serverControllerBuilder) Endpoint(ep options.Endpoint) ServerEndpointBuilder {
	return &serverEndpointBuilder{serverControllerBuilder: b, handlerStrategy: b.handlerStrategy, iface: b.iface, ep: ep, qualifier: b.qualifier}
}
