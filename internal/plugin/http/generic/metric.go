package generic

import (
	"github.com/555f/gg/pkg/file"

	. "github.com/dave/jennifer/jen"
)

func GenMetric() func(f *file.GoFile) {
	return func(f *file.GoFile) {
		prometheusPkg := "github.com/prometheus/client_golang/prometheus"
		promautoPkg := "github.com/prometheus/client_golang/prometheus/promauto"
		promhttpPkg := "github.com/prometheus/client_golang/prometheus/promhttp"
		httpPkg := "net/http"

		f.Func().Id("baseMetricMiddleware").Params(
			Id("registry").Qual(prometheusPkg, "Registerer"),
			Id("buckets").Index().Float64(),
			Id("handlerName").String(),
		).Params(
			Func().Params(Qual("net/http", "Handler")).Params(Qual("net/http", "Handler")),
		).Block(
			Id("reg").Op(":=").Qual(prometheusPkg, "WrapRegistererWith").Call(
				Qual(prometheusPkg, "Labels").Values(
					Lit("handler").Op(":").Id("handlerName"),
				),
				Id("registry"),
			),
			Id("requestsTotal").Op(":=").Qual(promautoPkg, "With").Call(
				Id("reg"),
			).Dot("NewCounterVec").Call(
				Qual(prometheusPkg, "CounterOpts").Values(
					Id("Name").Op(":").Lit("http_requests_total"),
					Id("Help").Op(":").Lit("Tracks the number of HTTP requests."),
				),
				Index().String().Values(Lit("method"), Lit("code")),
			),
			Id("requestDuration").Op(":=").Qual(promautoPkg, "With").Call(
				Id("reg"),
			).Dot("NewHistogramVec").Call(
				Qual(prometheusPkg, "HistogramOpts").Values(
					Id("Name").Op(":").Lit("http_request_duration_seconds"),
					Id("Help").Op(":").Lit("Tracks the latencies for HTTP requests."),
				),
				Index().String().Values(Lit("method"), Lit("code")),
			),
			Id("requestSize").Op(":=").Qual(promautoPkg, "With").Call(
				Id("reg"),
			).Dot("NewSummaryVec").Call(
				Qual(prometheusPkg, "SummaryOpts").Values(
					Id("Name").Op(":").Lit("http_request_size_bytes"),
					Id("Help").Op(":").Lit("Tracks the size of HTTP requests."),
				),
				Index().String().Values(Lit("method"), Lit("code")),
			),
			Id("responseSize").Op(":=").Qual(promautoPkg, "With").Call(
				Id("reg"),
			).Dot("NewSummaryVec").Call(
				Qual(prometheusPkg, "SummaryOpts").Values(
					Id("Name").Op(":").Lit("http_response_size_bytes"),
					Id("Help").Op(":").Lit("Tracks the size of HTTP responses."),
				),
				Index().String().Values(Lit("method"), Lit("code")),
			),
			Return(
				Func().Params(Id("next").Qual(httpPkg, "Handler")).Params(Qual(httpPkg, "Handler")).Block(
					Return(
						Qual(promhttpPkg, "InstrumentHandlerCounter").Call(Id("requestsTotal"),
							Qual(promhttpPkg, "InstrumentHandlerDuration").Call(Id("requestDuration"),
								Qual(promhttpPkg, "InstrumentHandlerRequestSize").Call(Id("requestSize"),
									Qual(promhttpPkg, "InstrumentHandlerRequestSize").Call(Id("responseSize"), Id("next")),
								),
							),
						),
					),
				),
			),
		)
	}
}
