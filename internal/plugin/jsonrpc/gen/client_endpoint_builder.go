package gen

import (
	"github.com/555f/gg/internal/plugin/jsonrpc/options"
	"github.com/555f/gg/pkg/gen"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

type clientEndpointBuilder struct {
	*BaseClientBuilder
	iface     options.Iface
	ep        options.Endpoint
	qualifier Qualifier
}

func (b *clientEndpointBuilder) BuildReqStruct() ClientEndpointBuilder {
	methodRequestName := b.iface.Name + b.ep.MethodName + "Request"
	resultName := b.iface.Name + b.ep.MethodName + "BatchResult"
	clientName := clientStructName(b.iface)
	recvName := strcase.ToLowerCamel(b.ep.MethodName)

	b.codes = append(b.codes,
		jen.Type().Id(methodRequestName).StructFunc(func(g *jen.Group) {
			g.Id("c").Op("*").Id(clientName)
			g.Id("params").StructFunc(func(g *jen.Group) {
				for _, param := range b.ep.Params {
					if len(param.Params) > 0 {
						for _, childParam := range param.Params {
							g.Add(makeRequestStructParam(param, childParam, b.qualifier.Qual))
						}
						continue
					}
					g.Add(makeRequestStructParam(nil, param, b.qualifier.Qual))
				}
			})
			g.Id("before").Index().Qual("github.com/555f/jsonrpc", "ClientBeforeFunc")
			g.Id("after").Index().Qual("github.com/555f/jsonrpc", "ClientAfterFunc")
			// g.Id("ctx").Qual("context", "Context")
		}),
		jen.Func().Params(jen.Id(recvName).Op("*").Id(methodRequestName)).Id("Before").Params().
			Index().Qual("github.com/555f/jsonrpc", "ClientBeforeFunc").Block(
			jen.Return(jen.Id(recvName).Dot("before")),
		),
		jen.Func().Params(jen.Id(recvName).Op("*").Id(methodRequestName)).Id("SetBefore").Params(
			jen.Id("before").Op("...").Qual("github.com/555f/jsonrpc", "ClientBeforeFunc"),
		).Op("*").Id(methodRequestName).Block(
			jen.Id(recvName).Dot("before").Op("=").Id("before"),
			jen.Return(jen.Id(recvName)),
		),
		jen.Func().Params(jen.Id(recvName).Op("*").Id(methodRequestName)).Id("After").Params().
			Index().Qual("github.com/555f/jsonrpc", "ClientAfterFunc").Block(
			jen.Return(jen.Id(recvName).Dot("after")),
		),
		jen.Func().Params(jen.Id(recvName).Op("*").Id(methodRequestName)).Id("SetAfter").Params(
			jen.Id("after").Op("...").Qual("github.com/555f/jsonrpc", "ClientAfterFunc"),
		).Op("*").Id(methodRequestName).Block(
			jen.Id(recvName).Dot("after").Op("=").Id("after"),
			jen.Return(jen.Id(recvName)),
		),
	)

	if len(b.ep.Results) > 0 {
		b.codes = append(b.codes,
			jen.Type().Id(resultName).StructFunc(func(g *jen.Group) {
				for _, result := range b.ep.Results {
					g.Id(result.FldNameExport).Add(types.Convert(result.Type, b.qualifier.Qual)).Tag(map[string]string{"json": result.Name})
				}
			}),
		)
	}

	for _, param := range b.ep.Params {
		if len(param.Params) > 0 {
			for _, childParam := range param.Params {
				b.codes = append(b.codes, makeSetFunc(recvName, methodRequestName, param, childParam, b.qualifier.Qual))
			}
		} else {
			if !param.Required {
				b.codes = append(b.codes, makeSetFunc(recvName, methodRequestName, nil, param, b.qualifier.Qual))
			}
		}
	}
	return b
}

// BuildExecuteMethod implements ClientEndpointBuilder.
func (b *clientEndpointBuilder) BuildExecuteMethod() ClientEndpointBuilder {
	methodRequestName := b.iface.Name + b.ep.MethodName + "Request"
	recvName := strcase.ToLowerCamel(b.ep.MethodName)
	resultName := b.iface.Name + b.ep.MethodName + "BatchResult"

	b.codes = append(b.codes,
		jen.Func().Params(jen.Id(recvName).Op("*").Id(methodRequestName)).Id("Execute").Params().
			ParamsFunc(func(g *jen.Group) {
				for _, result := range b.ep.Sig.Results {
					g.Id(result.Name).Add(types.Convert(result.Type, b.qualifier.Qual))
				}
			}).
			BlockFunc(func(g *jen.Group) {
				batchResultID := jen.Id("batchResult")
				resultAssignOp := ":="

				if len(b.ep.Results) == 0 && b.ep.Error == nil {
					batchResultID = jen.Id("_")
					resultAssignOp = "="
				}
				g.List(batchResultID, jen.Err()).Op(resultAssignOp).Id(recvName).Dot("c").Dot("Client").Dot("Execute").Call(jen.Id(recvName))
				g.Do(gen.CheckErr(jen.Return()))

				if b.ep.Error != nil {
					g.Err().Op("=").Id("batchResult").Dot("Error").Call(jen.Lit(0))
					g.Do(gen.CheckErr(jen.Return()))
				}

				if len(b.ep.Results) > 0 {
					g.Id("clientResult").Op(":=").Id("batchResult").Dot("At").Call(jen.Lit(0)).Assert(jen.Id(resultName))
				}

				g.ReturnFunc(func(g *jen.Group) {
					var ids []jen.Code
					for _, result := range b.ep.Sig.Results {
						if result.IsError {
							g.Id(result.Name)
							continue
						}
						g.Id("clientResult").Add(ids...).Dot(strcase.ToCamel(result.Name))
					}
				})
			}),
	)
	return b
}

func (b *clientEndpointBuilder) BuildMethod() ClientEndpointBuilder {
	methodRequestName := b.iface.Name + b.ep.MethodName + "Request"
	recvName := strcase.ToLowerCamel(b.iface.Name)
	clientName := clientStructName(b.iface)
	b.codes = append(b.codes,
		jen.Func().Params(jen.Id(recvName).Op("*").Id(clientName)).Id(b.ep.MethodName).
			ParamsFunc(func(g *jen.Group) {
				for _, param := range b.ep.Params {
					if param.Required {
						g.Id(param.Name).Add(types.Convert(param.Type, b.qualifier.Qual))
					}
				}
			}).Op("*").Id(methodRequestName).BlockFunc(func(g *jen.Group) {
			g.Id("r").Op(":=").Op("&").Id(methodRequestName).Values(
				// jen.Id("ctx").Op(":").Qual("context", "TODO").Call(),
				jen.Id("c").Op(":").Id(recvName),
			)
			for _, param := range b.ep.Params {
				if param.Required {
					g.Id("r").Dot("params").Dot(param.Name).Op("=").Id(param.Name)
				}
			}
			g.Return(jen.Id("r"))
		}),
	)
	return b
}

func (b *clientEndpointBuilder) BuildResultMethod() ClientEndpointBuilder {
	methodRequestName := b.iface.Name + b.ep.MethodName + "Request"
	recvName := strcase.ToLowerCamel(b.ep.MethodName)
	resultName := b.iface.Name + b.ep.MethodName + "BatchResult"
	b.codes = append(b.codes,
		jen.Func().Params(jen.Id(recvName).Op("*").Id(methodRequestName)).Id("MakeResult").Params(jen.Id("data").Index().Byte()).Params(jen.Any(), jen.Error()).
			BlockFunc(func(g *jen.Group) {
				if len(b.ep.Results) > 0 {
					g.Var().Id("result").Id(resultName)
					g.If(
						jen.Err().Op(":=").Do(b.qualifier.Qual("encoding/json", "Unmarshal")).Call(
							jen.Id("data"),
							jen.Op("&").Id("result"),
						),
						jen.Err().Op("!=").Nil(),
					).Block(
						jen.Return(jen.Nil(), jen.Err()),
					)
					g.Return(jen.Id("result"), jen.Nil())
				} else {
					g.Return(jen.Nil(), jen.Nil())
				}
			}),
	)
	return b
}

// BuildReqMethod implements ClientEndpointBuilder.
func (b *clientEndpointBuilder) BuildReqMethod() ClientEndpointBuilder {
	methodRequestName := b.iface.Name + b.ep.MethodName + "Request"
	recvName := strcase.ToLowerCamel(b.ep.MethodName)

	b.codes = append(b.codes,
		jen.Func().Params(jen.Id(recvName).Op("*").Id(methodRequestName)).Id("MakeRequest").Params().Params(jen.String(), jen.Any()).
			BlockFunc(func(g *jen.Group) {
				g.Var().Id("params").StructFunc(func(g *jen.Group) {
					for _, param := range b.ep.Params {
						jsonTag := param.Name
						fld := g.Id(param.FldName)
						if !param.Required {
							jsonTag += ",omitempty"
							fld.Op("*")
						}
						fld.Add(types.Convert(param.Type, b.qualifier.Qual)).Tag(map[string]string{"json": jsonTag})
					}
				})
				for _, param := range b.ep.Params {
					g.Id("params").Dot(param.FldName).Op("=").Id(recvName).Dot("params").Dot(param.Name)
				}
				g.Return(
					jen.Lit(b.ep.RPCMethodName),
					jen.Id("params"),
				)
			}),
	)
	return b
}

// BuildSetters implements ClientEndpointBuilder.
func (b *clientEndpointBuilder) BuildSetters() ClientEndpointBuilder {
	return b
}
