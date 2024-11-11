package typetransform

import (
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"

	stdtypes "go/types"
)

var _ Parser = &NullTypeParse{}

var nameToBasicType = map[string]*types.Basic{
	"Int":   types.BasicTyp[stdtypes.Int],
	"Float": types.BasicTyp[stdtypes.Float64],
	"Bool":  types.BasicTyp[stdtypes.Bool],
}

type NullTypeParse struct{}

func (s *NullTypeParse) Parse(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (parseCode []jen.Code, paramID jen.Code, hasError bool) {
	n := t.(*types.Named)

	if n.Name == "Time" {
		code, paramID, _ := For(types.NamedTyp["time.Time"]).
			SetAssignID(assignID).
			SetValueID(valueID).
			SetOp(op).
			Parse()

		if code != nil {
			parseCode = append(parseCode,
				code,
			)
		}

		paramID = jen.Do(qualFn(n.Pkg.Path, "TimeFrom")).Call(paramID)

		return parseCode, paramID, false
	}

	if basic, ok := nameToBasicType[n.Name]; ok {
		code, paramID, _ := For(basic).
			SetAssignID(assignID).
			SetValueID(valueID).
			SetOp(op).
			Parse()

		if code != nil {
			parseCode = append(parseCode,
				code,
			)
		}

		paramID = jen.Do(qualFn(n.Pkg.Path, n.Name+"From")).Call(paramID)

		return parseCode, paramID, false
	}

	return nil, nil, false
}

func (s *NullTypeParse) Format(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (formatCode []jen.Code, paramID jen.Code, hasError bool) {
	n := t.(*types.Named)
	switch n.Name {
	case "String", "Time", "Int", "Float", "Bool":
		return nil, jen.Add(valueID).Dot("ValueOrZero").Call(), false
	}
	return
}

func (s *NullTypeParse) Support(t any) bool {
	switch t := t.(type) {
	case *types.Named:
		return t.Pkg.Path == "gopkg.in/guregu/null.v4"
	}
	return false
}
