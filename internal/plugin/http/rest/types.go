package rest

import (
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gen"
	. "github.com/dave/jennifer/jen"
)

var (
	epFunc = Func().Params(
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
	}
}
