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
	}
}

func GenJSONRPCClient() func(f *file.GoFile) {
	return func(f *file.GoFile) {
		f.Type().Id("ClientBeforeFunc").Func().Params(
			Qual("context", "Context"),
			Op("*").Qual("net/http", "Request"),
		).Qual("context", "Context")
		f.Type().Id("ClientAfterFunc").Func().Params(
			Qual("context", "Context"),
			Op("*").Qual("net/http", "Response"),
			Qual("encoding/json", "RawMessage"),
		).Qual("context", "Context")

		f.Type().Id("clientRequester").Interface(
			Id("makeRequest").Params().Params(String(), Any()),
			Id("makeResult").Params(Id("data").Index().Byte()).Params(Any(), Error()),
			Id("before").Params().Index().Id("ClientBeforeFunc"),
			Id("after").Params().Index().Id("ClientAfterFunc"),
			Id("context").Params().Do(f.Import("context", "Context")),
		)

		f.Type().Id("BatchResult").Struct(
			Id("results").Index().Any(),
		)

		f.Func().Params(Id("r").Op("*").Id("BatchResult")).Id("At").Params(Id("i").Int()).Any().Block(
			Return(Id("r").Dot("results").Index(Id("i"))),
		)

		f.Func().Params(Id("r").Op("*").Id("BatchResult")).Id("Len").Params().Int().Block(
			Return(Len(Id("r").Dot("results"))),
		)

		f.Type().Id("clientReq").Struct(
			Id("ID").Uint64().Tag(map[string]string{"json": "id"}),
			Id("Version").String().Tag(map[string]string{"json": "jsonrpc"}),
			Id("Method").String().Tag(map[string]string{"json": "method"}),
			Id("Params").Any().Tag(map[string]string{"json": "params"}),
		)

		f.Type().Id("clientResp").Struct(
			Id("ID").Uint64().Tag(map[string]string{"json": "id"}),
			Id("Version").String().Tag(map[string]string{"json": "jsonrpc"}),
			Id("Error").Op("*").Id("clientError").Tag(map[string]string{"json": "error"}),
			Id("Result").Do(f.Import("encoding/json", "RawMessage")).Tag(map[string]string{"json": "result"}),
		)

		f.Type().Id("clientError").Struct(
			Id("Code").Error().Tag(map[string]string{"json": "code"}),
			Id("Message").String().Tag(map[string]string{"json": "message"}),
			Id("Data").Any().Tag(map[string]string{"json": "data"}),
		)
	}
}

func GenClient() func(f *file.GoFile) {
	return func(f *file.GoFile) {
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
