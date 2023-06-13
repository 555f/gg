package rest

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/strcase"

	. "github.com/dave/jennifer/jen"
)

func GenErrorEncoder(errorWrapper *options.ErrorWrapper) func(f *file.GoFile) {
	return func(f *file.GoFile) {
		f.Func().Id("serverErrorEncoder").Params(
			//Id("ctx").Qual("context", "Context"),
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
				g.If(Id("jsonErr").Op("!=").Nil()).Block(
					List(Id("_"), Id("_")).Op("=").Id("w").Dot("Write").Call(Index().Byte().Call(Lit("unexpected marshal error"))),
					Return(),
				)

				g.Id("h").Dot("Set").Call(Lit("Content-Type"), Lit("application/json; charset=utf-8"))
			}

			g.Id("w").Dot("WriteHeader").Call(Id("statusCode"))

			if errorWrapper != nil {
				g.List(Id("_"), Id("_")).Op("=").Id("w").Dot("Write").Call(Id("data"))
			}
		})
	}
}
