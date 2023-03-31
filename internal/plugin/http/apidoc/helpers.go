package apidoc

import (
	"bytes"
	"encoding/json"

	"github.com/555f/curlbuilder"
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/types"
)

func makeTypeRecursive(t any) (result interface{}) {
	switch t := t.(type) {
	case *types.Slice:
		return []interface{}{0}
	case *types.Basic:
		if t.IsSigned() || t.IsUnsigned() {
			return 123
		}
		if t.IsFloat() {
			return 1.0
		}
		return "abc"
	case *types.Struct:
		result := map[string]interface{}{}
		for _, field := range t.Fields {
			name := field.Var.Name
			if t, err := field.SysTags.Get("json"); err == nil {
				name = t.Value()
			}
			result[name] = makeTypeRecursive(field.Var.Type)
		}
		return result
	case *types.Named:
		return makeTypeRecursive(t.Type)
	}
	return
}

func paramType(t any) string {
	switch t := t.(type) {
	case *types.Basic:
		return t.Name
	case *types.Array:
		return paramType(t.Value) + " " + "array"
	case *types.Slice:
		return paramType(t.Value) + " " + "array"
	case *types.Map:
		return "object"
	case *types.Named:
		return t.Name
	}
	return ""
}

func requestCURL(ep *options.Endpoint) []byte {
	b := curlbuilder.New().
		SetURL("http://127.0.0.1" + ep.Path).
		SetMethod(ep.HTTPMethod)

	switch ep.HTTPMethod {
	case "POST", "PUT", "PATCH", "DELETE":
		if len(ep.BodyParams) > 0 {
			b.SetBody(requestJSON(ep))
		}
	}
	return []byte(b.String())
}

func responseJSON(ep *options.Endpoint) []byte {
	results := map[string]interface{}{}
	for _, result := range ep.BodyResults {
		results[result.Name] = makeTypeRecursive(result.Type)
	}
	var buf bytes.Buffer
	e := json.NewEncoder(&buf)
	e.SetIndent("", "  ")
	_ = e.Encode(results)
	return buf.Bytes()
}

func requestJSON(ep *options.Endpoint) []byte {
	results := map[string]interface{}{}
	for _, param := range ep.BodyParams {
		results[param.Name] = makeTypeRecursive(param.Type)
	}
	var buf bytes.Buffer
	e := json.NewEncoder(&buf)
	e.SetIndent("", "  ")
	_ = e.Encode(results)
	return bytes.Trim(buf.Bytes(), "\n")
}

func paramsByRequired(ep *options.Endpoint, required bool) (results []*options.EndpointParam) {
	for _, param := range ep.PathParams {
		if param.Required != required {
			continue
		}
		results = append(results, param)
	}
	for _, param := range ep.BodyParams {
		if param.Required != required {
			continue
		}
		results = append(results, param)
	}
	for _, param := range ep.QueryParams {
		if param.Required != required {
			continue
		}
		results = append(results, param)
	}
	for _, param := range ep.HeaderParams {
		if param.Required != required {
			continue
		}
		results = append(results, param)
	}
	for _, param := range ep.CookieParams {
		if param.Required != required {
			continue
		}
		results = append(results, param)
	}
	return
}

type contract struct {
	Name        string
	Title       string
	Description string
	Fields      []contractField
}

type contractField struct {
	Name  string
	Title string
	Type  string
}

func structTypesRecursive(t any, visited map[string]struct{}) (contracts []*contract) {
	if t, ok := t.(*types.Named); ok {
		if st, ok := t.Type.(*types.Struct); ok {
			_, ok := visited[t.Name]
			if ok {
				return
			}

			visited[t.Name] = struct{}{}

			c := &contract{
				Name:        t.Name,
				Title:       t.Title,
				Description: t.Description,
			}
			for _, field := range st.Fields {
				name := field.Var.Name
				if tag, err := field.SysTags.Get("json"); err == nil {
					name = tag.Value()
				}
				c.Fields = append(c.Fields, contractField{
					Name:  name,
					Title: field.Var.Title,
					Type:  paramType(field.Var.Type),
				})
				contracts = append(contracts, structTypesRecursive(field.Var.Type, visited)...)
			}
			contracts = append(contracts, c)
		}
	}
	return
}

func structTypes(ep *options.Endpoint) (contracts []*contract) {
	visited := map[string]struct{}{}
	for _, result := range ep.Results {
		contracts = append(contracts, structTypesRecursive(result.Type, visited)...)
	}
	return
}
