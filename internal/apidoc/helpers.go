package apidoc

import (
	"strings"

	"github.com/555f/gg/pkg/types"
)

func ParamType(t any) (v string, isArray bool) {
	switch t := t.(type) {
	case *types.Basic:
		switch {
		case t.IsSigned():
			return "int", false
		case t.IsUnsigned():
			return "uint", false
		case t.IsString():
			return "string", false
		case t.IsFloat():
			return "float", false
		case t.IsBool():
			return "bool", false
		}
	case *types.Array:
		v, _ = ParamType(t.Value)
		isArray = true
		return
	case *types.Slice:
		isArray = true
		v, _ = ParamType(t.Value)
		return
	case *types.Map:
		return "object", false
	case *types.Named:
		return "$" + t.Name, false
	}
	return "", false
}

type Schema struct {
	Name        string  `json:"name"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Params      []Param `json:"params"`
}

func SchemaTypes(t any, schemas map[string]*Schema, schemasNames map[string]struct{}) {
	if s, ok := t.(*types.Slice); ok {
		t = s.Value
	}
	if t, ok := t.(*types.Named); ok {
		if st, ok := t.Type.(*types.Struct); ok {
			_, ok := schemas[t.Name]
			if ok {
				return
			}

			title := t.Title
			if title == "" {
				title = t.Name
			}

			schema := &Schema{
				Name:        t.Name,
				Title:       title,
				Description: t.Description,
			}
			for _, field := range st.Fields {
				name := field.Name
				if tag, err := field.SysTags.Get("json"); err == nil {
					name = tag.Value()
				}
				paramType, isArray := ParamType(field.Type)
				title := strings.TrimSpace(field.Title)
				schema.Params = append(schema.Params, Param{
					Name:     name,
					Title:    title,
					Type:     paramType,
					In:       "",
					Required: false,
					Array:    isArray,
					Example:  ExampleValue(field.Type),
				})

				schemasNames[schema.Name] = struct{}{}
				schemas[t.Name] = schema

				SchemaTypes(field.Type, schemas, schemasNames)
			}
		}
	}
}

func ExampleValue(t any) string {
	switch t := t.(type) {
	case *types.Basic:
		switch {
		case t.IsSigned():
			return "-1"
		case t.IsUnsigned():
			return "1"
		case t.IsString():
			return "abc"
		case t.IsFloat():
			return "1.0"
		case t.IsBool():
			return "true"
		}
	case *types.Array:
		return ExampleValue(t.Value)
	case *types.Slice:
		return ExampleValue(t.Value)

	case *types.Map:
		return "{}"
	case *types.Named:
		switch t.Pkg.Path {
		case "github.com/google/uuid", "github.com/satori/go.uuid":
			switch t.Name {
			case "UUID":
				return `"25200f0a-c0e0-4b7f-bd36-904e26a4456e"`
			}
		}
		return "{}"
	}
	return ""
}
