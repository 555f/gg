package rest

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/gen"
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
				// epName := endpointFuncName(b.iface, ep)
				reqDecName := reqDecFuncName(b.iface, ep)
				// respEncName := respEncFuncName(b.iface, ep)

				// hasResults := len(ep.Results) > 0
				hasParams := len(ep.Params) > 0
				// epResultOp := ":="
				// epResultName := "resp"

				// if hasParams && !hasResults {
				// epResultOp = "="
				// }
				// if !hasResults {
				// epResultName = "_"
				// }

				var resultID jen.Code

				// errorEncoder := jen.Id("o").Dot("errorEncoder").Call(
				// jen.Id(b.handlerStrategy.RespArgName()),
				// jen.Err(),
				// )

				handlerFuncBody := jen.CustomFunc(jen.Options{Multi: true}, func(g *jen.Group) {
					g.Id("reqCtx").Op(":=").Qual(contextPkg, "TODO").Call()

					if hasParams {
						g.List(jen.Id("req"), jen.Err()).Op(":=").Id(reqDecName).Call(jen.Id(b.handlerStrategy.ReqArgName()))
						g.Do(gen.CheckErr(jen.Return()))
					}
					// g.List(jen.Id(epResultName), jen.Err()).Op(epResultOp).Id(epName).Call(jen.Id("svc")).CallFunc(func(g *jen.Group) {
					// 	g.Id("reqCtx")
					// 	if hasParams {
					// 		g.Id("req")
					// 	} else {
					// 		g.Nil()
					// 	}
					// })
					// g.Do(gen.CheckErr(
					// errorEncoder,
					// jen.Return(),
					// ))
					// if hasResults {
					// resultID = jen.Id("result")
					// g.List(resultID, jen.Err()).Op(":=").Id(respEncName).Call(jen.Id("resp"))
					// g.Do(gen.CheckErr(
					// errorEncoder,
					// jen.Return(),
					// ))
					// }
				})

				g.Add(b.handlerStrategy.HandlerFunc(
					ep.HTTPMethod,
					ep.Path,
					resultID,
					jen.Append(jen.Id("o").Dot("middleware"), jen.Id("o").Dot("middleware"+ep.MethodName).Op("...")),
					handlerFuncBody,
				))
			}
		}))
	return b
}

func (b *serverControllerBuilder) Endpoint(ep options.Endpoint) ServerEndpointBuilder {
	return &serverEndpointBuilder{serverControllerBuilder: b, handlerStrategy: b.handlerStrategy, iface: b.iface, ep: ep, qualifier: b.qualifier}
}
