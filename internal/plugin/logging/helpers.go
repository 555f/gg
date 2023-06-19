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

func makeParamLog(p *types.Var) *jen.Statement {
	st := jen.Empty()
	switch t := p.Type.(type) {
	default:
		st.Id(p.Name)
	case *types.Named:
		if hasMethodString(t) {
			st.Id(p.Name).Dot("String").Call()
		} else {
			st.Lit("")
		}
	case *types.Slice, *types.Array, *types.Map:
		st.Len(jen.Id(p.Name))
	}
	return st
}
