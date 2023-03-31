package generic

import (
	"github.com/555f/gg/pkg/file"

	. "github.com/dave/jennifer/jen"
)

func GenServer() func(f *file.GoFile) {
	return func(f *file.GoFile) {
		f.Type().Id("statusCoder").Interface(
			Id("StatusCode").Params().Int(),
		)
		f.Type().Id("dataer").Interface(
			Id("Data").Params().Qual("bytes", "Buffer"),
		)
		f.Type().Id("contenter").Interface(
			Id("ContentType").Params().String(),
		)
		f.Type().Id("headerer").Interface(
			Id("Headers").Params().Qual("net/http", "Header"),
		)
		f.Type().Id("cookier").Interface(
			Id("Cookies").Params().Index().Qual("net/http", "Cookie"),
		)

		f.Type().Id("serverOptions").Struct(
			Id("before").Index().Id("ServerBeforeFunc"),
			Id("after").Index().Id("ServerAfterFunc"),
			Id("serverMiddleware").Index().Id("ServerMiddlewareFunc"),
			Id("middleware").Index().Id("EndpointMiddleware"),
		)
		f.Type().Id("ServerErrorEncoder").Func().
			Params(
				Id("ctx").Qual("context", "Context"),
				Err().Error(),
				Id("w").Qual("net/http", "ResponseWriter"),
			)
		f.Type().Id("ServerBeforeFunc").Func().Params(
			Qual("context", "Context"),
			Op("*").Qual("net/http", "Request"),
		).Qual("context", "Context")
		f.Type().Id("ServerAfterFunc").Func().Params(
			Qual("context", "Context"),
			Qual("net/http", "ResponseWriter"),
		).Qual("context", "Context")
		f.Type().Id("ServerMiddlewareFunc").Func().Params(
			Qual("net/http", "Handler"),
		).Qual("net/http", "Handler")

		f.Type().Id("ServerOption").Func().Params(Op("*").Id("Server"))

		f.Type().Id("Endpoint").Op("=").Func().Params(
			Id("ctx").Qual("context", "Context"),
			Id("request").Interface(),
		).Params(Id("response").Interface(), Err().Error())

		f.Type().Id("EndpointMiddleware").Op("=").Func().Params(Id("Endpoint")).Id("Endpoint")

		f.Func().Id("middlewareChain").Params(Id("middlewares").Index().Id("EndpointMiddleware")).Id("EndpointMiddleware").Block(
			Return(
				Func().Params(Id("next").Id("Endpoint")).Id("Endpoint").Block(
					If(Len(Id("middlewares")).Op("==").Lit(0)).Block(
						Return(Id("next")),
					),
					Id("outer").Op(":=").Id("middlewares").Index(Lit(0)),
					Id("others").Op(":=").Id("middlewares").Index(Lit(1).Op(":")),
					For(Id("i").Op(":=").Len(Id("others")).Op("-").Lit(1), Id("i").Op(">=").Lit(0), Id("i").Op("--")).Block(
						Id("next").Op("=").Id("others").Index(Id("i")).Call(Id("next")),
					),
					Return(Id("outer").Call(Id("next"))),
				),
			),
		)
		f.Func().Id("serverMiddlewareChain").Params(Id("serverMiddlewares").Index().Id("ServerMiddlewareFunc")).Id("ServerMiddlewareFunc").Block(
			Return(
				Func().Params(Id("next").Qual("net/http", "Handler")).Qual("net/http", "Handler").Block(
					If(Len(Id("serverMiddlewares")).Op("==").Lit(0)).Block(
						Return(Id("next")),
					),
					Id("outer").Op(":=").Id("serverMiddlewares").Index(Lit(0)),
					Id("others").Op(":=").Id("serverMiddlewares").Index(Lit(1).Op(":")),
					For(Id("i").Op(":=").Len(Id("others")).Op("-").Lit(1), Id("i").Op(">=").Lit(0), Id("i").Op("--")).Block(
						Id("next").Op("=").Id("others").Index(Id("i")).Call(Id("next")),
					),
					Return(Id("outer").Call(Id("next"))),
				),
			),
		)
	}
}
