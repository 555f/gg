package webview

import (
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
)

func makeBindName(iface *gg.Interface, method *types.Func) string {
	return strcase.ToLowerCamel(iface.Named.Name) + "_" + method.Name
}
