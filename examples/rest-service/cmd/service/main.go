package main

import (
	"net/http"

	"github.com/f555/gg-examples/internal/interface/controller"
	"github.com/f555/gg-examples/internal/server"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	registry := prometheus.NewRegistry()

	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	t := server.NewRESTServer(
		server.ProfileController(
			new(controller.ProfileController),
			server.ProfileControllerMetricMiddleware(registry, nil),
		),
	)

	t.AddRoute("GET", "/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}), nil, nil, nil)

	_ = http.ListenAndServe(":8080", t)
}
