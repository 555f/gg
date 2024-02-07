package main

import (
	"fmt"
	"os"

	"github.com/555f/gg/examples/rest-service-echo/internal/config"
	"github.com/555f/gg/examples/rest-service-echo/internal/interface/controller"
	"github.com/555f/gg/examples/rest-service-echo/internal/server"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"golang.org/x/exp/slog"
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

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	server.SetupRoutesProfileController(new(controller.ProfileController), e)

	e.GET("/metrics", echo.WrapHandler(promhttp.HandlerFor(registry, promhttp.HandlerOpts{})))

	if err := e.Start(fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port)); err != nil {
		slog.Error("start server", err)
	}
}
