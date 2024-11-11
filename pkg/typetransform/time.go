package typetransform

import (
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

var _ Parser = &TimeTypeParse{}

type TimeTypeParse struct{}

func (s *TimeTypeParse) Parse(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (parseCode []jen.Code, paramID jen.Code, hasError bool) {
	n := t.(*types.Named)
	switch n.Name {
	case "Time":
		parseCode = []jen.Code{jen.List(assignID, jen.Err()).Op(op).Do(qualFn("time", "Parse")).Call(jen.Do(qualFn("time", "RFC3339")), valueID)}
	case "Duration":
		parseCode = []jen.Code{jen.List(assignID, jen.Err()).Op(op).Do(qualFn("time", "ParseDuration")).Call(valueID)}
	}
	return parseCode, assignID, true
}

func (s *TimeTypeParse) Format(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (formatCode []jen.Code, paramID jen.Code, hasError bool) {
	return nil, jen.Add(valueID).Dot("Format").Call(jen.Do(qualFn("time", "RFC3339"))), false
}

func (s *TimeTypeParse) Support(t any) bool {
	switch t := t.(type) {
	case *types.Named:
		return t.Pkg.Path == "time"
	}
	return false
}
