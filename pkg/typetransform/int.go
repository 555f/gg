package typetransform

import (
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

var _ Parser = &IntTypeParse{}

type IntTypeParse struct{}

func (s *IntTypeParse) Parse(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (parseCode []jen.Code, paramID jen.Code, hasError bool) {
	b := t.(*types.Basic)
	var parseFunc string
	switch {
	case b.IsSigned():
		parseFunc = "ParseInt"

	case b.IsUnsigned():
		parseFunc = "ParseUint"
	}
	parseCode = []jen.Code{jen.List(assignID, jen.Err()).Op(op).Do(qualFn("github.com/555f/go-strings", parseFunc)).Types(types.Convert(t, qualFn)).Call(valueID, jen.Lit(10), jen.Lit(b.BitSize()))}
	return parseCode, assignID, true
}

func (s *IntTypeParse) Format(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (formatCode []jen.Code, paramID jen.Code, hasError bool) {
	b := t.(*types.Basic)
	switch {
	case b.IsSigned():
		return nil, jen.Qual("strconv", "FormatInt").CallFunc(func(g *jen.Group) {
			if b.BitSize() == 64 && !b.IsInt() {
				g.Add(valueID)
			} else {
				g.Id("int64").Call(valueID)
			}
			g.Lit(10)
		}), false
	case b.IsUnsigned():
		return nil, jen.Qual("strconv", "FormatUint").CallFunc(func(g *jen.Group) {
			if b.BitSize() == 64 && b.IsUint() {
				g.Add(valueID)
			} else {
				g.Id("uint64").Call(valueID)
			}
			g.Lit(10)
		}), false
	}

	return nil, nil, false
}

func (s *IntTypeParse) Support(t any) bool {
	switch t := t.(type) {
	case *types.Basic:
		switch t.Name() {
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
			return true
		}
	}
	return false
}
