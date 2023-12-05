package gen

import (
	"github.com/555f/gg/internal/plugin/jsonrpc/options"
	"github.com/dave/jennifer/jen"
)

type BaseClientBuilder struct {
	qualifier Qualifier
	codes     []jen.Code
}

// BuildConstruct implements ClientBuilder.
func (b *BaseClientBuilder) BuildConstruct(iface options.Iface) ClientBuilder {
	clientName := clientStructName(iface)

	b.codes = append(b.codes,
		jen.Func().Id("New"+iface.Name+"Client").Params(
			jen.Id("target").String(),
			jen.Id("opts").Op("...").Qual("github.com/555f/jsonrpc", "ClientOption"),
		).Op("*").Id(clientName).BlockFunc(
			func(g *jen.Group) {
				g.Return(
					jen.Op("&").Id(clientName).Values(
						jen.Id("Client").Op(":").Qual("github.com/555f/jsonrpc", "NewClient").Call(
							jen.Id("target"),
							jen.Id("opts").Op("..."),
						),
					),
				)
			},
		),
	)
	return b
}

// BuildTypes implements ClientBuilder.
func (b *BaseClientBuilder) BuildTypes() ClientBuilder {

	return b
}

// Endpoint implements ClientBuilder.
func (b *BaseClientBuilder) Endpoint(iface options.Iface, ep options.Endpoint) ClientEndpointBuilder {
	return &clientEndpointBuilder{BaseClientBuilder: b, iface: iface, ep: ep, qualifier: b.qualifier}
}

func (b *BaseClientBuilder) BuildStruct(iface options.Iface) ClientBuilder {
	clientName := iface.Name + "Client"
	b.codes = append(b.codes,
		jen.Type().Id(clientName).StructFunc(func(g *jen.Group) {
			g.Op("*").Qual("github.com/555f/jsonrpc", "Client")
		}),
	)
	return b
}

func (b *BaseClientBuilder) Build() jen.Code {
	return jen.Custom(jen.Options{Multi: true}, b.codes...)
}

func NewBaseClientBuilder(qualifier Qualifier) *BaseClientBuilder {
	return &BaseClientBuilder{qualifier: qualifier}
}
