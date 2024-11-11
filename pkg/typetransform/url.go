package typetransform

import (
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

var _ Parser = &URLTypeParse{}

type URLTypeParse struct{}

func (s *URLTypeParse) Parse(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (parseCode []jen.Code, paramID jen.Code, hasError bool) {
	parseCode = []jen.Code{jen.List(assignID, jen.Err()).Op(op).Do(qualFn("net/url", "Parse")).Call(valueID)}
	return parseCode, assignID, true
}

func (s *URLTypeParse) Format(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (formatCode []jen.Code, paramID jen.Code, hasError bool) {
	return
}

func (s *URLTypeParse) Support(t any) bool {
	switch t := t.(type) {
	case *types.Named:
		return t.Pkg.Path == "net/url" && t.Name == "URL"
	}
	return false
}
