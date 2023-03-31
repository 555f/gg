package jsonrpc

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gen"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"

	. "github.com/dave/jennifer/jen"
)

func GenClient(s options.Iface) func(f *file.GoFile) {
	return func(f *file.GoFile) {
		clientName := s.Name + "Client"
		clientRecvName := strcase.ToLowerCamel(s.Name)

		f.Type().Id(clientName).StructFunc(func(g *Group) {
			g.Id("client").Op("*").Qual("net/http", "Client")
			g.Id("target").String()
			g.Id("incrementID").Uint64()
			g.Id("methodOpts").Op("*").Id("clientMethodOptions")
		})

		f.Func().Params(Id(clientRecvName).Op("*").Id(clientName)).Id("autoIncrementID").Params().Uint64().Block(
			Return(Do(f.Import("sync/atomic", "AddUint64")).Call(Op("&").Id(clientRecvName).Dot("incrementID"), Lit(1))),
		)

		f.Func().Params(Id(clientRecvName).Op("*").Id(clientName)).Id("Execute").Params(Id("requests").Op("...").Id("clientRequester")).Params(Op("*").Id("BatchResult"), Error()).Block(
			Id(clientRecvName).Dot("incrementID").Op("=").Lit(0),
			List(Id("req"), Err()).Op(":=").Do(f.Import("net/http", "NewRequest")).Call(Lit("POST"), Id(clientRecvName).Dot("target"), Nil()),
			Do(gen.CheckErr(Return(Nil(), Err()))),
			Id("idsIndex").Op(":=").Id("make").Call(Map(Uint64()).Int(), Len(Id("requests"))),
			Id("rpcRequests").Op(":=").Id("make").Call(Index().Id("clientReq"), Len(Id("requests"))),
			For(List(Id("_"), Id("beforeFunc")).Op(":=").Range().Id(clientRecvName).Dot("methodOpts").Dot("before")).Block(
				Id("req").Op("=").Id("req").Dot("WithContext").Call(
					Id("beforeFunc").Call(Id("req").Dot("Context").Call(), Id("req")),
				),
			),
			For(List(Id("i"), Id("request")).Op(":=").Range().Id("requests")).Block(
				Id("req").Op("=").Id("req").Dot("WithContext").Call(Id("request").Dot("context").Call()),
				For(List(Id("_"), Id("beforeFunc")).Op(":=").Range().Id("request").Dot("before").Call()).Block(
					Id("req").Op("=").Id("req").Dot("WithContext").Call(
						Id("beforeFunc").Call(Id("req").Dot("Context").Call(), Id("req")),
					),
				),
				List(Id("methodName"), Id("params")).Op(":=").Id("request").Dot("makeRequest").Call(),
				Id("r").Op(":=").Id("clientReq").Values(
					Id("ID").Op(":").Id(clientRecvName).Dot("autoIncrementID").Call(),
					Id("Version").Op(":").Lit("2.0"),
					Id("Method").Op(":").Id("methodName"),
					Id("Params").Op(":").Id("params"),
				),
				Id("idsIndex").Index(Id("r").Dot("ID")).Op("=").Id("i"),
				Id("rpcRequests").Index(Id("i")).Op("=").Id("r"),
			),
			Id("buf").Op(":=").Do(f.Import("bytes", "NewBuffer")).Call(Nil()),
			If(Err().Op(":=").Do(f.Import("encoding/json", "NewEncoder")).Call(Id("buf")).Dot("Encode").Call(Id("rpcRequests")), Err().Op("!=").Nil()).Block(
				Return(Nil(), Err()),
			),
			Id("req").Dot("Body").Op("=").Do(f.Import("io", "NopCloser")).Call(Id("buf")),
			List(Id("resp"), Err()).Op(":=").Id(clientRecvName).Dot("client").Dot("Do").Call(Id("req")),
			Do(gen.CheckErr(Return(Nil(), Err()))),
			Id("responses").Op(":=").Id("make").Call(Index().Id("clientResp"), Len(Id("requests"))),
			If(Err().Op(":=").Do(f.Import("encoding/json", "NewDecoder")).Call(Id("resp").Dot("Body")).Dot("Decode").Call(Op("&").Id("responses")), Err().Op("!=").Nil()).Block(
				Return(Nil(), Err()),
			),
			Id("batchResult").Op(":=").Op("&").Id("BatchResult").Values(
				Id("results").Op(":").Id("make").Call(Index().Any(), Len(Id("requests"))),
			),

			For(List(Id("_"), Id("response")).Op(":=").Range().Id("responses")).Block(
				For(List(Id("_"), Id("afterFunc")).Op(":=").Range().Id(clientRecvName).Dot("methodOpts").Dot("after")).Block(
					Id("afterFunc").Call(Id("resp").Dot("Request").Dot("Context").Call(), Id("resp"), Id("response").Dot("Result")),
				),
				Id("i").Op(":=").Id("idsIndex").Index(Id("response").Dot("ID")),
				Id("request").Op(":=").Id("requests").Index(Id("i")),
				For(List(Id("_"), Id("afterFunc")).Op(":=").Range().Id("request").Dot("after").Call()).Block(
					Id("afterFunc").Call(Id("resp").Dot("Request").Dot("Context").Call(), Id("resp"), Id("response").Dot("Result")),
				),
				List(Id("result"), Err()).Op(":=").Id("request").Dot("makeResult").Call(Id("response").Dot("Result")),
				Do(gen.CheckErr(Return(Nil(), Err()))),
				Id("batchResult").Dot("results").Index(Id("i")).Op("=").Id("result"),
			),
			Return(Id("batchResult"), Nil()),
		)

		for _, endpoint := range s.Endpoints {
			methodRequestName := s.Name + endpoint.MethodName + "Request"
			recvName := strcase.ToLowerCamel(endpoint.MethodName)
			resultName := s.Name + endpoint.MethodName + "BatchResult"
			path := strcase.ToLowerCamel(s.Name) + "." + strcase.ToLowerCamel(endpoint.MethodName)
			if endpoint.Path != "" {
				path = endpoint.Path
			}

			if len(endpoint.BodyResults) > 0 {
				f.Type().Id(resultName).StructFunc(gen.WrapResponse(endpoint.WrapResponse, endpoint.BodyResults, f.Import))
			}

			f.Type().Id(methodRequestName).StructFunc(func(g *Group) {
				g.Id("c").Op("*").Id(clientName)
				g.Id("client").Op("*").Qual("net/http", "Client")
				g.Id("methodOpts").Op("*").Id("clientMethodOptions")
				g.Id("params").StructFunc(func(g *Group) {
					for _, param := range endpoint.Params {
						if len(param.Params) > 0 {
							for _, childParam := range param.Params {
								g.Add(makeRequestStructParam(param, childParam, f.Import))
							}
							continue
						}
						g.Add(makeRequestStructParam(nil, param, f.Import))
					}
				})
			})

			for _, param := range endpoint.Params {
				if len(param.Params) > 0 {
					for _, childParam := range param.Params {
						f.Add(makeSetFunc(recvName, methodRequestName, param, childParam, f.Import))
					}
				} else {
					if !param.Required {
						f.Add(makeSetFunc(recvName, methodRequestName, nil, param, f.Import))
					}
				}
			}

			f.Func().Params(Id(recvName).Op("*").Id(methodRequestName)).Id("makeRequest").Params().Params(String(), Any()).
				BlockFunc(func(g *Group) {
					g.Var().Id("params").StructFunc(func(g *Group) {
						for _, param := range endpoint.BodyParams {
							jsonTag := param.Name
							fld := g.Id(param.FldName)
							if !param.Required {
								jsonTag += ",omitempty"
								fld.Op("*")
							}
							fld.Add(types.Convert(param.Type, f.Import)).Tag(map[string]string{"json": jsonTag})
						}
					})

					for _, param := range endpoint.BodyParams {
						g.Id("params").Dot(param.FldName).Op("=").Id(recvName).Dot("params").Dot(param.Name)
					}

					g.Return(
						Lit(path),
						Id("params"),
					)
				})

			f.Func().Params(Id(recvName).Op("*").Id(methodRequestName)).Id("makeResult").Params(Id("data").Index().Byte()).Params(Any(), Error()).
				BlockFunc(func(g *Group) {
					if len(endpoint.BodyResults) > 0 {
						g.Var().Id("result").Id(resultName)
						g.If(
							Err().Op(":=").Do(f.Import("encoding/json", "Unmarshal")).Call(
								Id("data"),
								Op("&").Id("result"),
							),
							Err().Op("!=").Nil(),
						).Block(
							Return(Nil(), Err()),
						)
						g.Return(Id("result"), Nil())
					} else {
						g.Return(Nil(), Nil())
					}
				})

			f.Func().Params(Id(recvName).Op("*").Id(methodRequestName)).Id("before").Params().Index().Id("ClientBeforeFunc").Block(
				Return(Id(recvName).Dot("methodOpts").Dot("before")),
			)

			f.Func().Params(Id(recvName).Op("*").Id(methodRequestName)).Id("after").Params().Index().Id("ClientAfterFunc").Block(
				Return(Id(recvName).Dot("methodOpts").Dot("after")),
			)

			f.Func().Params(Id(recvName).Op("*").Id(methodRequestName)).Id("context").Params().Do(f.Import("context", "Context")).Block(
				Return(Id(recvName).Dot("methodOpts").Dot("ctx")),
			)

			f.Func().Params(Id(recvName).Op("*").Id(methodRequestName)).Id("Execute").Params().
				ParamsFunc(func(g *Group) {
					for _, result := range endpoint.Sig.Results {
						g.Id(result.Name).Add(types.Convert(result.Type, f.Import))
					}
				}).
				BlockFunc(func(g *Group) {
					batchResultID := Id("batchResult")
					resultAssignOp := ":="
					if len(endpoint.BodyResults) == 0 {
						batchResultID = Id("_")
						resultAssignOp = "="
					}
					g.List(batchResultID, Err()).Op(resultAssignOp).Id(recvName).Dot("c").Dot("Execute").Call(Id(recvName))
					g.Do(gen.CheckErr(Return()))

					if len(endpoint.BodyResults) > 0 {
						g.Id("clientResult").Op(":=").Id("batchResult").Dot("At").Call(Lit(0)).Assert(Id(resultName))
					}

					g.ReturnFunc(func(g *Group) {
						var ids []Code
						for _, name := range endpoint.WrapResponse {
							ids = append(ids, Dot(strcase.ToCamel(name)))
						}
						for _, result := range endpoint.Sig.Results {
							if result.IsError {
								g.Id(result.Name)
								continue
							}
							g.Id("clientResult").Add(ids...).Dot(strcase.ToCamel(result.Name))
						}
					})
				})
		}
		for _, endpoint := range s.Endpoints {
			methodRequestName := s.Name + endpoint.MethodName + "Request"
			recvName := strcase.ToLowerCamel(s.Name)

			f.Func().Params(Id(recvName).Op("*").Id(clientName)).Id(endpoint.MethodName).
				ParamsFunc(func(g *Group) {
					for _, param := range endpoint.Params {
						if param.Required {
							g.Id(param.Name).Add(types.Convert(param.Type, f.Import))
						}
					}
					g.Id("opts").Op("...").Id("ClientMethodOption")
				}).Op("*").Id(methodRequestName).BlockFunc(func(g *Group) {
				g.Id("m").Op(":=").Op("&").Id(methodRequestName).Values(
					Id("client").Op(":").Id(recvName).Dot("client"),
					Id("methodOpts").Op(":").Op("&").Id("clientMethodOptions").Values(
						Id("ctx").Op(":").Qual("context", "TODO").Call(),
					),
					Id("c").Op(":").Id(recvName),
				)
				for _, param := range endpoint.Params {
					if param.Required {
						g.Id("m").Dot("params").Dot(param.Name).Op("=").Id(param.Name)
					}
				}
				g.For(List(Id("_"), Id("o")).Op(":=").Range().Id("opts")).Block(
					Id("o").Call(Id("m").Dot("methodOpts")),
				)
				g.Return(Id("m"))
			})
		}
		f.Func().Id("New"+s.Name+"Client").Params(
			Id("target").String(),
			Id("opts").Op("...").Id("ClientMethodOption"),
		).Op("*").Id(clientName).BlockFunc(
			func(g *Group) {
				g.Id("c").Op(":=").Op("&").Id(clientName).Values(
					Id("target").Op(":").Id("target"),
					Id("client").Op(":").Qual("net/http", "DefaultClient"),
					Id("methodOpts").Op(":").Op("&").Id("clientMethodOptions").Values(
						Id("ctx").Op(":").Qual("context", "TODO").Call(),
					),
				)
				g.For(List(Id("_"), Id("o")).Op(":=").Range().Id("opts")).Block(
					Id("o").Call(Id("c").Dot("methodOpts")),
				)
				g.Return(Id("c"))
			},
		)
	}
}
