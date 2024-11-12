package typetransform

import (
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

var _ Parser = &MapTypeParse{}

type MapTypeParse struct{}

func (s *MapTypeParse) Parse(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (parseCode []jen.Code, paramID jen.Code, hasError bool) {
	slice := t.(*types.Map)
	basic := slice.Value.(*types.Basic)
	switch {
	case basic.IsSigned():
		parseCode = append(parseCode,
			jen.List(assignID, jen.Err()).Op(op).Qual("github.com/555f/go-strings", "SplitKeyValInt").Types(types.Convert(slice.Value, qualFn)).Call(valueID, jen.Lit(","), jen.Lit("="), jen.Lit(10), jen.Lit(64)),
		)
	case basic.IsUnsigned():
		parseCode = append(parseCode,
			jen.List(assignID, jen.Err()).Op(op).Qual("github.com/555f/go-strings", "SplitKeyValUint").Types(types.Convert(slice.Value, qualFn)).Call(valueID, jen.Lit(","), jen.Lit("="), jen.Lit(10), jen.Lit(64)),
		)
	case basic.IsFloat():
		parseCode = append(parseCode,
			jen.List(assignID, jen.Err()).Op(op).Qual("github.com/555f/go-strings", "SplitKeyValFloat").Types(types.Convert(slice.Value, qualFn)).Call(valueID, jen.Lit(","), jen.Lit("="), jen.Lit(64)),
		)
	case basic.IsString():
		parseCode = append(parseCode,
			jen.List(assignID, jen.Err()).Op(op).Qual("github.com/555f/go-strings", "SplitKeyValString").Types(types.Convert(slice.Value, qualFn)).Call(valueID, jen.Lit(","), jen.Lit("=")),
		)
	}
	return parseCode, assignID, true
}

func (s *MapTypeParse) Format(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (formatCode []jen.Code, paramID jen.Code, hasError bool) {
	m := t.(*types.Map)
	basic := m.Value.(*types.Basic)

	switch {
	case basic.IsInteger():
		return nil, jen.Do(qualFn("github.com/555f/go-strings", "JoinKeyValInt")).Types(types.Convert(m.Value, qualFn)).Call(valueID, jen.Lit(";"), jen.Lit("="), jen.Lit(10)), false
	case basic.IsFloat():
		return nil, jen.Do(qualFn("github.com/555f/go-strings", "JoinKeyValFloat")).Types(types.Convert(m.Value, qualFn)).Call(valueID, jen.Lit(";"), jen.Lit("="), jen.Lit('f'), jen.Lit(2), jen.Lit(64)), false
	case basic.IsString():
		return nil, jen.Do(qualFn("github.com/555f/go-strings", "JoinKeyValString")).Types(types.Convert(m.Value, qualFn)).Call(valueID, jen.Lit(";"), jen.Lit("=")), false
	}
	return
}

func (s *MapTypeParse) Support(t any) bool {
	switch t := t.(type) {
	case *types.Map:
		_, ok := t.Value.(*types.Basic)
		return ok
	}
	return false
}
