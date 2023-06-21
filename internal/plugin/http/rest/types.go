package rest

import (
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gen"
	. "github.com/dave/jennifer/jen"
)

var (
	echoMiddlewareFunc = Qual("github.com/labstack/echo/v4", "MiddlewareFunc")
	epFunc             = Func().Params(
		Id("ctx").Qual("context", "Context"),
		Id("request").Any(),
	).Params(
		Any(),
		Error(),
	)
	respEncFunc = Func().Params(
		Id("result").Any(),
	).Params(
		Any(),
		Error(),
	)

	reqDecFunc = Func().Params(
		Id("pathParams").Id("pathParams"),
		Id("request").Op("*").Qual("net/http", "Request"),
		Id("params").Qual("encoding/json", "RawMessage"),
	).Params(
		Id("result").Any(),
		Err().Error(),
	)

	serverErrorEncoder = gen.CheckErr(
		Id("serverErrorEncoder").Call(
			Id("rw"),
			Err(),
		),
		Return(),
	)
)

func GenTypes() func(f *file.GoFile) {
	return func(f *file.GoFile) {

		f.Type().Id("contentTypeInvalidError").Struct()

		f.Func().Params(Op("*").Id("contentTypeInvalidError")).Id("Error").Params().String().Block(Return(Lit("content type invalid")))
		f.Func().Params(Op("*").Id("contentTypeInvalidError")).Id("StatusCode").Params().Int().Block(Return(Lit(400)))

		f.Type().Id("pathParams").Interface(
			Id("Param").Params(String()).String(),
		)
		f.Type().Id("pathParamsNoop").Struct()
		f.Func().Params(Op("*").Id("pathParamsNoop")).Id("Param").Params(String()).String().Block(Return(Lit("")))

		// optionsName := "Options"
		// optionName := "Option"

		// f.Type().Id(optionName).Func().Params(Op("*").Id(optionsName))
		// f.Type().Id(optionsName).Struct(
		// 	Id("before").Index().Id("BeforeFunc"),
		// 	Id("after").Index().Id("AfterFunc"),
		// 	Id("middleware").Index().Id("MiddlewareFunc"),
		// ).Line()

		// f.Type().Id("BeforeFunc").Func().Params(
		// 	Qual("context", "Context"),
		// 	Op("*").Qual("net/http", "Request"),
		// ).Qual("context", "Context")
		// f.Type().Id("AfterFunc").Func().Params(
		// 	Qual("context", "Context"),
		// 	Qual("net/http", "ResponseWriter"),
		// ).Qual("context", "Context")
		// f.Type().Id("MiddlewareFunc").Func().Params(
		// 	Qual("net/http", "Handler"),
		// ).Qual("net/http", "Handler")
		// f.Type().Id("PopulatePathParamsFunc").Func().Params().Qual("net/url", "Values")

		// f.Func().Id("applyOptions").Params(Id("opts").Op("...").Id(optionName)).Id(optionName).Block(
		// 	Return(
		// 		Func().Params(Id("o").Op("*").Id(optionsName)).Block(
		// 			For(List(Id("_"), Id("opt")).Op(":=")).Range().Id("opts").Block(
		// 				Id("opt").Call(Id("o")),
		// 			),
		// 		),
		// 	),
		// )

		// f.Func().Id("applyHandlerOptions").Params(Id("opts").Index().Id(optionName)).Id("MiddlewareFunc").Block(
		// 	Id("o").Op(":=").Op("&").Id(optionsName).Values(),
		// 	Id("applyOptions").Call(Id("opts").Op("...")).Call(Id("o")),
		// 	Return(
		// 		Func().Params(Id("next").Qual("net/http", "Handler")).Qual("net/http", "Handler").Block(
		// 			Return(
		// 				Qual("net/http", "HandlerFunc").Call(
		// 					Func().Params(
		// 						Id("rw").Qual("net/http", "ResponseWriter"),
		// 						Id("r").Op("*").Qual("net/http", "Request"),
		// 					).Block(
		// 						Id("ctx").Op(":=").Id("r").Dot("Context").Call(),

		// 						For(List(Id("_"), Id("beforeFunc")).Op(":=")).Range().Id("o").Dot("before").Block(
		// 							Id("ctx").Op("=").Id("beforeFunc").Call(Id("ctx"), Id("r")),
		// 						),
		// 						Id("applyMiddleware").Call(Id("next"), Id("o").Dot("middleware").Op("...")).
		// 							Dot("ServeHTTP").Call(
		// 							Id("rw"),
		// 							Id("r").Dot("WithContext").Call(Id("ctx")),
		// 						),
		// 						For(List(Id("_"), Id("afterFunc")).Op(":=")).Range().Id("o").Dot("after").Block(
		// 							Id("afterFunc").Call(Id("ctx"), Id("rw")),
		// 						),
		// 					),
		// 				),
		// 			),
		// 		),
		// 	),
		// )

		// f.Func().Id("Before").Params(Id("before").Op("...").Id("BeforeFunc")).Id(optionName).Block(
		// 	Return(
		// 		Func().Params(Id("o").Op("*").Id(optionsName)).Block(
		// 			Id("o").Dot("before").Op("=").Append(Id("o").Dot("before"), Id("before").Op("...")),
		// 		),
		// 	),
		// )
		// f.Func().Id("After").Params(Id("after").Op("...").Id("AfterFunc")).Id(optionName).Block(
		// 	Return(
		// 		Func().Params(Id("o").Op("*").Id(optionsName)).Block(
		// 			Id("o").Dot("after").Op("=").Append(Id("o").Dot("after"), Id("after").Op("...")),
		// 		),
		// 	),
		// )
		// f.Func().Id("Middleware").Params(Id("middleware").Op("...").Id("MiddlewareFunc")).Id(optionName).Block(
		// 	Return(
		// 		Func().Params(Id("o").Op("*").Id(optionsName)).Block(
		// 			Id("o").Dot("middleware").Op("=").Append(Id("o").Dot("middleware"), Id("middleware").Op("...")),
		// 		),
		// 	),
		// )

		// f.Func().Id("applyMiddleware").Params(
		// 	Id("h").Qual("net/http", "Handler"),
		// 	Id("middleware").Op("...").Id("MiddlewareFunc"),
		// ).Qual("net/http", "Handler").Block(
		// 	For(
		// 		Id("i").Op(":=").Len(Id("middleware")).Op("-").Lit(1),
		// 		Id("i").Op(">=").Lit(0),
		// 		Id("i").Op("--"),
		// 	).Block(
		// 		Id("h").Op("=").Id("middleware").Index(Id("i")).Call(Id("h")),
		// 	),
		// 	Return(Id("h")),
		// )

		// f.Func().Id("pathParamsFromEchoContext").Params(
		// 	Id("c").Qual("github.com/labstack/echo/v4", "Context"),
		// ).Qual("net/url", "Values").Block(
		// 	Id("paramNames").Op(":=").Id("c").Dot("ParamNames").Call(),
		// 	Id("paramValues").Op(":=").Id("c").Dot("ParamValues").Call(),
		// 	Id("values").Op(":=").Qual("net/url", "Values").Values(),
		// 	For(
		// 		List(Id("i"), Id("name")).Op(":=").Range().Id("paramNames"),
		// 	).Block(

		// 		Id("values").Dot("Set").Call(Id("name"), Id("paramValues").Index(Id("i"))),
		// 	),
		// 	Return(Id("values")),
		// )
	}
}
