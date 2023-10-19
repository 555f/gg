package gen

import (
	"github.com/555f/gg/internal/plugin/jsonrpc/options"
	"github.com/555f/gg/pkg/gen"
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
		for _, p := range b.ep.Params {
			st := jen.Id(p.FldName).Add(types.Convert(p.Type, b.qualifier.Qual))
			if p.Name != "" {
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
			if p.Name != "" {
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

			bodyParams := jen.Id("param")
			if len(b.ep.Params) == 1 {
				bodyParams = bodyParams.Dot(b.ep.Params[0].FldName)
			}
			g.Err().Op("=").Qual("encoding/json", "Unmarshal").Call(jen.Id("params"), jen.Op("&").Add(bodyParams))
			g.Do(gen.CheckErr(
				jen.Return(jen.Nil(), jen.Err()),
			))

			g.Return(jen.Op("&").Id("param"), jen.Nil())
		}

		b.codes = append(b.codes,
			jen.Func().Id(reqDecName).ParamsFunc(func(g *jen.Group) {
				g.Id("params").Qual(jsonPkg, "RawMessage")
			}).Params(
				jen.Id("req").Any(),
				jen.Err().Error(),
			).BlockFunc(blockFunc),
		)
	}

	return b
}

// func (b *serverEndpointBuilder) BuildRespEnc() ServerEndpointBuilder {
// 	respEncName := b.respEncFuncName()
// 	respName := b.respStructName()

// 	if len(b.ep.Results) > 0 {
// 		blockFunc := func(g *jen.Group) {
// 			if len(b.ep.Results) > 0 {
// 				if len(b.ep.Results) == 1 {
// 					g.Id("result").Op("=").Id("result").Assert(jen.Op("*").Id(respName)).Dot(b.ep.Results[0].FldNameExport)
// 				} else {
// 					g.Return(jen.Id("result"), jen.Nil())
// 				}

// 			} else {
// 				g.Return(jen.Nil(), jen.Nil())
// 			}
// 		}

// 		b.codes = append(b.codes,
// 			jen.Func().Id(respEncName).Params(
// 				jen.Id("result").Any(),
// 			).Params(
// 				jen.Any(),
// 				jen.Error(),
// 			).BlockFunc(blockFunc))
// 	}
// 	return b
// }

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

// func (b *serverEndpointBuilder) respEncFuncName() string {
// 	return respEncFuncName(b.iface, b.ep)
// }
