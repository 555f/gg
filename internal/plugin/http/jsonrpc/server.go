package jsonrpc

import (
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gen"

	. "github.com/dave/jennifer/jen"
)

func GenServer() func(f *file.GoFile) {
	return func(f *file.GoFile) {
		f.Const().Id("jsonRPCParseError").Int().Op("=").Lit(-32700)
		f.Const().Id("jsonRPCInvalidRequestError").Int().Op("=").Lit(-32600)
		f.Const().Id("jsonRPCMethodNotFoundError").Int().Op("=").Lit(-32601)
		f.Const().Id("jsonRPCInvalidParamsError").Int().Op("=").Lit(-32602)
		f.Const().Id("jsonRPCInternalError").Int().Op("=").Lit(-32603)

		f.Type().Id("JSONRPCHandlerFunc").Func().Params(
			Id("ctx").Qual("context", "Context"),
			Id("rw").Qual("net/http", "ResponseWriter"),
			Id("r").Op("*").Qual("net/http", "Request"),
			Id("prams").Qual("encoding/json", "RawMessage"),
		).Params(Any(), Error())

		f.Type().Id("jsonRPCRoute").Struct(
			Id("handler").Id("JSONRPCHandlerFunc"),
			Id("before").Index().Id("ServerBeforeFunc"),
			Id("after").Index().Id("ServerAfterFunc"),
		)
		f.Type().Id("Server").StructFunc(func(g *Group) {
			g.Id("routes").Map(String()).Id("jsonRPCRoute")
		})

		f.Type().Id("jsonRPCError").Struct(
			Id("Code").Int().Tag(map[string]string{"json": "code"}),
			Id("Message").String().Tag(map[string]string{"json": "message"}),
			Id("Data").Any().Tag(map[string]string{"json": "data,omitempty"}),
		)

		f.Type().Id("jsonRPCRequestData").Struct(
			Id("requests").Index().Id("jsonRPCRequest"),
			Id("isBatch").Bool(),
		)

		f.Func().Params(Id("r").Op("*").Id("jsonRPCRequestData")).Id("UnmarshalJSON").Params(Id("b").Index().Byte()).Error().Block(
			If(Qual("bytes", "HasPrefix").Call(Id("b"), Index().Byte().Call(Lit("[")))).Block(
				Id("r").Dot("isBatch").Op("=").True(),
				Return(
					Qual("encoding/json", "Unmarshal").Call(Id("b"), Op("&").Id("r").Dot("requests")),
				),
			),
			Var().Id("req").Id("jsonRPCRequest"),
			If(Err().Op(":=").Qual("encoding/json", "Unmarshal").Call(Id("b"), Op("&").Id("req")), Err().Op("!=").Nil()).Block(
				Return(Err()),
			),
			Id("r").Dot("requests").Op("=").Append(Id("r").Dot("requests"), Id("req")),
			Return(Nil()),
		)

		f.Type().Id("jsonRPCRequest").Struct(
			Id("ID").Any().Tag(map[string]string{"json": "id"}),
			Id("Version").String().Tag(map[string]string{"json": "jsonrpc"}),
			Id("Method").String().Tag(map[string]string{"json": "method"}),
			Id("Params").Qual("encoding/json", "RawMessage").Tag(map[string]string{"json": "params"}),
		)

		f.Type().Id("jsonRPCResponse").Struct(
			Id("ID").Any().Tag(map[string]string{"json": "id"}),
			Id("Version").String().Tag(map[string]string{"json": "jsonrpc"}),
			Id("Error").Op("*").Id("jsonRPCError").Tag(map[string]string{"json": "error,omitempty"}),
			Id("Result").Qual("encoding/json", "RawMessage").Tag(map[string]string{"json": "result,omitempty"}),
		)

		f.Func().Params(Id("s").Op("*").Id("Server")).Id("addRoute").Params(
			Id("method").String(),
			Id("handler").Id("JSONRPCHandlerFunc"),
			Id("before").Index().Id("ServerBeforeFunc"),
			Id("after").Index().Id("ServerAfterFunc"),
		).Block(
			Id("s").Dot("routes").Index(Id("method")).Op("=").Id("jsonRPCRoute").Values(
				Id("handler").Op(":").Id("handler"),
				Id("before").Op(":").Id("before"),
				Id("after").Op(":").Id("after"),
			),
		)

		f.Func().Params(Id("s").Op("*").Id("Server")).Id("makeErrorResponse").Params(
			Id("id").Any(),
			Id("code").Int(),
			Id("message").String(),
		).Id("jsonRPCResponse").Block(
			Return(
				Id("jsonRPCResponse").Values(
					Id("ID").Op(":").Id("id"),
					Id("Version").Op(":").Lit("2.0"),
					Id("Error").Op(":").Op("&").Id("jsonRPCError").Values(
						Id("Code").Op(":").Id("code"),
						Id("Message").Op(":").Id("message"),
					),
				),
			),
		)

		f.Func().Params(Id("s").Op("*").Id("Server")).Id("handleRoute").Params(
			Id("route").Id("jsonRPCRoute"),
			Id("ctx").Qual("context", "Context"),
			Id("w").Qual("net/http", "ResponseWriter"),
			Id("r").Op("*").Qual("net/http", "Request"),
			Id("prams").Qual("encoding/json", "RawMessage"),
		).Params(Id("resp").Any(), Err().Error()).Block(
			For(List(Id("_"), Id("before")).Op(":=").Range().Id("route").Dot("before")).Block(
				Id("ctx").Op("=").Id("before").Call(Id("ctx"), Id("r")),
			),
			List(Id("resp"), Err()).Op("=").Id("route").Dot("handler").Call(Id("ctx"), Id("w"), Id("r"), Id("prams")),
			Do(gen.CheckErr(
				Return(Nil(), Err()),
			)),
			For(List(Id("_"), Id("after")).Op(":=").Range().Id("route").Dot("after")).Block(
				Id("ctx").Op("=").Id("after").Call(Id("ctx"), Id("w")),
			),
			Return(),
		)

		f.Func().Params(Id("s").Op("*").Id("Server")).Id("ServeHTTP").Params(
			Id("w").Qual("net/http", "ResponseWriter"),
			Id("r").Op("*").Qual("net/http", "Request"),
		).Block(
			Id("ctx").Op(":=").Id("r").Dot("Context").Call(),

			Var().Id("requestData").Id("jsonRPCRequestData"),
			Var().Id("responses").Index().Id("jsonRPCResponse"),

			If(Err().Op(":=").Qual("encoding/json", "NewDecoder").Call(Id("r").Dot("Body")).Dot("Decode").Call(Op("&").Id("requestData")), Err().Op("!=").Nil()).Block(
				Id("responses").Op("=").Append(Id("responses"), Id("s").Dot("makeErrorResponse").Call(Nil(), Id("jsonRPCParseError"), Err().Dot("Error").Call())),
			).Else().Block(
				For(List(Id("_"), Id("req")).Op(":=").Range().Id("requestData.requests")).Block(
					If(List(Id("route"), Id("ok")).Op(":=").Id("s").Dot("routes").Index(Id("req").Dot("Method")), Id("ok")).Block(
						List(Id("resp"), Err()).Op(":=").Id("s").Dot("handleRoute").Call(Id("route"), Id("ctx"), Id("w"), Id("r"), Id("req").Dot("Params")),
						Do(gen.CheckErr(
							Id("responses").Op("=").Append(
								Id("responses"),
								Id("s").Dot("makeErrorResponse").Call(Id("req").Dot("ID"), Id("jsonRPCInternalError"), Err().Dot("Error").Call()),
							),
							Continue(),
						)),
						List(Id("result"), Err()).Op(":=").Qual("encoding/json", "Marshal").Call(Id("resp")),
						Do(gen.CheckErr(
							Id("responses").Op("=").Append(
								Id("responses"),
								Id("s").Dot("makeErrorResponse").Call(Id("req").Dot("ID"), Id("jsonRPCInternalError"), Err().Dot("Error").Call()),
							),
							Continue(),
						)),
						Id("responses").Op("=").Append(
							Id("responses"),
							Id("jsonRPCResponse").Values(
								Id("ID").Op(":").Id("req").Dot("ID"),
								Id("Version").Op(":").Lit("2.0"),
								Id("Result").Op(":").Id("result"),
							),
						),
					),
				),
			),
			Var().Id("data").Any(),
			If(Id("requestData").Dot("isBatch")).Block(
				Id("data").Op("=").Id("responses"),
			).Else().Block(
				Id("data").Op("=").Id("responses").Index(Lit(0)),
			),
			Id("_").Op("=").Qual("encoding/json", "NewEncoder").Call(Id("w")).Dot("Encode").Call(Id("data")),
		)

		f.Func().Id("NewJSONRPCServer").Params(Id("opts").Op("...").Id("ServerOption")).Op("*").Id("Server").Block(
			Id("s").Op(":=").Op("&").Id("Server").Values(
				Id("routes").Op(":").Make(Map(String()).Id("jsonRPCRoute")),
			),
			For(List(Id("_"), Id("opt"))).Op(":=").Range().Id("opts").Block(
				Id("opt").Call(Id("s")),
			),
			Return(Id("s")),
		)
	}
}
