package http

import (
	_ "embed"
	"go/token"
	"os"
	"path"
	"path/filepath"

	"github.com/555f/gg/internal/openapi"
	"github.com/555f/gg/internal/plugin/http/apidoc"
	"github.com/555f/gg/internal/plugin/http/generic"
	"github.com/555f/gg/internal/plugin/http/httperror"
	"github.com/555f/gg/internal/plugin/http/jsonrpc"
	"github.com/555f/gg/internal/plugin/http/openapidoc"
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/internal/plugin/http/rest"
	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v3"
)

//go:embed files/apidoc.html
var apiDocTemplate string

//go:embed files/style.css
var styleCSS string

//go:embed files/vue.min.js
var vueJS string

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
	openapiOutput := filepath.Join(
		p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("openapi-output", "docs/openapi.yaml"),
	)
	apiDocOutput := filepath.Join(
		p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("apidoc-output", "docs/apidoc.html"),
	)
	httpReqOutput := filepath.Join(
		p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("httpreq-output", ".http"),
	)
	errorWrapperPath := p.ctx.Options.GetString("error-wrapper")
	defaultErrorPath := p.ctx.Options.GetString("error-default")

	var (
		serverServices  []options.Iface
		clientServices  []options.Iface
		openapiServices []options.Iface
		apidocServices  []options.Iface
		errorWrapper    *options.ErrorWrapper
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
		s, err := options.Decode(iface)
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
		if s.Server.Enable {
			serverServices = append(serverServices, s)
		}
	}

	if len(apidocServices) > 0 || len(openapiServices) > 0 {
		httpErrors := httperror.Load(p.ctx.Structs, errorWrapper)

		if len(apidocServices) > 0 {
			adFile := file.NewTxtFile(apiDocOutput)
			files = append(files, adFile)
			err := apidoc.Gen(apiDocTemplate, styleCSS, vueJS, openapiServices, httpErrors)(adFile)
			if err != nil {
				errs = multierror.Append(errs, errors.Error(err.Error(), token.Position{}))
			}
		}
		if len(openapiServices) > 0 {
			var openAPI openapi.OpenAPI
			openapiTmpl := p.ctx.Options.GetString("openapi-tpl")
			if openapiTmpl != "" {
				openapiTmplPath := path.Join(p.ctx.Workdir, openapiTmpl)
				data, err := os.ReadFile(openapiTmplPath)
				if err != nil {
					errs = multierror.Append(errs, errors.Error(err.Error(), token.Position{}))
				}
				if err := yaml.Unmarshal(data, &openAPI); err != nil {
					errs = multierror.Append(errs, errors.Error(err.Error(), token.Position{}))
				}
			}

			opFile := file.NewTxtFile(openapiOutput)
			files = append(files, opFile)
			openapidoc.Gen(openAPI, openapiServices, httpErrors)(opFile)
		}
	}

	fileSet := map[string]*file.GoFile{}
	if len(serverServices) > 0 {
		pkgPathVisited := map[string]struct{}{}
		for _, s := range serverServices {
			if s.HTTPReq != "" {
				hrf := file.NewTxtFile(filepath.Join(httpReqOutput, strcase.ToSnake(s.Name)+".http"))
				generic.GenHTTPReq(s)(hrf)
				files = append(files, hrf)
			}

			f, ok := fileSet[serverOutput]
			if !ok {
				f = file.NewGoFile(p.ctx.Module, serverOutput)

				rest.GenTypes()(f)
				rest.GenErrorEncoder(errorWrapper)(f)
				rest.GenHTTPHandler()(f)

				fileSet[serverOutput] = f
				files = append(files, f)
			}

			pkgPath := path.Dir(serverOutput) + s.Name

			if _, ok := pkgPathVisited[pkgPath]; !ok {
				switch s.Type {
				case "rest":
					rest.GenOptions(s)(f)
				}
				pkgPathVisited[pkgPath] = struct{}{}
			}

			rest.GenStruct(s)(f)
			rest.GenMetric(s)(f)
		}
	}

	if len(clientServices) > 0 {
		pkgPathVisited := map[string]struct{}{}
		for _, s := range clientServices {
			f, ok := fileSet[clientOutput]
			if !ok {
				f = file.NewGoFile(p.ctx.Module, clientOutput)
				fileSet[clientOutput] = f
				files = append(files, f)
			}
			pkgPath := path.Dir(clientOutput)
			if _, ok := pkgPathVisited[pkgPath]; !ok {
				switch s.Type {
				case "rest":
					rest.GenClientTypes()(f)
				}
				pkgPathVisited[pkgPath] = struct{}{}
			}
			switch s.Type {
			case "rest":
				rest.GenClient(s, errorWrapper)(f)
			case "jsonrpc":
				jsonrpc.GenClient(s)(f)
			}
		}
	}
	return
}

func (p *Plugin) Dependencies() []string {
	return []string{}
}
