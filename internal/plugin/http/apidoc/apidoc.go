package apidoc

import (
	"bytes"
	"text/template"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"

	"github.com/555f/gg/internal/plugin/http/httperror"
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/file"
)

type templateValues struct {
	Title      string
	Services   []options.Iface
	HTTPErrors []httperror.Error
	Style      string
	VueJS      string
}

func Gen(apiDocTemplate, styleCSS, vueJS string, services []options.Iface, httpErrors []httperror.Error) func(f *file.TxtFile) error {
	return func(f *file.TxtFile) error {
		t := template.Must(template.New("apidoc").
			Funcs(template.FuncMap{
				"requestJSON":      requestJSON,
				"responseJSON":     responseJSON,
				"requestCURL":      requestCURL,
				"paramsByRequired": paramsByRequired,
				"paramType":        paramType,
				"structTypes":      structTypes,
				"highlight": func(source []byte, lexer string) string {
					var result bytes.Buffer
					l := lexers.Get(lexer)
					l = chroma.Coalesce(l)
					f := html.New(
						html.WithClasses(true),
						html.Standalone(false),
					)
					it, _ := l.Tokenise(nil, string(source))
					_ = f.Format(&result, styles.Fallback, it)
					return result.String()
				},
			}).Parse(apiDocTemplate))
		return t.Execute(f, &templateValues{
			Services:   services,
			HTTPErrors: httpErrors,
			Style:      styleCSS,
			VueJS:      vueJS,
		})
	}
}
