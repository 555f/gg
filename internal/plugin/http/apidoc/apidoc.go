package apidoc

import (
	_ "embed"
	"strings"

	"github.com/555f/gg/internal/apidoc"
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/file"
)

func Gen(title string, services []options.Iface) func(f *file.TxtFile) error {
	return func(f *file.TxtFile) error {
		api := apidoc.API{
			Schemas: map[string]*apidoc.Schema{},
		}
		for _, service := range services {
			schemasNames := make(map[string]struct{})

			apiService := apidoc.Service{
				Name:        service.Name,
				Title:       service.Title,
				Description: service.Description,
			}
			for _, ep := range service.Endpoints {
				apiEndpoint := apidoc.Endpoint{
					Title:  ep.Title,
					Path:   ep.Path,
					Method: ep.HTTPMethod,
				}
				for _, p := range ep.Params {
					apidoc.SchemaTypes(p.Type, api.Schemas, schemasNames)
					paramType, isArray := apidoc.ParamType(p.Type)
					apiEndpoint.Params = append(apiEndpoint.Params, apidoc.Param{
						Name:     p.Name,
						Title:    strings.TrimSpace(p.Title),
						Type:     paramType,
						Array:    isArray,
						Required: p.Required,
						In:       string(p.HTTPType),
						Example:  p.Zero,
					})
				}
				for _, r := range ep.Results {
					apidoc.SchemaTypes(r.Type, api.Schemas, schemasNames)
					paramType, isArray := apidoc.ParamType(r.Type)
					apiEndpoint.Results = append(apiEndpoint.Results, apidoc.Param{
						Name:  r.Name,
						Title: strings.TrimSpace(r.Title),
						Type:  paramType,
						Array: isArray,
						In:    string(r.HTTPType),
					})
				}
				apiService.Endpoints = append(apiService.Endpoints, apiEndpoint)
			}
			api.Services = append(api.Services, apiService)
		}
		return apidoc.Gen(title, api)(f)
	}
}
