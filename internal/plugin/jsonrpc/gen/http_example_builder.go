package gen

import (
	"bytes"
	"fmt"

	"github.com/555f/curlbuilder"
	"github.com/555f/gg/internal/plugin/jsonrpc/options"
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
		uri := "{{scheme}}://{{host}}"

		b.write("### %s - %s\n", ep.Title, ep.Description)
		//for i := 0; i < len(ep.ParamsNameIdx); i++ {
		//	b.write("# @prompt %s\n", ep.ParamsNameIdx[i])
		//}

		//for _, h := range ep.OpenapiHeaders {
		//	b.write("# @prompt %s\n", strcase.ToLowerCamel(h.Name))
		//}
		//for _, h := range ep.HeaderParams {
		//	b.write("# @prompt %s\n", strcase.ToLowerCamel(h.Name))
		//}
		switch b.iface.HTTPReq {
		case "http":
			b.write("POST %s HTTP/1.1\n", uri)
			b.write("Content-Type: application/json")
			//for _, h := range ep.OpenapiHeaders {
			//	b.write(h.Name + ": \"{{" + strcase.ToLowerCamel(h.Name) + "}}\"")
			//}
			if len(ep.Params) > 0 {
				b.write("\n" + ep.Params.ToJSON())
			}
		case "curl":
			headers := []string{"Content-Type", "application/json"}
			//for _, h := range ep.OpenapiHeaders {
			//	headers = append(headers, h.Name, "{{"+strcase.ToLowerCamel(h.Name)+"}}")
			//}
			//for _, h := range ep.HeaderParams {
			//	headers = append(headers, h.Name, "{{"+strcase.ToLowerCamel(h.Name)+"}}")
			//}
			cb := curlbuilder.New()
			cb.SetMethod("POST")
			cb.SetURL(uri)
			cb.SetHeaders(headers...)
			if len(ep.Params) > 0 {
				cb.SetBody(`{"jsonrpc":"2.0", "method":"` + ep.RPCMethodName + `", "params": ` + ep.Params.ToJSON() + `}`)
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
