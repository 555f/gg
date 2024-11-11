package typetransform

import (
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

var _ Parser = &GoogleUUIDTypeParse{}

type GoogleUUIDTypeParse struct{}

func (s *GoogleUUIDTypeParse) Parse(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (parseCode []jen.Code, paramID jen.Code, hasError bool) {
	named := t.(*types.Named)
	parseCode = append(parseCode, jen.List(assignID, jen.Err()).Op(op).Qual(named.Pkg.Path, "Parse").Call(valueID))
	return parseCode, assignID, true
}

func (s *GoogleUUIDTypeParse) Format(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (formatCode []jen.Code, paramID jen.Code, hasError bool) {
	return
}

func (s *GoogleUUIDTypeParse) Support(t any) bool {
	switch t := t.(type) {
	case *types.Named:
		return t.Pkg.Path == "github.com/google/uuid" && t.Name == "UUID"
	}
	return false
}
