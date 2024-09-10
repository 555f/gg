//go:build !gg

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"

	"github.com/555f/gg/examples/rest-service-chi/internal/config"
	"github.com/555f/gg/examples/rest-service-chi/internal/interface/controller"
	"github.com/555f/gg/examples/rest-service-chi/internal/server"

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

	cfg, err := config.New()
	if err != nil {
		slog.Error("config load", err)
		os.Exit(1)
	}

	r := chi.NewRouter()

	server.SetupRoutesProfileController(new(controller.ProfileController), r)

	r.Get("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port), r); err != nil {
		slog.Error("start server", err)
	}
}
