package rest

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/gen"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

type serverControllerBuilder struct {
	*BaseServerBuilder
	handlerStrategy HandlerStrategy
	iface           options.Iface
}

func (b *serverControllerBuilder) Build() ServerControllerBuilder {
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
					g.Var().Id("err").Error()

					if len(ep.Params) > 0 {
						if len(ep.BodyParams) > 0 {
							bodyParams := jen.Id(reqName)
							if len(ep.BodyParams) == 1 && ep.NoWrapRequest {
								bodyParams = bodyParams.Dot(ep.BodyParams[0].FldName.String())
							}

							var stRequests []jen.Code
							if ep.ReqRootXMLName != "" {
								stRequests = append(stRequests, jen.Id("XMLName").Qual("encoding/xml", "Name").Tag(map[string]string{"xml": ep.ReqRootXMLName}))
							}
							for _, p := range ep.BodyParams {
								st := jen.Id(p.FldName.String()).Add(types.Convert(p.Type, b.qualifier.Qual))
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
								g.Var().Id("bodyData").Op("=").Make(jen.Index().Byte(), jen.Lit(0), jen.Lit(10485760)) // 10MB
								g.Id("buf").Op(":=").Qual("bytes", "NewBuffer").Call(jen.Id("bodyData"))
								g.List(jen.Id("written"), jen.Id("err")).Op(":=").Qual("io", "Copy").Call(jen.Id("buf"), b.handlerStrategy.BodyPathParam())
								g.Do(gen.CheckErr(
									jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.RespArgName()), jen.Err()),
									jen.Return(),
								))
								if len(ep.ContentTypes) > 0 {
									nameContentType, typ := b.handlerStrategy.HeaderParam("content-type")
									g.Add(typ)

									g.Id("parts").Op(":=").Qual("strings", "Split").Call(jen.Id(nameContentType), jen.Lit(";"))
									g.If(jen.Len(jen.Id("parts")).Op(">").Lit(0)).Block(
										jen.Id(nameContentType).Op("=").Id("parts").Index(jen.Lit(0)),
									)

									g.Switch(jen.Id(nameContentType)).BlockFunc(func(g *jen.Group) {
										g.Default().BlockFunc(func(g *jen.Group) {
											g.Err().Op("=").Qual("encoding/json", "Unmarshal").Call(jen.Id("bodyData").Index(jen.Op(":").Id("written")), jen.Op("&").Add(bodyParams))
											g.Do(gen.CheckErr(
												jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.RespArgName()), jen.Err()),
												jen.Return(),
											))
										})

										for _, contentType := range ep.ContentTypes {
											switch contentType {
											case "xml":
												g.Case(jen.Lit("application/xml")).BlockFunc(func(g *jen.Group) {
													g.Err().Op("=").Qual("encoding/xml", "Unmarshal").Call(jen.Id("bodyData").Index(jen.Op(":").Id("written")), jen.Op("&").Add(bodyParams))
													g.Do(gen.CheckErr(
														jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.RespArgName()), jen.Err()),
														jen.Return(),
													))
												})
											case "urlencoded":
												g.Case(jen.Lit("application/x-www-form-urlencoded")).BlockFunc(func(g *jen.Group) {
													typ, hasErr := b.handlerStrategy.FormParams()

													g.Add(typ)

													if hasErr {
														g.Do(gen.CheckErr(
															jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.RespArgName()), jen.Err()),
															jen.Return(),
														))
													}

													for _, p := range ep.BodyParams {
														_, typ := b.handlerStrategy.FormParam(p.Name)

														g.Add(gen.ParseValue(typ, jen.Add(bodyParams).Dot(p.FldName.String()), "=", p.Type, b.qualifier.Qual, func() jen.Code {
															return jen.Do(gen.CheckErr(
																jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.RespArgName()), jen.Err()),
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
															jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.RespArgName()), jen.Err()),
															jen.Return(),
														))
													}
													for _, p := range ep.BodyParams {
														_, typ := b.handlerStrategy.MultipartFormParam(p.Name)

														g.Add(gen.ParseValue(typ, jen.Add(bodyParams).Dot(p.FldName.String()), "=", p.Type, b.qualifier.Qual, func() jen.Code {
															return jen.Do(gen.CheckErr(
																jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.RespArgName()), jen.Err()),
																jen.Return(),
															))
														}))
													}
												})
											}
										}

									})
								} else {
									g.Err().Op("=").Qual("encoding/json", "Unmarshal").Call(jen.Id("bodyData").Index(jen.Op(":").Id("written")), jen.Op("&").Add(bodyParams))
									g.Do(gen.CheckErr(
										jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.RespArgName()), jen.Err()),
										jen.Return(),
									))
								}
							}
						}

						buildParams := func(g *jen.Group, params options.EndpointParams, f func(pathName string) (name string, typ jen.Code)) {
							for _, p := range params {
								paramName, typ := f(p.Name)
								paramVarName := "param" + p.FldName.String()

								g.Add(typ)

								g.Var().Id(paramVarName).Add(types.Convert(p.Type, b.qualifier.Qual))

								g.If(jen.Id(paramName).Op("!=").Lit("")).Block(
									jen.Add(gen.ParseValue(jen.Id(paramName), jen.Id(paramVarName), "=", p.Type, b.qualifier.Qual, func() jen.Code {
										return jen.Do(gen.CheckErr(
											jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.RespArgName()), jen.Err()),
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
							g.Add(b.handlerStrategy.QueryParams())

							buildParams(g, ep.QueryParams, b.handlerStrategy.QueryParam)
						}
					}

					g.Do(func(s *jen.Statement) {
						s.ListFunc(func(g *jen.Group) {
							for _, r := range ep.Results {
								g.Id(r.FldName.String())
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
							g.Add(b.handlerStrategy.Context())
						}
						for _, p := range ep.Params {
							switch p.HTTPType {
							default:
								g.Id("req").Dot(p.FldName.String())
							case options.PathHTTPType, options.CookieHTTPType, options.QueryHTTPType:
								g.Id("param" + p.FldName.String())
							}
						}
					})
					if ep.Error != nil {
						g.Do(gen.CheckErr(
							jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.RespArgName()), jen.Err()),
							jen.Return(),
						))
					}

					if len(ep.BodyResults) > 0 {
						if !ep.NoWrapResponse {
							g.Var().Id(respName).StructFunc(gen.WrapResponse(ep.WrapResponse, func(g *jen.Group) {
								for _, result := range ep.BodyResults {
									g.Id(result.FldName.Camel()).Add(types.Convert(result.Type, b.qualifier.Qual)).Tag(map[string]string{"json": result.Name})
								}
							}, b.qualifier.Qual))

							for _, p := range ep.BodyResults {
								g.Id(respName).Do(func(s *jen.Statement) {
									for _, name := range ep.WrapResponse {
										s.Dot(strcase.ToCamel(name))
									}
								}).Dot(p.FldName.Camel()).Op("=").Id(p.FldName.String())
							}
						} else if len(ep.BodyResults) == 1 {
							g.Id(respName).Op(":=").Id(ep.BodyResults[0].FldName.String())
						}

						g.Var().Id("respData").Index().Byte()

						if len(ep.AcceptTypes) > 0 {
							nameAcceptType, typ := b.handlerStrategy.HeaderParam("accept")
							g.Add(typ)

							g.Switch(jen.Id(nameAcceptType)).BlockFunc(func(g *jen.Group) {
								g.Default().Block(
									jen.Id(nameAcceptType).Op("=").Lit("application/json"),
									jen.List(jen.Id("respData"), jen.Err()).Op("=").Qual("encoding/json", "Marshal").Call(jen.Id("resp")),
									jen.Do(gen.CheckErr(
										jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.RespArgName()), jen.Err()),
										jen.Return(),
									)),
								)
								for _, acceptType := range ep.AcceptTypes {
									switch acceptType {
									case "xml":
										g.Case(jen.Lit("application/xml")).Block(
											jen.List(jen.Id("respData"), jen.Err()).Op("=").Qual("encoding/xml", "Marshal").Call(jen.Id("resp")),
											jen.Do(gen.CheckErr(
												jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.RespArgName()), jen.Err()),
												jen.Return(),
											)),
										)
									}
								}
							})
							g.Add(b.handlerStrategy.WriteBody(jen.Id("respData"), jen.Id(nameAcceptType), 200))
						} else {
							g.List(jen.Id("respData"), jen.Err()).Op("=").Qual("encoding/json", "Marshal").Call(jen.Id("resp"))
							g.Do(gen.CheckErr(
								jen.Id("o").Dot("errorEncoder").Call(jen.Id(b.handlerStrategy.RespArgName()), jen.Err()),
								jen.Return(),
							))
							g.Add(b.handlerStrategy.WriteBody(jen.Id("respData"), jen.Lit("application/json"), 200))
						}
					}
				}

				middlewares := jen.Append(jen.Id("o").Dot("middleware"), jen.Id("o").Dot("middleware"+ep.MethodName).Op("...")).Op("...")

				g.Add(b.handlerStrategy.HandlerFunc(ep.HTTPMethod, ep.Path, middlewares, handlerFunc))
			}
		}))
	return b
}
