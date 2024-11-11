package http

import (
	_ "embed"
	"go/token"

	"os"
	"path"
	"path/filepath"

	"github.com/555f/gg/internal/openapi"
	"github.com/555f/gg/internal/plugin/http/apidoc"
	"github.com/555f/gg/internal/plugin/http/clienttest"
	"github.com/555f/gg/internal/plugin/http/httperror"
	"github.com/555f/gg/internal/plugin/http/openapidoc"
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/internal/plugin/http/rest"
	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/jaswdr/faker/v2"
	"gopkg.in/yaml.v2"

	"github.com/hashicorp/go-multierror"
)

const (
	promCollectorName = "prometheusCollector"
	prometheusPkg     = "github.com/prometheus/client_golang/prometheus"
	jsonPkg           = "encoding/json"
)

type Plugin struct {
	ctx *gg.Context
}

func (p *Plugin) Name() string { return "http" }

func (p *Plugin) Exec() (files []file.File, errs error) {
	serverOutput := filepath.Join(
		p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("server-output", "internal/server/server.go"),
	)
	clientOutput := filepath.Join(
		p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("client-output", "internal/server/client.go"),
	)
	clientTestOutput := filepath.Join(
		p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("client-test-output", "internal/server/client_test.go"),
	)
	openapiOutput := filepath.Join(
		p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("openapi-output", "docs"),
	)
	apiDocOutput := filepath.Join(
		p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("apidoc-output", "docs/api.html"),
	)

	apiDocTitle := p.ctx.Options.GetStringWithDefault("apidoc-title", "APIDoc")

	httpReqOutput := filepath.Join(
		p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("httpreq-output", ".http"),
	)
	errorWrapperPath := p.ctx.Options.GetString("error-wrapper")
	defaultErrorPath := p.ctx.Options.GetString("error-default")
	isCheckStrict := p.ctx.Options.GetBoolWithDefault("strict", true)

	var (
		serverServices     []options.Iface
		clientServices     []options.Iface
		clientTestServices []options.Iface
		openapiServices    []options.Iface
		apidocServices     []options.Iface
		errorWrapper       *options.ErrorWrapper
	)

	if errorWrapperPath != "" && defaultErrorPath != "" {
		errorWrapperPath = filepath.Join(p.ctx.PkgPath, errorWrapperPath)
		defaultErrorPath = filepath.Join(p.ctx.PkgPath, defaultErrorPath)

		var err error
		errorWrapper, err = options.DecodeErrorWrapper(errorWrapperPath, defaultErrorPath, p.ctx.Structs)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	for _, iface := range p.ctx.Interfaces {
		s, err := options.Decode(iface, isCheckStrict)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		if s.Openapi.Enable {
			openapiServices = append(openapiServices, s)
		}
		if s.APIDoc.Enable {
			apidocServices = append(apidocServices, s)
		}
		if s.Client.Enable {
			clientServices = append(clientServices, s)
		}
		if s.Client.EnableTest {
			clientTestServices = append(clientTestServices, s)
		}
		if s.Server.Enable {
			serverServices = append(serverServices, s)
		}
	}

	if len(apidocServices) > 0 {

		adFile := file.NewTxtFile(apiDocOutput)
		files = append(files, adFile)
		err := apidoc.Gen(apiDocTitle, openapiServices)(adFile)
		if err != nil {
			errs = multierror.Append(errs, errors.Error(err.Error(), token.Position{}))
		}
	}

	if len(openapiServices) > 0 {
		httpErrors := httperror.Load(p.ctx.Structs, errorWrapper)

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
			openapidoc.Gen(openAPI, s, httpErrors)(opFile)
		}
	}

	serverFile := file.NewGoFile(p.ctx.Module, serverOutput)
	serverFile.SetVersion(p.ctx.Version)

	clientFile := file.NewGoFile(p.ctx.Module, clientOutput)
	clientFile.SetVersion(p.ctx.Version)

	serverBuilder := rest.NewServerBuilder(serverFile, errorWrapper)
	clientBuilder := rest.NewBaseClientBuilder(clientFile)

	serverBuilder.RegisterHandlerStrategy("echo", func() rest.HandlerStrategy {
		return rest.NewHandlerStrategyEcho()
	})
	serverBuilder.RegisterHandlerStrategy("chi", func() rest.HandlerStrategy {
		return rest.NewHandlerStrategyChi()
	})

	if len(serverServices) > 0 {
		for _, iface := range serverServices {
			if iface.HTTPReq != "" {
				hrf := file.NewTxtFile(filepath.Join(httpReqOutput, strcase.ToSnake(iface.Name)+".http"))
				hrf.WriteBytes(
					rest.NewHTTPExampleBuilder(iface).Build(),
				)
				files = append(files, hrf)
			}
			serverBuilder.Controller(iface).Build()
		}
		serverFile.Add(serverBuilder.Build())
		files = append(files, serverFile)
	}

	if len(clientTestServices) > 0 {
		clientTestFile := file.NewGoFile(p.ctx.Module, clientTestOutput, file.UseTestPkg())

		fake := faker.New()
		clientTestGen := clienttest.New(clientTestFile.Group, p.ctx.PkgPath, fake, clientTestFile.Import, errorWrapper)

		for _, iface := range clientTestServices {
			for _, ep := range iface.Endpoints {
				clientTestGen.Generate(iface, ep, []clienttest.Config{
					{StatusCode: 200},
					{StatusCode: 400, CheckError: true},
				})
			}
		}

		files = append(files, clientTestFile)
	}

	if len(clientServices) > 0 {
		clientBuilder.
			SetErrorWrapper(errorWrapper).
			BuildTypes()

		for _, iface := range clientServices {
			clientBuilder.BuildStruct(iface)
			for _, ep := range iface.Endpoints {
				clientBuilder.Endpoint(iface, ep).
					BuildReqStruct().
					BuildSetters().
					BuildMethod().
					BuildReqMethod().
					BuildExecuteMethod()
			}
			clientBuilder.BuildConstruct(iface)
		}
		clientFile.Add(clientBuilder.Build())
		files = append(files, clientFile)
	}
	return
}

func (p *Plugin) Dependencies() []string {
	return []string{}
}
