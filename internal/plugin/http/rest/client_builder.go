package rest

import (
	"path/filepath"

	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/strcase"
	"github.com/dave/jennifer/jen"
)

const clientOptionName = "clientOptions"

type BaseClientBuilder struct {
	errorWrapper *options.ErrorWrapper
	qualifier    Qualifier
	codes        []jen.Code
}

func (b *BaseClientBuilder) SetErrorWrapper(errorWrapper *options.ErrorWrapper) ClientBuilder {
	b.errorWrapper = errorWrapper
	return b
}

func (b *BaseClientBuilder) Build() jen.Code {
	return jen.Custom(jen.Options{Multi: true}, b.codes...)
}

func (b *BaseClientBuilder) BuildTypes() ClientBuilder {
	labelContextMethodShortName := jen.Id("labelFromContext").Call(jen.Lit("methodNameShort"), jen.Id("shortMethodContextKey"))
	labelContextMethodFullName := jen.Id("labelFromContext").Call(jen.Lit("methodNameFull"), jen.Id("methodContextKey"))
	labelContextScopeName := jen.Id("labelFromContext").Call(jen.Lit("scopeName"), jen.Id("scopeNameContextKey"))

	b.codes = append(b.codes,
		jen.Type().Id("contextKey").String(),
		jen.Const().Id("methodContextKey").Id("contextKey").Op("=").Lit("method"),
		jen.Const().Id("shortMethodContextKey").Id("contextKey").Op("=").Lit("shortMethod"),
		jen.Const().Id("scopeNameContextKey").Id("contextKey").Op("=").Lit("scopeName"),

		jen.Func().Id("labelFromContext").Params(
			jen.Id("lblName").String(),
			jen.Id("ctxKey").Id("contextKey"),
		).Qual(promhttpPkg, "Option").Block(
			jen.Return(
				jen.Qual(promhttpPkg, "WithLabelFromCtx").Call(
					jen.Id("lblName"),
					jen.Func().Params(jen.Id("ctx").Qual(ctxPkg, "Context")).String().Block(
						jen.List(jen.Id("v"), jen.Id("_")).Op(":=").Id("ctx").Dot("Value").Call(jen.Id("ctxKey")).Assert(jen.String()),
						jen.Return(jen.Id("v")),
					),
				),
			),
		),
		jen.Func().Id("instrumentRoundTripperErrCounter").Params(
			jen.Id("counter").Op("*").Qual(prometheusPkg, "CounterVec"),
			jen.Id("next").Qual(httpPkg, "RoundTripper"),
		).Qual(promhttpPkg, "RoundTripperFunc").Block(
			jen.Return(
				jen.Func().
					Params(
						jen.Id("r").Op("*").Qual(httpPkg, "Request"),
					).
					Params(
						jen.Op("*").Qual(httpPkg, "Response"),
						jen.Error(),
					).
					Block(
						jen.Id("labels").Op(":=").Qual(prometheusPkg, "Labels").Values(
							jen.Lit("method").Op(":").Qual(stringsPkg, "ToLower").Call(jen.Id("r").Dot("Method")),
						),
						jen.List(jen.Id("labels").Index(jen.Lit("methodNameFull")), jen.Id("_")).Op("=").Id("r").Dot("Context").Call().Dot("Value").Call(jen.Id("methodContextKey")).Assert(jen.String()),
						jen.List(jen.Id("labels").Index(jen.Lit("methodNameShort")), jen.Id("_")).Op("=").Id("r").Dot("Context").Call().Dot("Value").Call(jen.Id("shortMethodContextKey")).Assert(jen.String()),
						jen.List(jen.Id("labels").Index(jen.Lit("scopeName")), jen.Id("_")).Op("=").Id("r").Dot("Context").Call().Dot("Value").Call(jen.Id("scopeNameContextKey")).Assert(jen.String()),
						jen.List(jen.Id("labels").Index(jen.Lit("code"))).Op("=").Lit(""),
						jen.List(jen.Id("resp"), jen.Err()).Op(":=").Id("next").Dot("RoundTrip").Call(jen.Id("r")),
						jen.If(jen.Id("err").Op("!=").Nil()).Block(
							jen.Var().Id("errType").String(),
							jen.Switch(jen.Id("e").Op(":=").Err().Assert(jen.Id("type"))).Block(
								jen.Default().Block(
									jen.Id("errType").Op("=").Err().Dot("Error").Call(),
								),
								jen.Case(jen.Op("*").Qual(tlsPkg, "CertificateVerificationError")).Block(
									jen.Id("errType").Op("=").Lit("failedVerifyCertificate"),
								),
								jen.Case(jen.Qual(netPkg, "Error")).Block(
									jen.Id("errType").Op("+=").Lit("net."),
									jen.If(jen.Id("e").Dot("Timeout").Call()).Block(
										jen.Id("errType").Op("+=").Lit("timeout."),
									),
									jen.Switch(jen.Id("ee").Op(":=").Id("e").Assert(jen.Id("type"))).Block(
										jen.Case(jen.Op("*").Qual(netPkg, "ParseError")).Block(
											jen.Id("errType").Op("+=").Lit("parse"),
										),
										jen.Case(jen.Op("*").Qual(netPkg, "InvalidAddrError")).Block(
											jen.Id("errType").Op("+=").Lit("invalidAddr"),
										),
										jen.Case(jen.Op("*").Qual(netPkg, "UnknownNetworkError")).Block(
											jen.Id("errType").Op("+=").Lit("unknownNetwork"),
										),
										jen.Case(jen.Op("*").Qual(netPkg, "DNSError")).Block(
											jen.Id("errType").Op("+=").Lit("dns"),
										),
										jen.Case(jen.Op("*").Qual(netPkg, "OpError")).Block(
											jen.Id("errType").Op("+=").Id("ee").Dot("Net").Op("+").Lit(".").Op("+").Id("ee").Dot("Op"),
										),
									),
								),
							),
							jen.Id("labels").Index(jen.Lit("err")).Op("=").Id("errType"),
							jen.Id("counter").Dot("With").Call(jen.Id("labels")).Dot("Add").Call(jen.Lit(1)),
						).Else().If(jen.Id("resp").Dot("StatusCode").Op(">").Lit(399)).Block(
							jen.List(jen.Id("labels").Index(jen.Lit("code"))).Op("=").Qual(strconvPkg, "Itoa").Call(jen.Id("resp").Dot("StatusCode")),
							jen.Id("labels").Index(jen.Lit("err")).Op("=").Lit("respFailed"),
							jen.Id("counter").Dot("With").Call(jen.Id("labels")).Dot("Add").Call(jen.Lit(1)),
						),

						jen.Return(jen.Id("resp"), jen.Err()),
					),
			),
		),

		jen.Type().Id("PrometheusCollector").Interface(
			jen.Qual(prometheusPkg, "Collector"),
			// jen.Id("Inflight").Params().Params(jen.Qual(prometheusPkg, "Gauge")),
			jen.Id("Requests").Params().Params(jen.Op("*").Qual(prometheusPkg, "CounterVec")),
			jen.Id("ErrRequests").Params().Params(jen.Op("*").Qual(prometheusPkg, "CounterVec")),
			jen.Id("Duration").Params().Params(jen.Op("*").Qual(prometheusPkg, "HistogramVec")),
			// jen.Id("DNSDuration").Params().Params(jen.Op("*").Qual(prometheusPkg, "HistogramVec")),
			// jen.Id("TLSDuration").Params().Params(jen.Op("*").Qual(prometheusPkg, "HistogramVec")),
		),

		jen.Type().Id("ClientBeforeFunc").Func().Params(
			jen.Qual("context", "Context"),
			jen.Op("*").Qual("net/http", "Request"),
		).Params(jen.Qual("context", "Context"), jen.Error()),
		jen.Type().Id("ClientAfterFunc").Func().Params(
			jen.Qual("context", "Context"),
			jen.Op("*").Qual("net/http", "Response"),
		).Qual("context", "Context"),

		jen.Type().Id(clientOptionName).Struct(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("before").Index().Id("ClientBeforeFunc"),
			jen.Id("after").Index().Id("ClientAfterFunc"),
			jen.Id("client").Op("*").Qual(httpPkg, "Client"),
		),
		jen.Type().Id("ClientOption").Func().Params(jen.Op("*").Id(clientOptionName)),
		jen.Func().Id("WithContext").Params(jen.Id("ctx").Qual("context", "Context")).Id("ClientOption").Block(
			jen.Return(jen.Func().Params(jen.Id("o").Op("*").Id(clientOptionName)).Block(
				jen.Id("o").Dot("ctx").Op("=").Id("ctx"),
			)),
		),
		jen.Func().Id("WithClient").Params(jen.Id("client").Op("*").Qual(httpPkg, "Client")).Id("ClientOption").Block(
			jen.Return(jen.Func().Params(jen.Id("o").Op("*").Id(clientOptionName)).Block(
				jen.Id("o").Dot("client").Op("=").Id("client"),
			)),
		),
		jen.Func().Id("WithPrometheusCollector").Params(jen.Id("c").Id("PrometheusCollector")).Id("ClientOption").Block(
			jen.Return(jen.Func().Params(jen.Id("o").Op("*").Id(clientOptionName)).Block(
				jen.If(jen.Id("o").Dot("client").Dot("Transport").Op("==").Nil()).Block(
					jen.Panic(jen.Lit("no transport is set for the http client")),
				),
				// jen.Id("trace").Op(":=").Op("&").Qual(promhttpPkg, "InstrumentTrace").Values(),
				jen.Id("o").Dot("client").Dot("Transport").Op("=").
					Id("instrumentRoundTripperErrCounter").Call(jen.Id("c").Dot("ErrRequests").Call(),
					// jen.Qual(promhttpPkg, "InstrumentRoundTripperInFlight").Call(
					// jen.Id("c").Dot("Inflight").Call(),
					jen.Qual(promhttpPkg, "InstrumentRoundTripperCounter").Call(
						jen.Id("c").Dot("Requests").Call(),
						// jen.Qual(promhttpPkg, "InstrumentRoundTripperTrace").Call(
						// jen.Id("trace"),
						jen.Qual(promhttpPkg, "InstrumentRoundTripperDuration").Call(
							jen.Id("c").Dot("Duration").Call(),
							jen.Id("o").Dot("client").Dot("Transport"),
							labelContextMethodShortName,
							labelContextMethodFullName,
							labelContextScopeName,
						),
						// ),
						labelContextMethodShortName,
						labelContextMethodFullName,
						labelContextScopeName,
					),
					// ),
				),
			)),
		),
		jen.Func().Id("Before").Params(jen.Id("before").Op("...").Id("ClientBeforeFunc")).Id("ClientOption").Block(
			jen.Return(jen.Func().Params(jen.Id("o").Op("*").Id(clientOptionName)).Block(
				jen.Id("o").Dot("before").Op("=").Append(jen.Id("o").Dot("before"), jen.Id("before").Op("...")),
			)),
		),
		jen.Func().Id("After").Params(jen.Id("after").Op("...").Id("ClientAfterFunc")).Id("ClientOption").Block(
			jen.Return(jen.Func().Params(jen.Id("o").Op("*").Id(clientOptionName)).Block(
				jen.Id("o").Dot("after").Op("=").Append(jen.Id("o").Dot("after"), jen.Id("after").Op("...")),
			)),
		),
	)
	return b
}

func (b *BaseClientBuilder) BuildStruct(iface options.Iface) ClientBuilder {
	clientName := clientStructName(iface)
	b.codes = append(b.codes,
		jen.Const().Id(strcase.ToLowerCamel(iface.Name)+"ScopeName").Op("=").Lit(filepath.Base(iface.PkgPath)),
		jen.Type().Id(clientName).StructFunc(func(g *jen.Group) {
			g.Id("target").String()
			g.Id("opts").Op("*").Id(clientOptionName)
		}))
	return b
}

func (b *BaseClientBuilder) BuildConstruct(iface options.Iface) ClientBuilder {
	clientName := clientStructName(iface)
	b.codes = append(b.codes, jen.Func().Id("New"+iface.Name+"Client").
		Params(
			jen.Id("target").String(),
			jen.Id("opts").Op("...").Id("ClientOption"),
		).Op("*").Id(clientName).BlockFunc(
		func(g *jen.Group) {
			g.Id("c").Op(":=").Op("&").Id(clientName).Values(
				jen.Id("target").Op(":").Id("target"),
				jen.Id("opts").Op(":").Op("&").Id(clientOptionName).Values(
					jen.Id("client").Op(":").Qual(cleanhttpPkg, "DefaultClient").Call(),
				),
			)
			g.For(jen.List(jen.Id("_"), jen.Id("o")).Op(":=").Range().Id("opts")).Block(
				jen.Id("o").Call(jen.Id("c").Dot("opts")),
			)
			g.Return(jen.Id("c"))
		},
	))
	return b
}

func (b *BaseClientBuilder) Endpoint(iface options.Iface, ep options.Endpoint) ClientEndpointBuilder {
	b.codes = append(b.codes,
		jen.Const().Id(strcase.ToLowerCamel(ep.MethodName)+"ShortName").Op("=").Lit(ep.MethodShortName),
		jen.Const().Id(strcase.ToLowerCamel(ep.MethodName)+"FullName").Op("=").Lit(ep.MethodFullName),
	)
	return &clientEndpointBuilder{BaseClientBuilder: b, iface: iface, ep: ep, qualifier: b.qualifier}
}

func NewBaseClientBuilder(qualifier Qualifier) *BaseClientBuilder {
	return &BaseClientBuilder{qualifier: qualifier}
}
