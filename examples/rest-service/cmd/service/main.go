package main

import (
	"fmt"

	"github.com/555f/gg/examples/rest-service/internal/config"
	"github.com/555f/gg/examples/rest-service/internal/interface/controller"
	"github.com/555f/gg/examples/rest-service/internal/server"

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
	cfg, errs := config.New()
	if len(errs) > 0 {
		for i := 0; i < len(errs); i++ {
			slog.Error("config load", errs[i])
		}
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
