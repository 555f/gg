package pkg

import (
	"github.com/555f/gg/cmd/gg/command"
	_ "github.com/555f/gg/internal/plugin/cli"
	_ "github.com/555f/gg/internal/plugin/config"
	_ "github.com/555f/gg/internal/plugin/grpc"
	_ "github.com/555f/gg/internal/plugin/http"
	_ "github.com/555f/gg/internal/plugin/jsonrpc"
	_ "github.com/555f/gg/internal/plugin/klog"
	_ "github.com/555f/gg/internal/plugin/metrics"
	_ "github.com/555f/gg/internal/plugin/middleware"
	_ "github.com/555f/gg/internal/plugin/prommetrics"
	_ "github.com/555f/gg/internal/plugin/pwa"
	_ "github.com/555f/gg/internal/plugin/slog"
	_ "github.com/555f/gg/internal/plugin/webview"
)

func Main(version string) {
	command.Execute(version)
}
