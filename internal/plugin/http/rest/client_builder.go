package rest

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/gen"
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
	b.codes = append(b.codes,
		jen.Type().Id("contextKey").String(),
		jen.Const().Id("methodContextKey").Id("contextKey").Op("=").Lit("method"),
		jen.Const().Id("shortMethodContextKey").Id("contextKey").Op("=").Lit("shortMethod"),
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
						jen.List(jen.Id("resp"), jen.Err()).Op(":=").Id("next").Dot("RoundTrip").Call(jen.Id("r")),
						jen.If(jen.Id("err").Op("!=").Nil()).Block(
							jen.Id("labels").Op(":=").Qual(prometheusPkg, "Labels").Values(
								jen.Lit("method").Op(":").Qual(stringsPkg, "ToLower").Call(jen.Id("r").Dot("Method")),
							),
							jen.List(jen.Id("labels").Index(jen.Lit("methodNameFull")), jen.Id("_")).Op("=").Id("r").Dot("Context").Call().Dot("Value").Call(jen.Id("methodContextKey")).Assert(jen.String()),
							jen.List(jen.Id("labels").Index(jen.Lit("methodNameShort")), jen.Id("_")).Op("=").Id("r").Dot("Context").Call().Dot("Value").Call(jen.Id("shortMethodContextKey")).Assert(jen.String()),

							jen.Id("errType").Op(":=").Lit(""),
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
						),
						jen.Return(jen.Id("resp"), jen.Err()),
					),
			),
		),

		jen.Type().Id("outgoingInstrumentation").Struct(
			jen.Id("inflight").Qual(prometheusPkg, "Gauge"),
			jen.Id("errRequests").Op("*").Qual(prometheusPkg, "CounterVec"),
			jen.Id("requests").Op("*").Qual(prometheusPkg, "CounterVec"),
			jen.Id("duration").Op("*").Qual(prometheusPkg, "HistogramVec"),
			jen.Id("dnsDuration").Op("*").Qual(prometheusPkg, "HistogramVec"),
			jen.Id("tlsDuration").Op("*").Qual(prometheusPkg, "HistogramVec"),
		),
		jen.Func().Params(jen.Id("i").Op("*").Id("outgoingInstrumentation")).Id("Describe").Params(jen.Id("in").Chan().Op("<-").Op("*").Qual(prometheusPkg, "Desc")).Block(
			jen.Id("i").Dot("inflight").Dot("Describe").Call(jen.Id("in")),
			jen.Id("i").Dot("requests").Dot("Describe").Call(jen.Id("in")),
			jen.Id("i").Dot("errRequests").Dot("Describe").Call(jen.Id("in")),
			jen.Id("i").Dot("duration").Dot("Describe").Call(jen.Id("in")),
			jen.Id("i").Dot("dnsDuration").Dot("Describe").Call(jen.Id("in")),
			jen.Id("i").Dot("tlsDuration").Dot("Describe").Call(jen.Id("in")),
		),

		jen.Func().Params(jen.Id("i").Op("*").Id("outgoingInstrumentation")).Id("Collect").Params(jen.Id("in").Chan().Op("<-").Qual(prometheusPkg, "Metric")).Block(
			jen.Id("i").Dot("inflight").Dot("Collect").Call(jen.Id("in")),
			jen.Id("i").Dot("requests").Dot("Collect").Call(jen.Id("in")),
			jen.Id("i").Dot("errRequests").Dot("Collect").Call(jen.Id("in")),
			jen.Id("i").Dot("duration").Dot("Collect").Call(jen.Id("in")),
			jen.Id("i").Dot("dnsDuration").Dot("Collect").Call(jen.Id("in")),
			jen.Id("i").Dot("tlsDuration").Dot("Collect").Call(jen.Id("in")),
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
		jen.Func().Id("WithProm").Params(
			jen.Id("namespace").String(),
			jen.Id("subsystem").String(),
			jen.Id("reg").Qual(prometheusPkg, "Registerer"),
			jen.Id("constLabels").Map(jen.String()).String(),
		).Id("ClientOption").Block(
			jen.Return(jen.Func().Params(jen.Id("o").Op("*").Id(clientOptionName)).Block(

				jen.Id("i").Op(":=").Op("&").Id("outgoingInstrumentation").Values(
					jen.Id("inflight").Op(":").Qual(prometheusPkg, "NewGauge").Call(
						jen.Qual(prometheusPkg, "GaugeOpts").Values(
							jen.Id("Namespace").Op(":").Id("namespace"),
							jen.Id("Subsystem").Op(":").Id("subsystem"),
							jen.Id("Name").Op(":").Lit("in_flight_requests"),
							jen.Id("Help").Op(":").Lit("A gauge of in-flight outgoing requests for the client."),
							jen.Id("ConstLabels").Op(":").Id("constLabels"),
						),
					),
					jen.Id("requests").Op(":").Qual(prometheusPkg, "NewCounterVec").Call(
						jen.Qual(prometheusPkg, "CounterOpts").Values(
							jen.Id("Namespace").Op(":").Id("namespace"),
							jen.Id("Subsystem").Op(":").Id("subsystem"),
							jen.Id("Name").Op(":").Lit("requests_total"),
							jen.Id("Help").Op(":").Lit("A counter for outgoing requests from the client."),
							jen.Id("ConstLabels").Op(":").Id("constLabels"),
						),
						jen.Index().String().Values(jen.Lit("method"), jen.Lit("code"), jen.Lit("methodNameFull"), jen.Lit("methodNameShort")),
					),
					jen.Id("errRequests").Op(":").Qual(prometheusPkg, "NewCounterVec").Call(
						jen.Qual(prometheusPkg, "CounterOpts").Values(
							jen.Id("Namespace").Op(":").Id("namespace"),
							jen.Id("Subsystem").Op(":").Id("subsystem"),
							jen.Id("Name").Op(":").Lit("err_requests_total"),
							jen.Id("Help").Op(":").Lit("A counter for outgoing error requests from the client."),
						),
						jen.Index().String().Values(jen.Lit("method"), jen.Lit("err"), jen.Lit("methodNameFull"), jen.Lit("methodNameShort")),
					),
					jen.Id("duration").Op(":").Qual(prometheusPkg, "NewHistogramVec").Call(
						jen.Qual(prometheusPkg, "HistogramOpts").Values(
							jen.Id("Namespace").Op(":").Id("namespace"),
							jen.Id("Subsystem").Op(":").Id("subsystem"),
							jen.Id("Name").Op(":").Lit("request_duration_histogram_seconds"),
							jen.Id("Help").Op(":").Lit("A histogram of outgoing request latencies."),
							jen.Id("Buckets").Op(":").Qual(prometheusPkg, "DefBuckets"),
							jen.Id("ConstLabels").Op(":").Id("constLabels"),
						),
						jen.Index().String().Values(jen.Lit("method"), jen.Lit("code"), jen.Lit("methodNameFull"), jen.Lit("methodNameShort")),
					),
					jen.Id("dnsDuration").Op(":").Qual(prometheusPkg, "NewHistogramVec").Call(
						jen.Qual(prometheusPkg, "HistogramOpts").Values(
							jen.Id("Namespace").Op(":").Id("namespace"),
							jen.Id("Subsystem").Op(":").Id("subsystem"),
							jen.Id("Name").Op(":").Lit("dns_duration_histogram_seconds"),
							jen.Id("Help").Op(":").Lit("Trace dns latency histogram."),
							jen.Id("Buckets").Op(":").Qual(prometheusPkg, "DefBuckets"),
							jen.Id("ConstLabels").Op(":").Id("constLabels"),
						),
						jen.Index().String().Values(jen.Lit("method"), jen.Lit("code"), jen.Lit("methodNameFull"), jen.Lit("methodNameShort")),
					),
					jen.Id("tlsDuration").Op(":").Qual(prometheusPkg, "NewHistogramVec").Call(
						jen.Qual(prometheusPkg, "HistogramOpts").Values(
							jen.Id("Namespace").Op(":").Id("namespace"),
							jen.Id("Subsystem").Op(":").Id("subsystem"),
							jen.Id("Name").Op(":").Lit("tls_duration_histogram_seconds"),
							jen.Id("Help").Op(":").Lit("Trace tls latency histogram."),
							jen.Id("Buckets").Op(":").Qual(prometheusPkg, "DefBuckets"),
							jen.Id("ConstLabels").Op(":").Id("constLabels"),
						),
						jen.Index().String().Values(jen.Lit("method"), jen.Lit("code"), jen.Lit("methodNameFull"), jen.Lit("methodNameShort")),
					),
				),
				jen.Id("trace").Op(":=").Op("&").Qual(promhttpPkg, "InstrumentTrace").Values(),
				jen.Id("o").Dot("client").Dot("Transport").Op("=").
					Id("instrumentRoundTripperErrCounter").Call(jen.Id("i").Dot("errRequests"),
					jen.Qual(promhttpPkg, "InstrumentRoundTripperInFlight").Call(
						jen.Id("i").Dot("inflight"),
						jen.Qual(promhttpPkg, "InstrumentRoundTripperCounter").Call(
							jen.Id("i").Dot("requests"),
							jen.Qual(promhttpPkg, "InstrumentRoundTripperTrace").Call(
								jen.Id("trace"),
								jen.Qual(promhttpPkg, "InstrumentRoundTripperDuration").Call(
									jen.Id("i").Dot("duration"),
									jen.Qual(httpPkg, "DefaultTransport"),
									jen.Id("labelFromContext").Call(jen.Lit("methodNameShort"), jen.Id("shortMethodContextKey")),
									jen.Id("labelFromContext").Call(jen.Lit("methodNameFull"), jen.Id("methodContextKey")),
								),
							),
							jen.Id("labelFromContext").Call(jen.Lit("methodNameShort"), jen.Id("shortMethodContextKey")),
							jen.Id("labelFromContext").Call(jen.Lit("methodNameFull"), jen.Id("methodContextKey")),
						),
					),
				),

				jen.Err().Op(":=").Id("reg").Dot("Register").Call(jen.Id("i")),
				jen.Do(gen.CheckErr(jen.Panic(jen.Err()))),
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
	b.codes = append(b.codes, jen.Type().Id(clientName).StructFunc(func(g *jen.Group) {
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
					jen.Id("client").Op(":").Qual("net/http", "DefaultClient"),
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
	return &clientEndpointBuilder{BaseClientBuilder: b, iface: iface, ep: ep, qualifier: b.qualifier}
}

func NewBaseClientBuilder(qualifier Qualifier) *BaseClientBuilder {
	return &BaseClientBuilder{qualifier: qualifier}
}
