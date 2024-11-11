package typetransform

import (
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

var _ Parser = &FloatTypeParse{}

type FloatTypeParse struct{}

func (s *FloatTypeParse) Parse(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (parseCode []jen.Code, paramID jen.Code, hasError bool) {
	b := t.(*types.Basic)
	parseCode = []jen.Code{jen.List(assignID, jen.Err()).Op(op).Do(qualFn("github.com/555f/go-strings", "ParseFloat")).Types(types.Convert(t, qualFn)).Call(valueID, jen.Lit(10), jen.Lit(b.BitSize()))}
	return parseCode, assignID, true
}

func (s *FloatTypeParse) Format(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (formatCode []jen.Code, paramID jen.Code, hasError bool) {
	return
}

func (s *FloatTypeParse) Support(t any) bool {
	switch t := t.(type) {
	case *types.Basic:
		switch t.Name() {
		case "float32", "float64":
			return true
		}
	}
	return false
}
