package webview

import (
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
)

func makeBindName(iface *gg.Interface, method *types.Func) string {
	return strcase.ToLowerCamel(iface.Named.Name) + "_" + method.Name
}

func methodByType(t any) string {
	if t, ok := t.(*types.Basic); ok {
		switch {
		case t.IsBool():
			return "Bool"
		case t.IsInteger():
			return "Int"
		case t.IsFloat():
			return "Float"
		case t.IsString():
			return "String"
		}
	}
	return ""
}

var indexNames = [7]string{"i", "j", "k", "l", "m", "n", "k"}

func makeIndexName(i int) string {
	if i > 6 {
		panic("index names max 7")
	}
	return indexNames[i]
}
