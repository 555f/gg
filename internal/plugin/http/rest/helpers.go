package rest

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

func clientStructName(iface options.Iface) string {
	return iface.Name + "Client"
}

func isNamedType(t any) (ok bool) {
	_, ok = t.(*types.Named)
	return
}

func normalizeVarName(name string) string {
	return strcase.ToLowerCamel(name)
}

func wrapResponse(names []string, completeFn func(g *jen.Group), qualFunc types.QualFunc) func(g *jen.Group) {
	return func(g *jen.Group) {
		if len(names) > 0 {
			g.Id(strcase.ToCamel(names[0])).StructFunc(wrapResponse(names[1:], completeFn, qualFunc)).Tag(map[string]string{"json": names[0]})
		} else {
			completeFn(g)
		}
	}
}
