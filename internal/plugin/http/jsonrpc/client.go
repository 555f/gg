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
		f.Type().Id(clientName).StructFunc(func(g *Group) {
			g.Op("*").Qual("github.com/555f/jsonrpc", "Client")
		})
		for _, ep := range s.Endpoints {
			methodRequestName := s.Name + ep.MethodName + "Request"
			recvName := strcase.ToLowerCamel(ep.MethodName)
			resultName := s.Name + ep.MethodName + "BatchResult"
			path := strcase.ToLowerCamel(s.Name) + "." + strcase.ToLowerCamel(ep.MethodName)
			if ep.Path != "" {
				path = ep.Path
			}

			if len(ep.BodyResults) > 0 {
				f.Type().Id(resultName).StructFunc(gen.WrapResponse(ep.WrapResponse, ep.BodyResults, f.Import))
			}

			f.Type().Id(methodRequestName).StructFunc(func(g *Group) {
				g.Id("c").Op("*").Id(clientName)
				g.Id("params").StructFunc(func(g *Group) {
					for _, param := range ep.Params {
						if len(param.Params) > 0 {
							for _, childParam := range param.Params {
								g.Add(makeRequestStructParam(param, childParam, f.Import))
							}
							continue
						}
						g.Add(makeRequestStructParam(nil, param, f.Import))
					}
				})
				g.Id("before").Index().Qual("github.com/555f/jsonrpc", "ClientBeforeFunc")
				g.Id("after").Index().Qual("github.com/555f/jsonrpc", "ClientAfterFunc")
				g.Id("ctx").Qual("context", "Context")
			})

			for _, param := range ep.Params {
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

			f.Func().Params(Id(recvName).Op("*").Id(methodRequestName)).Id("MakeRequest").Params().Params(String(), Any()).
				BlockFunc(func(g *Group) {
					g.Var().Id("params").StructFunc(func(g *Group) {
						for _, param := range ep.BodyParams {
							jsonTag := param.Name
							fld := g.Id(param.FldName)
							if !param.Required {
								jsonTag += ",omitempty"
								fld.Op("*")
							}
							fld.Add(types.Convert(param.Type, f.Import)).Tag(map[string]string{"json": jsonTag})
						}
					})

					for _, param := range ep.BodyParams {
						g.Id("params").Dot(param.FldName).Op("=").Id(recvName).Dot("params").Dot(param.Name)
					}

					g.Return(
						Lit(path),
						Id("params"),
					)
				})

			f.Func().Params(Id(recvName).Op("*").Id(methodRequestName)).Id("SetBefore").Params(
				Id("before").Op("...").Qual("github.com/555f/jsonrpc", "ClientBeforeFunc"),
			).Op("*").Id(methodRequestName).Block(
				Id(recvName).Dot("before").Op("=").Id("before"),
				Return(Id(recvName)),
			)

			f.Func().Params(Id(recvName).Op("*").Id(methodRequestName)).Id("SetAfter").Params(
				Id("after").Op("...").Qual("github.com/555f/jsonrpc", "ClientAfterFunc"),
			).Op("*").Id(methodRequestName).Block(
				Id(recvName).Dot("after").Op("=").Id("after"),
				Return(Id(recvName)),
			)

			f.Func().Params(Id(recvName).Op("*").Id(methodRequestName)).Id("WithContext").Params(
				Id("ctx").Qual("context", "Context"),
			).Op("*").Id(methodRequestName).Block(
				Id(recvName).Dot("ctx").Op("=").Id("ctx"),
				Return(Id(recvName)),
			)

			f.Func().Params(Id(recvName).Op("*").Id(methodRequestName)).Id("RawExecute").Params().Params(
				Index().Byte(),
				Map(Uint64()).Int(),
				Op("*").Qual("net/http", "Response"),
				Error(),
			).Block(
				Return(Id(recvName).Dot("c").Dot("RawExecute").Call(Id(recvName))),
			)

			f.Func().Params(Id(recvName).Op("*").Id(methodRequestName)).Id("MakeResult").Params(Id("data").Index().Byte()).Params(Any(), Error()).
				BlockFunc(func(g *Group) {
					if len(ep.BodyResults) > 0 {
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

			f.Func().Params(Id(recvName).Op("*").Id(methodRequestName)).Id("Before").Params().Index().Qual("github.com/555f/jsonrpc", "ClientBeforeFunc").Block(
				Return(Id(recvName).Dot("before")),
			)

			f.Func().Params(Id(recvName).Op("*").Id(methodRequestName)).Id("After").Params().Index().Qual("github.com/555f/jsonrpc", "ClientAfterFunc").Block(
				Return(Id(recvName).Dot("after")),
			)

			f.Func().Params(Id(recvName).Op("*").Id(methodRequestName)).Id("Context").Params().Do(f.Import("context", "Context")).Block(
				Return(Id(recvName).Dot("ctx")),
			)

			f.Func().Params(Id(recvName).Op("*").Id(methodRequestName)).Id("Execute").Params().
				ParamsFunc(func(g *Group) {
					for _, result := range ep.Sig.Results {
						g.Id(result.Name).Add(types.Convert(result.Type, f.Import))
					}
				}).
				BlockFunc(func(g *Group) {
					batchResultID := Id("batchResult")
					resultAssignOp := ":="
					if len(ep.BodyResults) == 0 {
						batchResultID = Id("_")
						resultAssignOp = "="
					}
					g.List(batchResultID, Err()).Op(resultAssignOp).Id(recvName).Dot("c").Dot("Client").Dot("Execute").Call(Id(recvName))
					g.Do(gen.CheckErr(Return()))

					if len(ep.BodyResults) > 0 {
						g.Id("clientResult").Op(":=").Id("batchResult").Dot("At").Call(Lit(0)).Assert(Id(resultName))
					}

					g.ReturnFunc(func(g *Group) {
						var ids []Code
						for _, name := range ep.WrapResponse {
							ids = append(ids, Dot(strcase.ToCamel(name)))
						}
						for _, result := range ep.Sig.Results {
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
				}).Op("*").Id(methodRequestName).BlockFunc(func(g *Group) {
				g.Id("r").Op(":=").Op("&").Id(methodRequestName).Values(
					Id("ctx").Op(":").Qual("context", "TODO").Call(),
					Id("c").Op(":").Id(recvName),
				)
				for _, param := range endpoint.Params {
					if param.Required {
						g.Id("r").Dot("params").Dot(param.Name).Op("=").Id(param.Name)
					}
				}
				g.Return(Id("r"))
			})
		}
		f.Func().Id("New"+s.Name+"Client").Params(
			Id("target").String(),
			Id("opts").Op("...").Qual("github.com/555f/jsonrpc", "ClientOption"),
		).Op("*").Id(clientName).BlockFunc(
			func(g *Group) {
				g.Return(
					Op("&").Id(clientName).Values(
						Id("Client").Op(":").Qual("github.com/555f/jsonrpc", "NewClient").Call(
							Id("target"),
							Id("opts").Op("..."),
						),
					),
				)
			},
		)
	}
}
