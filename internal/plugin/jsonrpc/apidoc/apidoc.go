package apidoc

import (
	_ "embed"

	"strings"

	"github.com/555f/gg/internal/apidoc"
	"github.com/555f/gg/internal/plugin/jsonrpc/options"
	"github.com/555f/gg/pkg/file"
)

func Gen(title string, services []options.Iface) func(f *file.TxtFile) error {
	return func(f *file.TxtFile) error {
		api := apidoc.API{
			Title:   title,
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
					Title:       ep.Title,
					Description: ep.Description,
					Path:        ep.RPCMethodName,
					ContentType: "application/json",
					Method:      "POST",
				}

				var params []apidoc.Param

				for _, p := range ep.Params {
					apidoc.SchemaTypes(p.Type, api.Schemas, schemasNames)

					paramType, isArray := apidoc.ParamType(p.Type)
					title := strings.TrimSpace(p.Title)

					params = append(params, apidoc.Param{
						Name:     p.Name,
						Title:    title,
						Type:     paramType,
						Required: p.Required,
						Array:    isArray,
						In:       "body",
						Example:  apidoc.ExampleValue(p.Type),
					})
				}

				apiEndpoint.Params = append(apiEndpoint.Params,
					apidoc.Param{
						Name:    "jsonrpc",
						Type:    "string",
						In:      "body",
						Example: "2.0",
					},
					apidoc.Param{
						Name:    "id",
						Type:    "string",
						In:      "body",
						Example: "1",
					},
					apidoc.Param{
						Name:    "method",
						Type:    "string",
						In:      "body",
						Example: ep.RPCMethodName,
					},
					apidoc.Param{
						Name:   "params",
						Params: params,
						Type:   "object",
						In:     "body",
					},
				)

				var results []apidoc.Param
				for _, r := range ep.Results {
					apidoc.SchemaTypes(r.Type, api.Schemas, schemasNames)

					paramType, isArray := apidoc.ParamType(r.Type)
					title := strings.TrimSpace(r.Title)

					results = append(results, apidoc.Param{
						Name:    r.Name,
						Title:   title,
						Type:    paramType,
						Example: apidoc.ExampleValue(r.Type),
						Array:   isArray,
						In:      "body",
					})
				}
				apiEndpoint.Results = append(apiEndpoint.Results,
					apidoc.Param{
						Name:    "jsonrpc",
						Type:    "string",
						In:      "body",
						Example: "2.0",
					},
					apidoc.Param{
						Name:    "id",
						Type:    "string",
						In:      "body",
						Example: "1",
					},
					apidoc.Param{
						Name:   "result",
						Type:   "object",
						In:     "body",
						Params: results,
					},
				)
				apiService.Endpoints = append(apiService.Endpoints, apiEndpoint)
			}

			api.Services = append(api.Services, apiService)
		}
		return apidoc.Gen(title, api)(f)
	}
}
