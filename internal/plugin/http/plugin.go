package http

import (
	_ "embed"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"

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

type Plugin struct{}

func (p *Plugin) Name() string { return "http" }

func (p *Plugin) Exec(ctx *gg.Context) (files []file.File, errs error) {
	serverOutput := ctx.Options.GetStringWithDefault("server-output", "internal/server/server.go")
	clientOutput := ctx.Options.GetStringWithDefault("client-output", "internal/server/client.go")
	openapiOutput := ctx.Options.GetStringWithDefault("openapi-output", "docs/openapi.yaml")
	apiDocOutput := ctx.Options.GetStringWithDefault("apidoc-output", "docs/apidoc.html")
	httpReqOutput := ctx.Options.GetStringWithDefault("httpreq-output", ".http")

	serverOutput = filepath.Join(ctx.Module.Dir, serverOutput)
	clientOutput = filepath.Join(ctx.Module.Dir, clientOutput)
	openapiOutput = filepath.Join(ctx.Module.Dir, openapiOutput)
	apiDocOutput = filepath.Join(ctx.Module.Dir, apiDocOutput)
	httpReqOutput = filepath.Join(ctx.Module.Dir, httpReqOutput)

	var (
		serverServices  []options.Iface
		clientServices  []options.Iface
		openapiServices []options.Iface
		apidocServices  []options.Iface
		errorWrapper    *options.ErrorWrapper
	)

	errorWrapperPath := ctx.Options.GetString("error-wrapper")
	defaultErrorPath := ctx.Options.GetString("error-default")

	if errorWrapperPath != "" && defaultErrorPath != "" {
		var err error
		defaultErrorPath = strings.Replace(defaultErrorPath, "~", ctx.Module.Path, 1)
		errorWrapperPath = strings.Replace(errorWrapperPath, "~", ctx.Module.Path, 1)
		errorWrapper, err = options.DecodeErrorWrapper(errorWrapperPath, defaultErrorPath, ctx.Structs)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	for _, iface := range ctx.Interfaces {
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
		httpErrors := httperror.Load(ctx.Structs, errorWrapper)

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
			openapiTmpl := ctx.Options.GetString("openapi-tpl")
			if openapiTmpl != "" {
				openapiTmplPath := path.Join(ctx.Module.Dir, openapiTmpl)
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
				f = file.NewGoFile(ctx.Module, serverOutput)

				fileSet[serverOutput] = f
				files = append(files, f)
			}

			pkgPath := path.Dir(serverOutput)

			if _, ok := pkgPathVisited[pkgPath]; !ok {
				rest.GenTypes()(f)
				rest.GenErrorEncoder(errorWrapper)(f)
				rest.GenHTTPHandler()(f)
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
				f = file.NewGoFile(ctx.Module, clientOutput)
				fileSet[clientOutput] = f
				files = append(files, f)
			}
			pkgPath := path.Dir(clientOutput)
			if _, ok := pkgPathVisited[pkgPath]; !ok {
				switch s.Type {
				case "rest":
					generic.GenRESTClient()(f)
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
