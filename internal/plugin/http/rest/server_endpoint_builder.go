package rest

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/gen"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

var (
	epFunc = jen.Func().Params(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("request").Any(),
	).Params(
		jen.Any(),
		jen.Error(),
	)
)

type serverEndpointBuilder struct {
	*serverControllerBuilder
	handlerStrategy HandlerStrategy
	iface           options.Iface
	ep              options.Endpoint
	qualifier       Qualifier
}

func (b *serverEndpointBuilder) Build() {
	epFuncName := b.endpointFuncName()
	reqName := b.reqStructName()
	respName := b.respStructName()

	b.codes = append(b.codes,
		jen.Func().Id(epFuncName).Params(
			jen.Id("svc").Qual(b.iface.PkgPath, b.iface.Name),
		).Add(epFunc).Block(
			jen.Return(jen.Add(epFunc).BlockFunc(func(g *jen.Group) {
				if len(b.ep.Params) > 0 {
					g.Id("r").Op(":=").Id("request").Assert(jen.Op("*").Id(reqName))
				}
				g.Do(func(s *jen.Statement) {
					s.ListFunc(func(g *jen.Group) {
						for _, r := range b.ep.Results {
							g.Id(r.FldName)
						}
						if b.ep.Error != nil {
							g.Id(b.ep.Error.Name)
						}
					})
					if len(b.ep.Results) > 0 || b.ep.Error != nil {
						s.Op(":=")
					}
				}).Id("svc").Dot(b.ep.MethodName).CallFunc(func(g *jen.Group) {
					if b.ep.Context != nil {
						g.Id("ctx")
					}
					for _, p := range b.ep.Params {
						g.Id("r").Dot(p.FldName)
					}
				})
				if b.ep.Error != nil {
					g.Do(gen.CheckErr(jen.Return(jen.Nil(), jen.Err())))
				}
				g.ReturnFunc(func(g *jen.Group) {
					if len(b.ep.Results) > 0 {
						g.Op("&").Id(respName).ValuesFunc(func(g *jen.Group) {
							for _, p := range b.ep.Results {
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
		),
	)
}

func (b *serverEndpointBuilder) BuildReqStruct() ServerEndpointBuilder {
	if len(b.ep.Params) > 0 {
		var stRequests []jen.Code

		reqName := b.reqStructName()

		if b.ep.ReqRootXMLName != "" {
			stRequests = append(stRequests, jen.Id("XMLName").Qual("encoding/xml", "Name").Tag(map[string]string{"xml": b.ep.ReqRootXMLName}))
		}
		for _, p := range b.ep.Params {
			st := jen.Id(p.FldName).Add(types.Convert(p.Type, b.qualifier.Qual))
			if p.Name != "" && p.HTTPType == "body" {
				st.Tag(map[string]string{"json": p.Name})
			} else {
				st.Tag(map[string]string{"json": "-"})
			}
			stRequests = append(stRequests, st)
		}
		b.codes = append(b.codes,
			jen.Type().Id(reqName).Struct(stRequests...),
		)
	}
	return b
}

func (b *serverEndpointBuilder) BuildRespStruct() ServerEndpointBuilder {
	if len(b.ep.Results) > 0 {
		var stResponses []jen.Code

		respName := b.respStructName()

		for _, p := range b.ep.Results {
			st := jen.Id(p.FldNameExport).Add(types.Convert(p.Type, b.qualifier.Qual))
			if p.Name != "" && p.HTTPType == "body" {
				st.Tag(map[string]string{"json": p.Name})
			} else {
				st.Tag(map[string]string{"json": "-"})
			}
			stResponses = append(stResponses, st)
		}
		b.codes = append(b.codes,
			jen.Type().Id(respName).Struct(stResponses...),
		)

	}
	return b
}

func (b *serverEndpointBuilder) BuildReqDec() ServerEndpointBuilder {
	if len(b.ep.Params) > 0 {
		reqDecName := b.reqDecFuncName()
		reqName := b.reqStructName()

		blockFunc := func(g *jen.Group) {
			g.Var().Id("param").Id(reqName)

			if len(b.ep.BodyParams) > 0 {
				switch b.ep.HTTPMethod {
				case "POST", "PUT", "PATCH", "DELETE":
					nameContentType, typ := b.handlerStrategy.HeaderParam("content-type")

					g.Add(typ)

					g.Id("parts").Op(":=").Qual("strings", "Split").Call(jen.Id(nameContentType), jen.Lit(";"))
					g.If(jen.Len(jen.Id("parts")).Op("==").Lit(0)).Block(
						jen.Return(jen.Nil(), jen.Err()),
					)
					g.Id(nameContentType).Op("=").Id("parts").Index(jen.Lit(0))

					g.Var().Id("bodyData").Op("=").Make(jen.Index().Byte(), jen.Lit(0), jen.Lit(10485760)) // 10MB
					g.Id("buf").Op(":=").Qual("bytes", "NewBuffer").Call(jen.Id("bodyData"))
					g.List(jen.Id("written"), jen.Id("err")).Op(":=").Qual("io", "Copy").Call(jen.Id("buf"), b.handlerStrategy.BodyPathParam())
					g.Do(gen.CheckErr(
						jen.Return(),
					))
					g.Switch(jen.Id(nameContentType)).BlockFunc(func(g *jen.Group) {
						g.Default().Block(
							jen.Return(jen.Nil(), jen.Op("&").Id("contentTypeInvalidError").Values()),
						)
						for _, contentType := range b.ep.ContentTypes {
							bodyParams := jen.Id("param")
							if len(b.ep.BodyParams) == 1 && b.ep.NoWrapRequest {
								bodyParams = bodyParams.Dot(b.ep.BodyParams[0].FldName)
							}

							switch contentType {
							case "xml":
								g.Case(jen.Lit("application/xml")).BlockFunc(func(g *jen.Group) {
									g.Err().Op("=").Qual("encoding/xml", "Unmarshal").Call(jen.Id("bodyData").Index(jen.Op(":").Id("written")), jen.Op("&").Add(bodyParams))
									g.Do(gen.CheckErr(
										jen.Return(jen.Nil(), jen.Err()),
									))
								})
							case "json":
								g.Case(jen.Lit("application/json")).BlockFunc(func(g *jen.Group) {
									g.Err().Op("=").Qual("encoding/json", "Unmarshal").Call(jen.Id("bodyData").Index(jen.Op(":").Id("written")), jen.Op("&").Add(bodyParams))
									g.Do(gen.CheckErr(
										jen.Return(jen.Nil(), jen.Err()),
									))
								})
							case "urlencoded":
								g.Case(jen.Lit("application/x-www-form-urlencoded")).BlockFunc(func(g *jen.Group) {
									g.Add(b.handlerStrategy.FormParams())

									for _, p := range b.ep.BodyParams {
										formParamName, typ := b.handlerStrategy.FormParam(p.Name)

										g.Add(typ)

										g.Add(gen.ParseValue(jen.Id(formParamName), jen.Id("param").Dot(p.FldName), "=", p.Type, b.qualifier.Qual, func() jen.Code {
											return jen.Do(gen.CheckErr(
												jen.Return(jen.Nil(), jen.Err()),
											))
										}))
									}
								})
							case "multipart":
								g.Case(jen.Lit("multipart/form-data")).BlockFunc(func(g *jen.Group) {
									g.Add(b.handlerStrategy.MultipartFormParams(b.ep.MultipartMaxMemory))
									for _, p := range b.ep.BodyParams {
										formParamName, typ := b.handlerStrategy.MultipartFormParam(p.Name)

										g.Add(typ)

										g.Add(gen.ParseValue(jen.Id(formParamName), jen.Id("param").Dot(p.FldName), "=", p.Type, b.qualifier.Qual, func() jen.Code {
											return jen.Do(gen.CheckErr(
												jen.Return(jen.Nil(), jen.Err()),
											))
										}))
									}
								})
							}
						}
					})
				}
			}

			if b.handlerStrategy.UsePathParams() && len(b.ep.PathParams) > 0 {
				for _, p := range b.ep.PathParams {
					pathParamName, typ := b.handlerStrategy.PathParam(p.Name)

					g.Add(typ)

					g.If(jen.Id(pathParamName).Op("!=").Lit("")).Block(
						jen.Add(gen.ParseValue(jen.Id(pathParamName), jen.Id("param").Dot(p.FldName), "=", p.Type, b.qualifier.Qual, func() jen.Code {
							return jen.Do(gen.CheckErr(
								jen.Return(jen.Nil(), jen.Err()),
							))
						})),
						jen.Do(gen.CheckErr(jen.Return())),
					)
				}
			}

			if len(b.ep.QueryParams) > 0 {
				g.Add(b.handlerStrategy.QueryParams())
				for _, param := range b.ep.QueryParams {
					queryParamName, typ := b.handlerStrategy.QueryParam(param.Name)

					g.Add(typ)

					paramID := jen.Id("param").Dot(param.FldName)
					if param.Parent != nil {
						paramID = jen.Id("param").Dot(param.Parent.FldName).Dot(param.FldName)
					}

					g.If(jen.Id(queryParamName).Op("!=").Lit("")).Block(
						jen.Add(gen.ParseValue(jen.Id(queryParamName), paramID, "=", param.Type, b.qualifier.Qual, func() jen.Code {
							return jen.Do(gen.CheckErr(jen.Return()))
						})),
					)
				}
			}

			if len(b.ep.HeaderParams) > 0 {
				for _, param := range b.ep.HeaderParams {
					queryParamName, typ := b.handlerStrategy.HeaderParam(param.Name)

					g.Add(typ)

					g.If(jen.Id(queryParamName).Op("!=").Lit("")).Block(
						jen.Add(gen.ParseValue(jen.Id(queryParamName), jen.Id("param").Dot(param.FldName), "=", param.Type, b.qualifier.Qual, func() jen.Code {
							return jen.Do(gen.CheckErr(jen.Return()))
						})),
					)
				}
			}
			if len(b.ep.Params) > 0 {
				g.Return(jen.Op("&").Id("param"), jen.Nil())
			} else {
				g.Return()
			}
		}

		b.codes = append(b.codes,
			jen.Func().Id(reqDecName).ParamsFunc(func(g *jen.Group) {
				g.Id(b.handlerStrategy.ReqArgName()).Add(b.handlerStrategy.ReqType())
			}).Params(
				jen.Id("result").Any(),
				jen.Err().Error(),
			).BlockFunc(blockFunc),
		)
	}
	return b
}

func (b *serverEndpointBuilder) BuildRespEnc() ServerEndpointBuilder {
	respEncName := b.respEncFuncName()
	respName := b.respStructName()

	if len(b.ep.Results) > 0 {
		blockFunc := func(g *jen.Group) {
			if len(b.ep.BodyResults) > 0 {
				if !b.ep.NoWrapResponse {
					g.Var().Id("wrapResult").StructFunc(gen.WrapResponse(b.ep.WrapResponse, b.ep.Results, b.qualifier.Qual))
					for _, r := range b.ep.Results {
						g.Id("wrapResult").Do(func(s *jen.Statement) {
							for _, name := range b.ep.WrapResponse {
								s.Dot(strcase.ToCamel(name))
							}
						}).Dot(r.FldNameExport).Op("=").Id("result").Assert(jen.Op("*").Id(respName)).Dot(r.FldNameExport)
					}

					g.Id("result").Op("=").Id("wrapResult")

				} else if len(b.ep.BodyResults) == 1 {
					g.Id("result").Op("=").Id("result").Assert(jen.Op("*").Id(respName)).Dot(b.ep.BodyResults[0].FldNameExport)
				}
				g.Return(jen.Id("result"), jen.Nil())
			} else {
				g.Return(jen.Nil(), jen.Nil())
			}
		}

		b.codes = append(b.codes,
			jen.Func().Id(respEncName).Params(
				jen.Id("result").Any(),
			).Params(
				jen.Any(),
				jen.Error(),
			).BlockFunc(blockFunc))
	}
	return b
}

func (b *serverEndpointBuilder) endpointFuncName() string {
	return endpointFuncName(b.iface, b.ep)
}

func (b *serverEndpointBuilder) reqStructName() string {
	return reqStructName(b.iface, b.ep)
}

func (b *serverEndpointBuilder) respStructName() string {
	return respStructName(b.iface, b.ep)
}

func (b *serverEndpointBuilder) reqDecFuncName() string {
	return reqDecFuncName(b.iface, b.ep)
}

func (b *serverEndpointBuilder) respEncFuncName() string {
	return respEncFuncName(b.iface, b.ep)
}
