package openapi

import (
	"bytes"
	"net/http"

	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"

	"gopkg.in/yaml.v2"
)

type Builder struct {
	openAPI OpenAPI
}

func (b *Builder) OpenAPI() OpenAPI {
	return b.openAPI
}

func (b *Builder) AddServer(url, description string, variables map[string]Variable) *Builder {
	b.openAPI.Servers = append(b.openAPI.Servers, Server{
		URL:         url,
		Description: description,
		Variables:   variables,
	})
	return b
}

func (b *Builder) AddTag(name, description string, externalDocs ExternalDocs) *Builder {
	b.openAPI.Tags = append(b.openAPI.Tags, Tag{
		Name:         name,
		Description:  description,
		ExternalDocs: externalDocs,
	})
	return b
}

func (b *Builder) AddScheme(scheme string) *Builder {
	b.openAPI.Schemes = append(b.openAPI.Schemes, scheme)
	return b
}

func (b *Builder) AddPath(method, path string, operation *Operation, opts ...PathOption) *Builder {
	opPath, ok := b.openAPI.Paths[path]
	if !ok {
		opPath = &Path{}
		b.openAPI.Paths[path] = opPath
	}
	for _, opt := range opts {
		opt(opPath)
	}
	switch method {
	case http.MethodGet:
		opPath.Get = operation
	case http.MethodPost:
		opPath.Post = operation
	case http.MethodPut:
		opPath.Put = operation
	case http.MethodPatch:
		opPath.Patch = operation
	case http.MethodDelete:
		opPath.Delete = operation
	}
	return b
}

func (b *Builder) makeRef(named *types.Named) string {
	return "#/components/schemas/" + strcase.ToCamel(named.Name)
}

func (b *Builder) schemaByTypeRecursive(title, description string, schema *Schema, t interface{}) {
	switch t := t.(type) {
	case *types.Named:
		switch t.Pkg.Path {
		default:
			refSchema := &Schema{
				Properties:  map[string]*Schema{},
				Description: description,
			}

			if _, ok := b.openAPI.Components.Schemas[t.Name]; !ok {
				b.AddComponent(t.Name, refSchema)

				b.schemaByTypeRecursive(title, description, refSchema, t.Type)

			}

			schema.Ref = b.makeRef(t)

			return
		case "encoding/json":
			schema.Description = title
			schema.Type = "object"
			schema.Properties = Properties{}
			return
		case "time":
			switch t.Name {
			case "Duration":
				schema.Description = title
				schema.Type = "string"
				schema.Example = "1h3m30s"
			case "Time":
				schema.Description = title
				schema.Type = "string"
				schema.Format = "date-time"
				schema.Example = "1985-04-02T01:30:00.00Z"
			}
			return
		case "gopkg.in/guregu/null.v4":
			switch t.Name {
			case "String":
				schema.Description = title
				schema.Type = "string"
			case "Int":
				schema.Description = title
				schema.Type = "number"
			case "Float":
				schema.Description = title
				schema.Type = "float"
			case "Bool":
				schema.Description = title
				schema.Type = "bool"
			case "Time":
				schema.Description = title
				schema.Type = "string"
				schema.Format = "date-time"
				schema.Example = "1985-04-02T01:30:00.00Z"
			}
			return
		case "github.com/pborman/uuid", "github.com/google/uuid":
			schema.Description = title
			schema.Type = "string"
			schema.Format = "uuid"
			schema.Example = "d5c02d83-6fbc-4dd7-8416-9f85ed80de46"
			return
		}
	case *types.Struct:
		for _, field := range t.Fields {
			name := field.Var.Name
			if tag, err := field.SysTags.Get("json"); err == nil {
				name = tag.Name
			}
			if name == "-" {
				continue
			}
			filedSchema := &Schema{
				Properties: Properties{},
			}
			schema.Properties[name] = filedSchema
			b.schemaByTypeRecursive(field.Var.Title, "", filedSchema, field.Var.Type)
		}
	case *types.Map:
		mapSchema := &Schema{
			Properties: Properties{},
		}
		schema.Description = title
		schema.Properties = Properties{"key": mapSchema}
		b.schemaByTypeRecursive(title, description, mapSchema, t.Value)
		return
	case *types.Array:
		schema.Description = title
		schema.Type = "array"
		schema.Items = &Schema{
			Properties: Properties{},
		}
		b.schemaByTypeRecursive(title, description, schema.Items, t.Value)
		return
	case *types.Slice:
		schema.Description = description
		if tp, ok := t.Value.(*types.Basic); ok && tp.IsByte() {
			schema.Type = "string"
			schema.Format = "byte"
			schema.Example = "U3dhZ2dlciByb2Nrcw=="
		} else {
			schema.Type = "array"
			schema.Items = &Schema{
				Properties: Properties{},
			}
			b.schemaByTypeRecursive(title, description, schema.Items, t.Value)
		}
		return
	case *types.Interface:
		schema.Description = title
		schema.Type = "object"
		schema.Description = "Can be any value - string, number, boolean, array or object."
		schema.Properties = Properties{}
		schema.Example = "null"
		schema.AnyOf = []Schema{
			{Type: "string", Example: "abc"},
			{Type: "integer", Example: 1},
			{Type: "number", Format: "float", Example: 1.11},
			{Type: "boolean", Example: true},
			{Type: "array"},
			{Type: "object"},
		}
		return
	case *types.Basic:
		if t.IsString() {
			schema.Description = title
			schema.Type = "string"
			schema.Example = "abc"
			return
		}
		if t.IsBool() {
			schema.Description = title
			schema.Type = "boolean"
			schema.Example = "true"
		}
		if t.IsNumeric() {
			if t.IsInt32() || t.IsUint32() {
				schema.Description = title
				schema.Type = "integer"
				schema.Format = "int32"
				schema.Example = 1
				return
			}
			if t.IsInt64() || t.IsUint64() {
				schema.Description = title
				schema.Type = "integer"
				schema.Format = "int64"
				schema.Example = 1
				return
			}
			if t.IsFloat32() || t.IsFloat64() {
				schema.Description = title
				schema.Type = "number"
				schema.Format = "float"
				schema.Example = 1.11
				return
			}
			schema.Description = title
			schema.Type = "integer"
			schema.Example = 1
			return
		}
	}
}

func (b *Builder) SchemaByType(title, description string, t interface{}) (schema *Schema) {
	schema = &Schema{
		Properties: Properties{},
	}
	b.schemaByTypeRecursive(title, description, schema, t)
	return
}

func (b *Builder) AddComponent(name string, schema *Schema) *Builder {
	if b.openAPI.Components.Schemas == nil {
		b.openAPI.Components.Schemas = map[string]*Schema{}
	}
	if _, ok := b.openAPI.Components.Schemas[name]; !ok {
		b.openAPI.Components.Schemas[name] = schema
	}
	return b
}

func (b *Builder) Build() ([]byte, error) {
	var buf bytes.Buffer
	if err := yaml.NewEncoder(&buf).Encode(b.openAPI); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func NewBuilder(openAPI OpenAPI) *Builder {
	if openAPI.Paths == nil {
		openAPI.Paths = map[string]*Path{}
	}
	return &Builder{openAPI: openAPI}
}
