package rest

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
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
