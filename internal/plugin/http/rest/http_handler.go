package rest

import (
	"github.com/555f/gg/pkg/file"
	. "github.com/dave/jennifer/jen"
)

func GenHTTPHandler() func(f *file.GoFile) {
	return func(f *file.GoFile) {
		f.Func().Id("httpHandler").Params(
			Id("ep").Add(epFunc),
			Id("reqDec").Add(reqDecFunc),
			Id("respEnc").Add(respEncFunc),
			Id("pathParams").Id("pathParams"),
		).Qual("net/http", "HandlerFunc").Block(
			Return(
				Func().Params(
					Id("rw").Qual("net/http", "ResponseWriter"),
					Id("r").Op("*").Qual("net/http", "Request"),
				).BlockFunc(func(g *Group) {
					g.Var().Id("params").Any()
					g.Var().Id("result").Any()
					g.Var().Err().Error()

					g.If(Id("reqDec").Op("!=").Nil()).Block(
						Var().Id("wb").Op("=").Make(Index().Byte(), Lit(0), Lit(10485760)), // 10MB
						Id("buf").Op(":=").Qual("bytes", "NewBuffer").Call(Id("wb")),
						List(Id("written"), Id("err")).Op(":=").Qual("io", "Copy").Call(Id("buf"), Id("r").Dot("Body")),
						Do(serverErrorEncoder),
						List(Id("params"), Err()).Op("=").Id("reqDec").Call(
							Id("pathParams"),
							Id("r"),
							Id("wb").Index(Op(":").Id("written")),
						),
						Do(serverErrorEncoder),
					)

					g.List(Id("result"), Err()).Op("=").Id("ep").Call(
						Id("r").Dot("Context").Call(),
						Id("params"),
					)
					g.Do(serverErrorEncoder)
					g.Id("statusCode").Op(":=").Lit(204)
					g.If(Id("respEnc").Op("!=").Nil()).Block(
						Id("statusCode").Op("=").Lit(200),
						List(Id("result"), Err()).Op("=").Id("respEnc").Call(
							Id("result"),
						),
						Do(serverErrorEncoder),

						Id("rw").Dot("Header").Call().Dot("Set").Call(Lit("Content-Type"), Lit("application/json")),
						List(Id("data"), Err()).Op(":=").Qual("encoding/json", "Marshal").Call(Id("result")),
						Do(serverErrorEncoder),
						If(
							List(Id("_"), Err()).Op(":=").Id("rw").Dot("Write").Call(Id("data")),
							Err().Op("!=").Nil(),
						).Block(
							Return(),
						),
					)
					g.Id("rw").Dot("WriteHeader").Call(Id("statusCode"))
				}),
			),
		)
	}

}
