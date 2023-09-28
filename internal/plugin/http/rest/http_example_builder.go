package rest

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/555f/curlbuilder"
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/strcase"
)

type HTTPExampleBuilder struct {
	iface options.Iface
	buf   bytes.Buffer
}

func (b *HTTPExampleBuilder) write(format string, a ...any) {
	fmt.Fprintf(&b.buf, format, a...)
}

func (b *HTTPExampleBuilder) Build() []byte {
	for _, ep := range b.iface.Endpoints {
		pathParts := strings.Split(ep.Path, "/")
		for i := 0; i < len(pathParts); i++ {
			s := pathParts[i]
			if strings.HasPrefix(s, ":") {
				name := s[1:]
				pathParts[i] = fmt.Sprintf("{{%s}}", name)
			}
		}

		uri := "{{scheme}}://{{host}}" + strings.Join(pathParts, "/")

		b.write("### %s - %s\n", ep.Title, ep.Description)
		for i := 0; i < len(ep.ParamsNameIdx); i++ {
			b.write("# @prompt %s\n", ep.ParamsNameIdx[i])
		}

		for _, h := range ep.OpenapiHeaders {
			b.write("# @prompt %s\n", strcase.ToLowerCamel(h.Name))
		}
		for _, h := range ep.HeaderParams {
			b.write("# @prompt %s\n", strcase.ToLowerCamel(h.Name))
		}
		switch b.iface.HTTPReq {
		case "http":
			b.write("%s %s HTTP/1.1\n", ep.HTTPMethod, uri)
			b.write("Content-Type: application/json")
			for _, h := range ep.OpenapiHeaders {
				b.write(h.Name + ": \"{{" + strcase.ToLowerCamel(h.Name) + "}}\"")
			}
			if len(ep.BodyParams) > 0 {
				b.write("\n" + ep.BodyParams.ToJSON())
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
				cb.SetBody(ep.BodyParams.ToJSON())
			}
			b.write(cb.String())
		}
		b.write("\n\n")
	}
	return b.buf.Bytes()
}

func NewHTTPExampleBuilder(iface options.Iface) *HTTPExampleBuilder {
	return &HTTPExampleBuilder{iface: iface}
}
