package logging

import (
	"github.com/555f/gg/pkg/types"

	"github.com/dave/jennifer/jen"
)

func hasMethodString(v *types.Named) bool {
	for _, method := range v.Methods {
		if method.Name != "String" {
			continue
		}
		if len(method.Sig.Params) == 0 && len(method.Sig.Results) == 1 {
			if t, ok := method.Sig.Results[0].Type.(*types.Basic); ok {
				return t.IsString()
			}
		}
	}
	return false
}

func makeLog(name string, t interface{}) *jen.Statement {
	st := jen.Lit(name).Op(",")
	switch t := t.(type) {
	default:
		return nil
	case *types.Basic:
		st.Id(name)
	case *types.Named:
		if hasMethodString(t) {
			st.Dot("String").Call()
		} else {
			return nil
		}
	case *types.Slice, *types.Array, *types.Map:
		st.Len(jen.Id(name))
	}
	return st
}
