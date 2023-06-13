package rest

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/file"
)

func GenMetric(s options.Iface) func(f *file.GoFile) {
	return func(f *file.GoFile) {
		//prometheusPkg := "github.com/prometheus/client_golang/prometheus"
		//
		//f.Func().Id(s.Name+"MetricMiddleware").Params(
		//	Id("registry").Qual(prometheusPkg, "Registerer"),
		//	Id("buckets").Index().Float64(),
		//).Id(s.Name + "ServerOption").Block(
		//	Return(Id(s.Name + "ApplyOptions").CallFunc(func(g *Group) {
		//		for _, endpoint := range s.Endpoints {
		//			g.Id(s.Name + endpoint.MethodName + "ServerMiddlewareFunc").Call(
		//				Id("baseMetricMiddleware").Call(Id("registry"), Id("buckets"), Lit(s.Name+"."+endpoint.MethodName)),
		//			)
		//		}
		//	})),
		//)
	}
}
