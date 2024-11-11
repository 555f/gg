package typetransform

import (
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

var _ Parser = &StringTypeParse{}

type StringTypeParse struct{}

func (s *StringTypeParse) Parse(valueID, assignID jen.Code, op string, t any, qualFunc types.QualFunc) (parseCode []jen.Code, paramID jen.Code, hasError bool) {
	return nil, valueID, false
}

func (s *StringTypeParse) Format(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (formatCode []jen.Code, paramID jen.Code, hasError bool) {
	return nil, valueID, false
}

func (s *StringTypeParse) Support(t any) bool {
	switch t := t.(type) {
	case *types.Basic:
		return t.Name() == "string"
	}
	return false
}
