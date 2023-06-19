package openapidoc

import (
	openapi2 "github.com/555f/gg/internal/openapi"

	"github.com/555f/gg/internal/plugin/http/httperror"
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/file"
)

func Gen(openAPI openapi2.OpenAPI, services []options.Iface, httpErrors []httperror.Error) func(f *file.TxtFile) {
	return func(f *file.TxtFile) {
		f.WriteText("# OpenAPI docs\n")
		b := openapi2.NewBuilder(openAPI)
		httpErrorsMap := make(map[string]httperror.Error, len(httpErrors))
		for _, e := range httpErrors {
			httpErrorsMap[e.Name] = e
			errorProps := openapi2.Properties{}
			for _, field := range e.Fields {
				errorProps[field.Name] = b.SchemaByType("", "", field.FldType)
			}
			b.AddComponent("Error"+e.Name, &openapi2.Schema{
				Type:       "object",
				Properties: errorProps,
			})
		}
		for _, service := range services {
			httpErrors := service.Server.Errors
			for _, ep := range service.Endpoints {
				tags := ep.OpenapiTags
				if len(tags) == 0 {
					tags = service.Openapi.Tags
				}
				o := &openapi2.Operation{
					Tags:        tags,
					Summary:     ep.Title,
					Description: ep.Description,
					OperationID: "",
					Consumes:    nil,
					Produces:    nil,
					Parameters:  nil,
					RequestBody: nil,
					Responses:   openapi2.Responses{},
					Security:    nil,
				}
				for _, errorName := range append(httpErrors, ep.Errors...) {
					if e, ok := httpErrorsMap[errorName]; ok {
						o.Responses[e.StatusCode] = &openapi2.Response{
							Description: e.Title + " " + e.Description,
							Content:     openapi2.Content{},
						}
						for _, bodyType := range ep.AcceptTypes {
							var contentType string
							switch bodyType {
							case "json":
								contentType = "application/json"
							case "xml":
								contentType = "application/xml"
							}
							o.Responses[e.StatusCode].Content[contentType] = openapi2.Media{Schema: &openapi2.Schema{
								Ref: "#/components/schemas/Error" + e.Name,
							}}
						}
					}
				}
				switch ep.HTTPMethod {
				case "POST", "PUT", "DELETE", "PATCH":
					if len(ep.BodyParams) > 0 {
						requestProperties := openapi2.Properties{}
						for _, param := range ep.BodyParams {
							requestProperties[param.Name] = b.SchemaByType(param.Title, "", param.Type)
						}
						o.RequestBody = &openapi2.RequestBody{
							Required: true,
							Content:  map[string]openapi2.Media{},
						}
						for _, bodyType := range ep.ContentTypes {
							switch bodyType {
							case "json":
								o.RequestBody.Content["application/json"] = openapi2.Media{
									Schema: &openapi2.Schema{
										Type:       "object",
										Properties: requestProperties,
									},
								}
							case "xml":
								o.RequestBody.Content["application/xml"] = openapi2.Media{
									Schema: &openapi2.Schema{
										Type: "object",
										XML: &openapi2.XML{
											Name: ep.ReqRootXMLName,
										},
										Properties: requestProperties,
									},
								}
							}
						}
					}
				}
				o.Responses[200] = &openapi2.Response{
					Description: "OK",
					Content:     openapi2.Content{},
				}
				if len(ep.BodyResults) > 0 {
					responseSchema := &openapi2.Schema{
						Type:       "object",
						Properties: openapi2.Properties{},
					}
					if len(ep.BodyResults) == 1 && ep.NoWrapResponse {
						responseSchema = b.SchemaByType(ep.BodyResults[0].Title, "", ep.BodyResults[0].Type)
					} else {
						for _, result := range ep.BodyResults {
							responseSchema.Properties[result.Name] = b.SchemaByType(result.Title, "", result.Type)
						}
					}
					for _, bodyType := range ep.AcceptTypes {
						switch bodyType {
						//case "custom":
						//	contentType, _ = resultContentType(ep.BodyResults)
						case "json":
							o.Responses[200].Content["application/json"] = openapi2.Media{
								Schema: responseSchema,
							}
						case "xml":
							responseSchema.XML = &openapi2.XML{Name: ep.RespRootXMLName}
							o.Responses[200].Content["application/xml"] = openapi2.Media{
								Schema: responseSchema,
							}
						}
					}
				}
				for _, param := range ep.QueryParams {
					o.Parameters = append(o.Parameters, openapi2.Parameter{
						In:          "query",
						Name:        param.Name,
						Description: param.Title,
						Required:    param.Required,
						Schema:      b.SchemaByType(param.Title, "", param.Type),
					})
				}
				for _, param := range ep.PathParams {
					o.Parameters = append(o.Parameters, openapi2.Parameter{
						In:          "path",
						Name:        param.Name,
						Description: param.Title,
						Required:    param.Required,
						Schema:      b.SchemaByType(param.Title, "", param.Type),
					})
				}
				for _, param := range ep.HeaderParams {
					o.Parameters = append(o.Parameters, openapi2.Parameter{
						In:          "header",
						Name:        param.Name,
						Description: param.Title,
						Required:    param.Required,
						Schema:      b.SchemaByType(param.Title, "", param.Type),
					})
				}
				for _, h := range ep.OpenapiHeaders {
					o.Parameters = append(o.Parameters, openapi2.Parameter{
						In:          "header",
						Name:        h.Name,
						Description: h.Title,
						// Required:    param.Required,
					})
				}

				b.AddPath(ep.HTTPMethod, ep.Path, o)
			}
		}
		p, _ := b.Build()
		f.WriteBytes(p)
	}
}
