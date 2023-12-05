package jsonrpc

import (
	"path/filepath"

	"github.com/555f/gg/internal/plugin/jsonrpc/gen"
	"github.com/555f/gg/internal/plugin/jsonrpc/options"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
	"github.com/hashicorp/go-multierror"
)

type Plugin struct {
	ctx *gg.Context
}

// Dependencies implements gg.Plugin.
func (*Plugin) Dependencies() []string {
	return nil
}

// Exec implements gg.Plugin.
func (p *Plugin) Exec() (files []file.File, errs error) {
	serverOutput := filepath.Join(
		p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("server-output", "internal/server/server.go"),
	)
	clientOutput := filepath.Join(
		p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("client-output", "internal/server/client.go"),
	)
	// openapiOutput := filepath.Join(
	// 	p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("openapi-output", "docs/openapi.yaml"),
	// )
	// apiDocOutput := filepath.Join(
	// 	p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("apidoc-output", "docs/apidoc.html"),
	// )
	var (
		serverServices  []options.Iface
		clientServices  []options.Iface
		openapiServices []options.Iface
	)
	for _, iface := range p.ctx.Interfaces {
		s, err := options.Decode(iface)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		if s.Openapi.Enable {
			openapiServices = append(openapiServices, s)
		}
		if s.Client.Enable {
			clientServices = append(clientServices, s)
		}
		if s.Server.Enable {
			serverServices = append(serverServices, s)
		}
	}

	serverFile := file.NewGoFile(p.ctx.Module, serverOutput)
	serverFile.SetVersion(p.ctx.Version)

	clientFile := file.NewGoFile(p.ctx.Module, clientOutput)
	clientFile.SetVersion(p.ctx.Version)

	serverBuilder := gen.NewServerBuilder(serverFile)
	clientBuilder := gen.NewBaseClientBuilder(clientFile)

	if len(serverServices) > 0 {
		serverBuilder.RegisterHandlerStrategy("default", func() gen.HandlerStrategy {
			return gen.NewHandlerStrategyJSONRPC()
		})

		for _, iface := range serverServices {
			controllerBuilder := serverBuilder.Controller(iface)

			controllerBuilder.BuildHandlers()

			for _, ep := range iface.Endpoints {
				controllerBuilder.Endpoint(ep).BuildReqStruct().
					BuildReqDec().
					BuildRespStruct().
					Build()
			}
		}

		files = append(files, serverFile)
		serverFile.Add(serverBuilder.Build())
	}

	if len(clientServices) > 0 {
		for _, iface := range clientServices {
			clientBuilder.BuildStruct(iface)
			for _, ep := range iface.Endpoints {
				clientBuilder.Endpoint(iface, ep).
					BuildReqStruct().
					BuildSetters().
					BuildMethod().
					BuildReqMethod().
					BuildResultMethod().
					BuildExecuteMethod()
			}
			clientBuilder.BuildConstruct(iface)
		}
		clientFile.Add(clientBuilder.Build())
		files = append(files, clientFile)
	}
	return
}

// Name implements gg.Plugin.
func (*Plugin) Name() string {
	return "jsonrpc"
}
