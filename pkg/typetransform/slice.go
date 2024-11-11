package typetransform

import (
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

var _ Parser = &SliceTypeParse{}

type SliceTypeParse struct{}

func (s *SliceTypeParse) Parse(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (parseCode []jen.Code, paramID jen.Code, hasError bool) {
	slice := t.(*types.Slice)
	basic := slice.Value.(*types.Basic)
	switch {
	case basic.IsString():
		parseCode = append(parseCode,
			jen.List(assignID, jen.Err()).Op(op).Do(qualFn("github.com/555f/go-strings", "Split")).Call(valueID, jen.Lit(";")),
		)
	case basic.IsNumeric():
		parseCode = append(parseCode,
			jen.List(assignID, jen.Err()).Op(op).Do(qualFn("github.com/555f/go-strings", "SplitInt")).Types(types.Convert(slice.Value, qualFn)).Call(valueID, jen.Lit(";"), jen.Lit(10), jen.Lit(64)),
		)
	}
	return parseCode, assignID, true
}

func (s *SliceTypeParse) Format(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (formatCode []jen.Code, paramID jen.Code, hasError bool) {
	slice := t.(*types.Slice)
	basic := slice.Value.(*types.Basic)

	switch {
	case basic.IsInteger():
		return nil, jen.Do(qualFn("github.com/555f/go-strings", "JoinInt")).Types(types.Convert(slice.Value, qualFn)).Call(valueID, jen.Lit(","), jen.Lit(10)), false
	case basic.IsFloat():
		return nil, jen.Do(qualFn("github.com/555f/go-strings", "JoinFloat")).Types(types.Convert(slice.Value, qualFn)).Call(valueID, jen.Lit(","), jen.Lit('f'), jen.Lit(2), jen.Lit(64)), false
	case basic.IsString():
		return nil, jen.Do(qualFn("strings", "Join")).Call(valueID, jen.Lit(",")), false
	}

	return
}

func (s *SliceTypeParse) Support(t any) bool {
	switch t := t.(type) {
	case *types.Slice:
		_, ok := t.Value.(*types.Basic)
		return ok
	}
	return false
}
