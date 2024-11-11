package typetransform

import (
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

type DefaultTypeParse struct{}

func (s *DefaultTypeParse) Parse(valueID, assignID jen.Code, op string, t any, qualFunc types.QualFunc, errorStatements []jen.Code) (parseCode []jen.Code, paramID jen.Code, hasError bool) {
	switch t := t.(type) {
	case *types.Named:
		if basic, ok := t.Type.(*types.Basic); ok {
			if basic.IsString() {
				parseCode = append(parseCode,
					jen.Add(assignID).Op(op).Do(qualFunc(t.Pkg.Path, t.Name)).Call(valueID),
				)
				return parseCode, assignID, false
			} else {
				code, paramID, _ := For(basic).
					SetAssignID(assignID).
					SetValueID(valueID).
					SetOp(":=").
					SetQualFunc(qualFunc).
					SetErrStatements(errorStatements...).
					Parse()

				return []jen.Code{code}, jen.Qual(t.Pkg.Path, t.Name).Call(paramID), false
			}
		}
	}

	return nil, nil, false
}
