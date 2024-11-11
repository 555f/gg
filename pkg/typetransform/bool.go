package typetransform

import (
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

var _ Parser = &BoolTypeParse{}

type BoolTypeParse struct{}

func (s *BoolTypeParse) Parse(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (parseCode []jen.Code, paramID jen.Code, hasError bool) {
	parseCode = []jen.Code{jen.List(assignID, jen.Err()).Op(op).Do(qualFn("github.com/555f/go-strings", "ParseBool")).Types(types.Convert(t, qualFn)).Call(valueID)}
	return parseCode, assignID, true
}

func (s *BoolTypeParse) Format(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (formatCode []jen.Code, paramID jen.Code, hasError bool) {
	return nil, jen.Qual("strconv", "FormatBool").Call(valueID), false
}

func (s *BoolTypeParse) Support(t any) bool {
	switch t := t.(type) {
	case *types.Basic:
		return t.Name() == "bool"
	}
	return false
}
