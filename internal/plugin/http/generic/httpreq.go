package generic

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/555f/curlbuilder"
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
)

func jsonFromType(name string, t any) string {
	var buf bytes.Buffer

	buf.WriteString("  " + strconv.Quote(name) + ":")

	switch v := t.(type) {
	case *types.Named:
		name := v.Pkg.Path + "." + v.Name
		switch name {
		default:
			if st := v.Struct(); st != nil {
				buf.WriteString("{\n")
				for _, f := range v.Struct().Fields {
					name := f.Var.Name
					if t, err := f.SysTags.Get("json"); err == nil {
						name = t.Value()
					}
					buf.WriteString("  " + jsonFromType(name, f.Var.Type))
				}
				buf.WriteString("  }")
			} else {
				buf.WriteString("\"\"")
			}
		case "time.Time":
			buf.WriteString("\"\"")
		}

	case *types.Slice, *types.Array:
		buf.WriteString("[]")
	case *types.Basic:
		if v.IsNumeric() {
			buf.WriteString("0")
		} else {
			buf.WriteString("\"\"")
		}
	}
	return buf.String()
}

func jsonFromParams(params []*options.EndpointParam) string {
	var buf bytes.Buffer
	buf.WriteString("{\n")
	for i, p := range params {
		if i > 0 {
			buf.WriteString(",\n")
		}
		buf.WriteString(jsonFromType(p.Name, p.Type))
	}
	buf.WriteString("}\n")
	return buf.String()
}

func GenHTTPReq(s options.Iface) func(f *file.TxtFile) {
	return func(f *file.TxtFile) {
		for _, ep := range s.Endpoints {
			pathParts := strings.Split(ep.Path, "/")
			for i := 0; i < len(pathParts); i++ {
				s := pathParts[i]
				if strings.HasPrefix(s, ":") {
					name := s[1:]
					pathParts[i] = fmt.Sprintf("{{%s}}", name)
				}
			}

			uri := "{{scheme}}://{{host}}" + strings.Join(pathParts, "/")

			f.WriteText("### %s - %s\n", ep.Title, ep.Description)
			for i := 0; i < len(ep.ParamsNameIdx); i++ {
				f.WriteText("# @prompt %s\n", ep.ParamsNameIdx[i])
			}

			for _, h := range ep.OpenapiHeaders {
				f.WriteText("# @prompt %s\n", strcase.ToLowerCamel(h.Name))
			}
			for _, h := range ep.HeaderParams {
				f.WriteText("# @prompt %s\n", strcase.ToLowerCamel(h.Name))
			}

			switch s.HTTPReq {
			case "http":
				f.WriteText("%s %s HTTP/1.1\n", ep.HTTPMethod, uri)
				f.WriteText("Content-Type: application/json")
				for _, h := range ep.OpenapiHeaders {
					f.WriteText(h.Name + ": \"{{" + strcase.ToLowerCamel(h.Name) + "}}\"")
				}
				if len(ep.BodyParams) > 0 {
					f.WriteText("\n" + jsonFromParams(ep.BodyParams))
				}
			case "curl":
				headers := []string{"Content-Type", "application/json"}
				for _, h := range ep.OpenapiHeaders {
					headers = append(headers, h.Name, "{{"+strcase.ToLowerCamel(h.Name)+"}}")
				}
				for _, h := range ep.HeaderParams {
					headers = append(headers, h.Name, "{{"+strcase.ToLowerCamel(h.Name)+"}}")
				}
				cb := curlbuilder.New()
				cb.SetMethod(ep.HTTPMethod)
				cb.SetURL(uri)
				cb.SetHeaders(headers...)
				if len(ep.BodyParams) > 0 {
					cb.SetBody(jsonFromParams(ep.BodyParams))
				}
				f.WriteText(cb.String())
			}
			f.WriteText("\n\n")
		}
	}
}
