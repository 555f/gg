package prommetrics

import (
	"path/filepath"

	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gen"
	"github.com/555f/gg/pkg/gg"
	"github.com/dave/jennifer/jen"
)

var (
	promCollectorName = "PrometheusCollector"
	prometheusPkg     = "github.com/prometheus/client_golang/prometheus"
	// promhttpPkg       = "github.com/prometheus/client_golang/prometheus/promhttp"
	// httpPkg           = "net/http"
)

type Plugin struct {
	ctx *gg.Context
}

func (p *Plugin) Name() string { return "prommetrics" }

func (p *Plugin) Exec() (files []file.File, errs error) {
	output := filepath.Join(p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("output", "internal/metrics/prom_metrics.go"))

	f := file.NewGoFile(p.ctx.Module, output)

	labels := jen.Index().String().Values(
		jen.Lit("method"),
		jen.Lit("code"),
		jen.Lit("scopeName"),
		jen.Lit("methodNameFull"),
		jen.Lit("methodNameShort"),
		jen.Lit("err"),
	)

	f.Type().Id(promCollectorName).Struct(
		// jen.Id("inflight").Qual(prometheusPkg, "Gauge"),
		jen.Id("errRequests").Op("*").Qual(prometheusPkg, "CounterVec"),
		jen.Id("requests").Op("*").Qual(prometheusPkg, "CounterVec"),
		jen.Id("duration").Op("*").Qual(prometheusPkg, "HistogramVec"),
		// jen.Id("dnsDuration").Op("*").Qual(prometheusPkg, "HistogramVec"),
		// jen.Id("tlsDuration").Op("*").Qual(prometheusPkg, "HistogramVec"),
	)

	f.Func().Params(jen.Id("i").Op("*").Id(promCollectorName)).Id("Requests").Params().Params(jen.Op("*").Qual(prometheusPkg, "CounterVec")).Block(
		jen.Return(jen.Id("i").Dot("requests")),
	)

	f.Func().Params(jen.Id("i").Op("*").Id(promCollectorName)).Id("ErrRequests").Params().Params(jen.Op("*").Qual(prometheusPkg, "CounterVec")).Block(
		jen.Return(jen.Id("i").Dot("errRequests")),
	)

	f.Func().Params(jen.Id("i").Op("*").Id(promCollectorName)).Id("Duration").Params().Params(jen.Op("*").Qual(prometheusPkg, "HistogramVec")).Block(
		jen.Return(jen.Id("i").Dot("duration")),
	)

	// f.Func().Params(jen.Id("i").Op("*").Id(promCollectorName)).Id("DNSDuration").Params().Params(jen.Op("*").Qual(prometheusPkg, "HistogramVec")).Block(
	// 	jen.Return(jen.Id("i").Dot("dnsDuration")),
	// )

	// f.Func().Params(jen.Id("i").Op("*").Id(promCollectorName)).Id("TLSDuration").Params().Params(jen.Op("*").Qual(prometheusPkg, "HistogramVec")).Block(
	// 	jen.Return(jen.Id("i").Dot("tlsDuration")),
	// )

	f.Func().Params(jen.Id("i").Op("*").Id(promCollectorName)).Id("Describe").Params(jen.Id("in").Chan().Op("<-").Op("*").Qual(prometheusPkg, "Desc")).Block(
		// jen.Id("i").Dot("inflight").Dot("Describe").Call(jen.Id("in")),
		jen.Id("i").Dot("requests").Dot("Describe").Call(jen.Id("in")),
		jen.Id("i").Dot("errRequests").Dot("Describe").Call(jen.Id("in")),
		jen.Id("i").Dot("duration").Dot("Describe").Call(jen.Id("in")),
		// jen.Id("i").Dot("dnsDuration").Dot("Describe").Call(jen.Id("in")),
		// jen.Id("i").Dot("tlsDuration").Dot("Describe").Call(jen.Id("in")),
	)

	f.Func().Params(jen.Id("i").Op("*").Id(promCollectorName)).Id("Collect").Params(jen.Id("in").Chan().Op("<-").Qual(prometheusPkg, "Metric")).Block(
		// jen.Id("i").Dot("inflight").Dot("Collect").Call(jen.Id("in")),
		jen.Id("i").Dot("requests").Dot("Collect").Call(jen.Id("in")),
		jen.Id("i").Dot("errRequests").Dot("Collect").Call(jen.Id("in")),
		jen.Id("i").Dot("duration").Dot("Collect").Call(jen.Id("in")),
		// jen.Id("i").Dot("dnsDuration").Dot("Collect").Call(jen.Id("in")),
		// jen.Id("i").Dot("tlsDuration").Dot("Collect").Call(jen.Id("in")),
	)

	f.Func().Id("RegisterPrometheusCollector").Params(
		jen.Id("namespace").String(),
		jen.Id("subsystem").String(),
		jen.Id("reg").Qual(prometheusPkg, "Registerer"),
		jen.Id("constLabels").Map(jen.String()).String(),
	).Params(jen.Op("*").Id(promCollectorName), jen.Error()).Block(
		jen.Id("c").Op(":=").Op("&").Id(promCollectorName).Values(
			// jen.Id("inflight").Op(":").Qual(prometheusPkg, "NewGauge").Call(
			// 	jen.Qual(prometheusPkg, "GaugeOpts").Values(
			// 		jen.Id("Namespace").Op(":").Id("namespace"),
			// 		jen.Id("Subsystem").Op(":").Id("subsystem"),
			// 		jen.Id("Name").Op(":").Lit("in_flight_requests"),
			// 		jen.Id("Help").Op(":").Lit("A gauge of in-flight outgoing requests for the client."),
			// 		jen.Id("ConstLabels").Op(":").Id("constLabels"),
			// 	),
			// ),
			jen.Id("requests").Op(":").Qual(prometheusPkg, "NewCounterVec").Call(
				jen.Qual(prometheusPkg, "CounterOpts").Values(
					jen.Id("Namespace").Op(":").Id("namespace"),
					jen.Id("Subsystem").Op(":").Id("subsystem"),
					jen.Id("Name").Op(":").Lit("requests_total"),
					jen.Id("Help").Op(":").Lit("A counter for outgoing requests from the client."),
					jen.Id("ConstLabels").Op(":").Id("constLabels"),
				),
				labels,
			),
			jen.Id("errRequests").Op(":").Qual(prometheusPkg, "NewCounterVec").Call(
				jen.Qual(prometheusPkg, "CounterOpts").Values(
					jen.Id("Namespace").Op(":").Id("namespace"),
					jen.Id("Subsystem").Op(":").Id("subsystem"),
					jen.Id("Name").Op(":").Lit("err_requests_total"),
					jen.Id("Help").Op(":").Lit("A counter for outgoing error requests from the client."),
				),
				labels,
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
				labels,
			),
			// jen.Id("dnsDuration").Op(":").Qual(prometheusPkg, "NewHistogramVec").Call(
			// 	jen.Qual(prometheusPkg, "HistogramOpts").Values(
			// 		jen.Id("Namespace").Op(":").Id("namespace"),
			// 		jen.Id("Subsystem").Op(":").Id("subsystem"),
			// 		jen.Id("Name").Op(":").Lit("dns_duration_histogram_seconds"),
			// 		jen.Id("Help").Op(":").Lit("Trace dns latency histogram."),
			// 		jen.Id("Buckets").Op(":").Qual(prometheusPkg, "DefBuckets"),
			// 		jen.Id("ConstLabels").Op(":").Id("constLabels"),
			// 	),
			// 	labels,
			// ),
			// jen.Id("tlsDuration").Op(":").Qual(prometheusPkg, "NewHistogramVec").Call(
			// 	jen.Qual(prometheusPkg, "HistogramOpts").Values(
			// 		jen.Id("Namespace").Op(":").Id("namespace"),
			// 		jen.Id("Subsystem").Op(":").Id("subsystem"),
			// 		jen.Id("Name").Op(":").Lit("tls_duration_histogram_seconds"),
			// 		jen.Id("Help").Op(":").Lit("Trace tls latency histogram."),
			// 		jen.Id("Buckets").Op(":").Qual(prometheusPkg, "DefBuckets"),
			// 		jen.Id("ConstLabels").Op(":").Id("constLabels"),
			// 	),
			// 	labels,
			// ),
		),
		jen.Err().Op(":=").Id("reg").Dot("Register").Call(jen.Id("c")),
		jen.Do(gen.CheckErr(jen.Return(
			jen.Nil(),
			jen.Err(),
		))),
		jen.Return(jen.Id("c"), jen.Nil()),
	)
	return []file.File{f}, errs
}

func (p *Plugin) Dependencies() []string { return nil }
