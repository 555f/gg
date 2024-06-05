package openapidoc

import (
	"github.com/555f/gg/internal/openapi"
	openapi2 "github.com/555f/gg/internal/openapi"

	"github.com/555f/gg/internal/plugin/jsonrpc/options"
	"github.com/555f/gg/pkg/file"
)

type jsonRPCError struct {
	httpCode string
	code     int
	message  string
}

var jsonRPCStandartErrors = []jsonRPCError{
	{

		httpCode: "x-32600",
		code:     -32600,
		message:  "Invalid Request",
	},
	{
		httpCode: "x-32601",
		code:     -32601,
		message:  "Method not found",
	},
	{
		httpCode: "x-32602",
		code:     -32602,
		message:  "Invalid params",
	},
	{
		httpCode: "x-32603",
		code:     -32603,
		message:  "Internal error",
	},
}

func makeJSONRpcError(code int, message string) *openapi.Schema {
	return &openapi2.Schema{
		Type: "object",
		Properties: openapi2.Properties{
			"code": &openapi2.Schema{
				Type:    "integer",
				Example: code,
			},
			"message": &openapi2.Schema{
				Type:    "string",
				Example: message,
			},
			"data": &openapi2.Schema{
				Type: "object",
			},
		},
	}
}

func BuildOpenAPI(openAPI openapi2.OpenAPI, service options.Iface) *openapi2.Builder {
	b := openapi2.NewBuilder(openAPI)

	tags := service.Openapi.Tags
	for _, ep := range service.Endpoints {
		o := &openapi2.Operation{
			Tags:        append(tags, ep.OpenapiTags...),
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
		if len(ep.Params) > 0 {
			requestProperties := openapi2.Properties{}
			for _, param := range ep.Params {
				requestProperties[param.Name] = b.SchemaByType(param.Title, "", param.Type)
			}
			o.RequestBody = &openapi2.RequestBody{
				Required: true,
				Content:  map[string]openapi2.Media{},
			}

			o.RequestBody.Content["application/json"] = openapi2.Media{
				Schema: &openapi2.Schema{
					Type: "object",
					Properties: openapi2.Properties{
						"id": &openapi2.Schema{
							OneOf: []openapi2.Schema{
								{Type: "string"},
								{Type: "number"},
							},
						},
						"jsonrpc": &openapi2.Schema{
							Type:    "string",
							Example: "2.0",
						},
						"method": &openapi2.Schema{
							Type:    "string",
							Example: ep.RPCMethodName,
						},
						"params": &openapi2.Schema{
							Type:       "object",
							Properties: requestProperties,
						},
					},
				},
			}
		}

		o.Responses["200"] = &openapi2.Response{
			Description: "OK",
			Content:     openapi2.Content{},
		}
		// if len(ep.Results) > 0 {
		responseSchema := &openapi2.Schema{
			Type:       "object",
			Properties: openapi2.Properties{},
		}
		if len(ep.Results) == 1 {
			responseSchema = b.SchemaByType(ep.Results[0].Title, "", ep.Results[0].Type)
		} else {
			for _, result := range ep.Results {
				responseSchema.Properties[result.Name] = b.SchemaByType(result.Title, "", result.Type)
			}
		}
		o.Responses["200"].Content["application/json"] = openapi2.Media{
			Schema: &openapi2.Schema{
				Type: "object",
				Properties: openapi2.Properties{
					"id": &openapi2.Schema{
						OneOf: []openapi2.Schema{
							{Type: "string"},
							{Type: "number"},
						},
					},
					"jsonrpc": &openapi2.Schema{
						Type:    "string",
						Example: "2.0",
					},
					"result": responseSchema,
				},
			},
		}

		for _, jsonRPCError := range jsonRPCStandartErrors {
			o.Responses[jsonRPCError.httpCode] = &openapi2.Response{
				Content: openapi2.Content{
					"application/json": openapi2.Media{
						Schema: &openapi2.Schema{
							Type: "object",
							Properties: openapi2.Properties{
								"id": &openapi2.Schema{
									OneOf: []openapi2.Schema{
										{Type: "string"},
										{Type: "number"},
									},
								},
								"jsonrpc": &openapi2.Schema{
									Type:    "string",
									Example: "2.0",
								},
								"error": makeJSONRpcError(jsonRPCError.code, jsonRPCError.message),
							},
						},
					},
				},
			}
		}
		b.AddPath("POST", "/"+ep.MethodName, o)
	}
	return b
}

func Gen(openAPI openapi2.OpenAPI, service options.Iface) func(f *file.TxtFile) {
	return func(f *file.TxtFile) {
		f.WriteText("# OpenAPI docs\n")

		b := BuildOpenAPI(openAPI, service)

		p, _ := b.Build()
		f.WriteBytes(p)
	}
}
