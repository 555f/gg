package jsonrpc

import (
	"go/token"
	"os"
	"path"
	"path/filepath"

	"github.com/555f/gg/internal/openapi"
	"github.com/555f/gg/internal/plugin/jsonrpc/apidoc"
	"github.com/555f/gg/internal/plugin/jsonrpc/gen"
	"github.com/555f/gg/internal/plugin/jsonrpc/openapidoc"
	"github.com/555f/gg/internal/plugin/jsonrpc/options"
	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v2"
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
	apiDocOutput := filepath.Join(
		p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("apidoc-output", "docs/api.html"),
	)
	apiDocTitle := p.ctx.Options.GetStringWithDefault("apidoc-title", "APIDoc")
	openapiOutput := filepath.Join(
		p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("openapi-output", "docs"),
	)
	httpReqOutput := filepath.Join(
		p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("httpreq-output", "docs"),
	)

	var (
		serverServices  []options.Iface
		clientServices  []options.Iface
		openapiServices []options.Iface
		apidocServices  []options.Iface
	)
	for _, iface := range p.ctx.Interfaces {
		s, err := options.Decode(iface)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		if s.APIDoc.Enable {
			apidocServices = append(apidocServices, s)
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

	if len(apidocServices) > 0 {
		adFile := file.NewTxtFile(apiDocOutput)
		files = append(files, adFile)
		err := apidoc.Gen(apiDocTitle, apidocServices)(adFile)
		if err != nil {
			errs = multierror.Append(errs, errors.Error(err.Error(), token.Position{}))
		}
	}

	if len(openapiServices) > 0 {
		openapiTmpl := p.ctx.Options.GetString("openapi-tpl")
		for _, s := range openapiServices {
			var openAPI openapi.OpenAPI
			if openapiTmpl != "" {
				openapiTmplPath := path.Join(p.ctx.Workdir, openapiTmpl, strcase.ToSnake(s.Name)+"_openapi.tpl.yaml")
				data, err := os.ReadFile(openapiTmplPath)
				if err != nil {
					errs = multierror.Append(errs, errors.Error(err.Error(), token.Position{}))
				}
				if err := yaml.Unmarshal(data, &openAPI); err != nil {
					errs = multierror.Append(errs, errors.Error(err.Error(), token.Position{}))
				}
			}
			opFile := file.NewTxtFile(filepath.Join(openapiOutput, strcase.ToSnake(s.Name)+".yaml"))
			files = append(files, opFile)
			openapidoc.Gen(openAPI, s)(opFile)
		}
	}

	serverFile := file.NewGoFile(p.ctx.Module, serverOutput)
	serverFile.SetVersion(p.ctx.Version)

	clientFile := file.NewGoFile(p.ctx.Module, clientOutput)
	clientFile.SetVersion(p.ctx.Version)

	serverBuilder := gen.NewServerBuilder(serverFile)
	clientBuilder := gen.NewBaseClientBuilder(clientFile)

	if len(serverServices) > 0 {
		for _, iface := range serverServices {
			if iface.HTTPReq != "" {
				hrf := file.NewTxtFile(filepath.Join(httpReqOutput, strcase.ToSnake(iface.Name)+".http"))
				hrf.WriteBytes(
					gen.NewHTTPExampleBuilder(iface).Build(),
				)
				files = append(files, hrf)
			}
		}
		serverFile.Add(serverBuilder.Build())
		files = append(files, serverFile)

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
