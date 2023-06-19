package generic

import (
	"github.com/555f/gg/pkg/file"

	. "github.com/dave/jennifer/jen"
)

func GenRESTClient() func(f *file.GoFile) {
	return func(f *file.GoFile) {
		f.Type().Id("ClientBeforeFunc").Func().Params(
			Qual("context", "Context"),
			Op("*").Qual("net/http", "Request"),
		).Qual("context", "Context")
		f.Type().Id("ClientAfterFunc").Func().Params(
			Qual("context", "Context"),
			Op("*").Qual("net/http", "Response"),
		).Qual("context", "Context")

		f.Type().Id("clientMethodOptions").Struct(
			Id("ctx").Qual("context", "Context"),
			Id("before").Index().Id("ClientBeforeFunc"),
			Id("after").Index().Id("ClientAfterFunc"),
		)
		f.Type().Id("ClientMethodOption").Func().Params(Op("*").Id("clientMethodOptions"))
		f.Func().Id("WithContext").Params(Id("ctx").Qual("context", "Context")).Id("ClientMethodOption").Block(
			Return(Func().Params(Id("o").Op("*").Id("clientMethodOptions")).Block(
				Id("o").Dot("ctx").Op("=").Id("ctx"),
			)),
		)
		f.Func().Id("Before").Params(Id("before").Op("...").Id("ClientBeforeFunc")).Id("ClientMethodOption").Block(
			Return(Func().Params(Id("o").Op("*").Id("clientMethodOptions")).Block(
				Id("o").Dot("before").Op("=").Append(Id("o").Dot("before"), Id("before").Op("...")),
			)),
		)
		f.Func().Id("After").Params(Id("after").Op("...").Id("ClientAfterFunc")).Id("ClientMethodOption").Block(
			Return(Func().Params(Id("o").Op("*").Id("clientMethodOptions")).Block(
				Id("o").Dot("after").Op("=").Append(Id("o").Dot("after"), Id("after").Op("...")),
			)),
		)
	}
}
