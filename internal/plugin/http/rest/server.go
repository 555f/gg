package rest

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/strcase"

	. "github.com/dave/jennifer/jen"
)

func GenServer(errorWrapper *options.ErrorWrapper) func(f *file.GoFile) {
	return func(f *file.GoFile) {
		f.Func().
			Id("encodeJSONResponse").
			Params(
				Id("ctx").Qual("context", "Context"),
				Id("w").Qual("net/http", "ResponseWriter"),
				Id("response").Any(),
			).
			Block(
				Id("statusCode").Op(":=").Lit(200),
				Var().Id("data").Qual("bytes", "Buffer"),
				If(Id("response").Op("!=").Nil()).Block(
					If(List(Id("v"), Id("ok")).Op(":=").Id("response").Assert(Id("statusCoder")), Id("ok")).Block(
						Id("statusCode").Op("=").Id("v").Dot("StatusCode").Call(),
					),
					If(List(Id("v"), Id("ok")).Op(":=").Id("response").Assert(Id("contenter")), Id("ok")).Block(
						Id("w").Dot("Header").Call().Dot("Set").Call(Lit("Content-Type"), Id("v").Dot("ContentType").Call()),
					).Else().Block(
						Id("w").Dot("Header").Call().Dot("Set").Call(Lit("Content-Type"), Lit("application/json; charset=utf-8")),
					),

					If(List(Id("v"), Id("ok")).Op(":=").Id("response").Assert(Id("dataer")), Id("ok")).Block(
						Id("data").Op("=").Id("v").Dot("Data").Call(),
					).Else().Block(
						If(Err().Op(":=").Qual("encoding/json", "NewEncoder").Call(Op("&").Id("data")).Dot("Encode").Call(Id("response")), Err().Op("!=").Nil()).Block(
							Return(),
						),
					),
					If(List(Id("v"), Id("ok")).Op(":=").Id("response").Assert(Id("headerer")), Id("ok")).Block(
						For(List(Id("key"), Id("values"))).Op(":=").Range().Id("v").Dot("Headers").Call().Block(
							For(List(Id("_"), Id("val"))).Op(":=").Range().Id("values").Block(
								Id("w").Dot("Header").Call().Dot("Add").Call(Id("key"), Id("val")),
							),
						),
					),
					If(List(Id("v"), Id("ok")).Op(":=").Id("response").Assert(Id("cookier")), Id("ok")).Block(
						For(List(Id("_"), Id("c"))).Op(":=").Range().Id("v").Dot("Cookies").Call().Block(
							Qual("net/http", "SetCookie").Call(Id("w"), Op("&").Id("c")),
						),
					),
				).Else().Block(
					Id("statusCode").Op("=").Lit(204),
				),
				Id("w").Dot("WriteHeader").Call(Id("statusCode")),
				If(List(Id("_"), Err()).Op(":=").Id("w").Dot("Write").Call(Id("data").Dot("Bytes").Call()), Err().Op("!=").Nil()).Block(
					Panic(Err()),
				),
			)

		f.Func().Id("serverErrorEncoder").Params(
			Id("ctx").Qual("context", "Context"),
			Id("w").Qual("net/http", "ResponseWriter"),
			Err().Error(),
		).BlockFunc(func(g *Group) {

			g.Var().Id("statusCode").Int().Op("=").Qual("net/http", "StatusInternalServerError")
			g.Id("h").Op(":=").Id("w").Dot("Header").Call()

			g.If(List(Id("e"), Id("ok")).Op(":=").Err().Assert(Interface(Id("StatusCode").Params().Int())), Id("ok")).Block(
				Id("statusCode").Op("=").Id("e").Dot("StatusCode").Call(),
			)
			g.If(List(Id("headerer"), Id("ok")).Op(":=").Err().Assert(Interface(Id("Headers").Params().Qual("net/http", "Header"))), Id("ok")).Block(
				For(List(Id("k"), Id("values"))).Op(":=").Range().Id("headerer").Dot("Headers").Call().Block(
					For(List(Id("_"), Id("v"))).Op(":=").Range().Id("values").Block(
						Id("h").Dot("Add").Call(Id("k"), Id("v")),
					),
				),
			)

			if errorWrapper != nil {
				errorWrapperName := strcase.ToLowerCamel(errorWrapper.Struct.Named.Name)
				g.Id(errorWrapperName).Op(":=").Do(f.Import(errorWrapper.Struct.Named.Pkg.Path, errorWrapper.Struct.Named.Name)).Values()
				for _, field := range errorWrapper.Fields {
					g.If(List(Id("e"), Id("ok")).Op(":=").Err().Assert(Interface(Id(field.Interface))), Id("ok")).Block(
						Id(errorWrapperName).Dot(field.FldName).Op("=").Id("e").Op(".").Id(field.MethodName).Call(),
					)
				}
				g.List(Id("data"), Id("jsonErr")).Op(":=").Qual("encoding/json", "Marshal").Call(Id(errorWrapperName))
			}

			g.If(Id("jsonErr").Op("!=").Nil()).Block(
				List(Id("_"), Id("_")).Op("=").Id("w").Dot("Write").Call(Index().Byte().Call(Lit("unexpected marshal error"))),
				Return(),
			)
			g.Id("h").Dot("Set").Call(Lit("Content-Type"), Lit("application/json; charset=utf-8"))
			g.Id("w").Dot("WriteHeader").Call(Id("statusCode"))
			g.List(Id("_"), Id("_")).Op("=").Id("w").Dot("Write").Call(Id("data"))
		})

		f.Type().Id("Regex").Struct(
			Id("Pattern").String(),
			Id("Params").Index().String(),
		)

		f.Type().Id("route").Struct(
			Id("method").String(),
			Id("regex").Op("*").Qual("regexp", "Regexp"),
			Id("pattern").String(),
			Id("params").Index().String(),
			Id("handler").Qual("net/http", "Handler"),
			Id("before").Index().Id("ServerBeforeFunc"),
			Id("after").Index().Id("ServerAfterFunc"),
		)
		f.Type().Id("Server").StructFunc(func(g *Group) {
			g.Id("staticRoutes").Map(String()).Op("*").Id("route")
			g.Id("regexRoutes").Map(String()).Index().Op("*").Id("route")
		})

		f.Type().Id("RouteOption").Func().Params(Id("r").Op("*").Id("route"))
		f.Func().Id("After").Params(Id("after").Op("...").Id("ServerAfterFunc")).Id("RouteOption").Block(
			Return(
				Func().Params(Id("r").Op("*").Id("route")).Block(
					Id("r").Dot("after").Op("=").Append(Id("r").Dot("after"), Id("after").Op("...")),
				),
			),
		)
		f.Func().Id("Before").Params(Id("before").Op("...").Id("ServerBeforeFunc")).Id("RouteOption").Block(
			Return(
				Func().Params(Id("r").Op("*").Id("route")).Block(
					Id("r").Dot("before").Op("=").Append(Id("r").Dot("before"), Id("before").Op("...")),
				),
			),
		)
		f.Func().Id("Middleware").Params(Id("middleware").Op("...").Id("ServerMiddlewareFunc")).Id("RouteOption").Block(
			Return(
				Func().Params(Id("r").Op("*").Id("route")).Block(
					Id("r").Dot("handler").Op("=").Id("serverMiddlewareChain").Call(Id("middleware")).Call(Id("r").Dot("handler")),
				),
			),
		)

		f.Func().Params(Id("s").Op("*").Id("Server")).Id("AddRoute").Params(
			Id("method").String(),
			Id("pattern").Any(),
			Id("handler").Qual("net/http", "Handler"),
			Id("opts").Op("...").Id("RouteOption"),
		).Block(
			Id("r").Op(":=").Op("&").Id("route").Values(
				Id("method").Op(":").Id("method"),
				Id("handler").Op(":").Id("handler"),
			),
			For(List(Id("_"), Id("opt")).Op(":=").Range().Id("opts")).Block(
				Id("opt").Call(Id("r")),
			),
			Switch(Id("t").Op(":=").Id("pattern").Assert(Id("type"))).Block(
				Default().Block(Panic(Lit("pattern must be string or Regex type"))),
				Case(Id("string")).Block(
					Id("s").Dot("staticRoutes").Index(Id("method").Op("+").Id("t")).Op("=").Id("r"),
				),
				Case(Id("Regex")).Block(
					Id("r").Dot("regex").Op("=").Qual("regexp", "MustCompile").Call(Id("t").Dot("Pattern")),
					Id("r").Dot("params").Op("=").Id("t").Dot("Params"),
					Id("s").Dot("regexRoutes").Index(Id("method")).Op("=").Append(Id("s").Dot("regexRoutes").Index(Id("method")), Id("r")),
				),
			),
		)

		f.Func().Params(Id("s").Op("*").Id("Server")).Id("handleRoute").Params(
			Id("route").Op("*").Id("route"),
			Id("rw").Qual("net/http", "ResponseWriter"),
			Id("r").Op("*").Qual("net/http", "Request"),
		).Block(
			For(List(Id("_"), Id("before")).Op(":=").Range().Id("route").Dot("before")).Block(
				Id("r").Op("=").Id("r").Dot("WithContext").Call(Id("before").Call(Id("r").Dot("Context").Call(), Id("r"))),
			),
			Id("route").Dot("handler").Dot("ServeHTTP").Call(Id("rw"), Id("r")),
			For(List(Id("_"), Id("after")).Op(":=").Range().Id("route").Dot("after")).Block(
				Id("r").Op("=").Id("r").Dot("WithContext").Call(Id("after").Call(Id("r").Dot("Context").Call(), Id("rw"))),
			),
		)

		f.Func().Params(Id("s").Op("*").Id("Server")).Id("ServeHTTP").Params(
			Id("w").Qual("net/http", "ResponseWriter"),
			Id("r").Op("*").Qual("net/http", "Request"),
		).Block(
			Id("requestPath").Op(":=").Id("r").Dot("URL").Dot("Path"),
			If(List(Id("route"), Id("ok")).Op(":=").Id("s").Dot("staticRoutes").Index(Id("r").Dot("Method").Op("+").Id("requestPath")), Id("ok")).Block(
				Id("s").Dot("handleRoute").Call(Id("route"), Id("w"), Id("r")),
				Return(),
			),
			If(List(Id("routes"), Id("ok")).Op(":=").Id("s").Dot("regexRoutes").Index(Id("r").Dot("Method")), Id("ok")).Block(
				For(List(Id("_"), Id("route"))).Op(":=").Range().Id("routes").Block(
					If(Op("!").Id("route").Dot("regex").Dot("MatchString").Call(Id("requestPath"))).Block(Continue()),
					Id("matches").Op(":=").Id("route").Dot("regex").Dot("FindStringSubmatch").Call(Id("requestPath")),
					If(Len(Id("matches").Index(Lit(0))).Op("!=").Len(Id("requestPath"))).Block(Continue()),
					If(Len(Id("route").Dot("params")).Op(">").Lit(0)).Block(
						Id("values").Op(":=").Id("r").Dot("URL").Dot("Query").Call(),
						For(List(Id("i"), Id("match"))).Op(":=").Range().Id("matches").Index(Lit(1).Op(":")).Block(
							Id("values").Dot("Add").Call(Id("route").Dot("params").Index(Id("i")), Id("match")),
						),
						Id("r").Dot("URL").Dot("RawQuery").Op("=").Id("values").Dot("Encode").Call(),
					),
					Id("s").Dot("handleRoute").Call(Id("route"), Id("w"), Id("r")),
					Return(),
				),
			),
			Id("w").Dot("WriteHeader").Call(Qual("net/http", "StatusMethodNotAllowed")),
		)

		f.Func().Id("NewRESTServer").Params(Id("opts").Op("...").Id("ServerOption")).Op("*").Id("Server").Block(
			Id("s").Op(":=").Op("&").Id("Server").Values(
				Id("staticRoutes").Op(":").Make(Map(String()).Op("*").Id("route")),
				Id("regexRoutes").Op(":").Make(Map(String()).Index().Op("*").Id("route")),
			),
			For(List(Id("_"), Id("opt"))).Op(":=").Range().Id("opts").Block(
				Id("opt").Call(Id("s")),
			),
			Return(Id("s")),
		)
	}
}
