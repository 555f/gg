package pkg

import (
	"github.com/555f/gg/cmd/gg/command"
	_ "github.com/555f/gg/internal/plugin/config"
	_ "github.com/555f/gg/internal/plugin/http"
	_ "github.com/555f/gg/internal/plugin/logging"
	_ "github.com/555f/gg/internal/plugin/middleware"
)

func Main(version string) {
	command.Execute(version)
}
