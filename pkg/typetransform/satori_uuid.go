package typetransform

import (
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

var _ Parser = &SatoriUUIDTypeParse{}

type SatoriUUIDTypeParse struct{}

func (s *SatoriUUIDTypeParse) Parse(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (parseCode []jen.Code, paramID jen.Code, hasError bool) {
	named := t.(*types.Named)
	parseCode = append(parseCode, jen.List(assignID, jen.Err()).Op(op).Qual(named.Pkg.Path, "FromString").Call(valueID))
	return parseCode, assignID, true
}

func (s *SatoriUUIDTypeParse) Format(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (formatCode []jen.Code, paramID jen.Code, hasError bool) {
	return nil, jen.Add(valueID).Dot("String").Call(), false
}

func (s *SatoriUUIDTypeParse) Support(t any) bool {
	switch t := t.(type) {
	case *types.Named:
		return t.Pkg.Path == "github.com/satori/go.uuid" && t.Name == "UUID"
	}
	return false
}
