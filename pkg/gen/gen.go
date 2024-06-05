package gen

import (
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"

	. "github.com/dave/jennifer/jen"
)

func CheckErr(statements ...Code) func(s *Statement) {
	return func(s *Statement) {
		s.If(Err().Op("!=").Nil()).Block(statements...)
	}
}

func FormatValue(id Code, t any, qualFunc types.QualFunc, timeFormat string) (s *Statement) {
	s = new(Statement)
	switch t := t.(type) {
	case *types.Basic:
		s.Do(formatFunc(id, t))
	case *types.Slice:
		if basic, ok := t.Value.(*types.Basic); ok {
			if basic.IsNumeric() {
				s.Qual("github.com/555f/go-strings", "JoinInt").Types(types.Convert(t.Value, qualFunc)).Call(id, Lit(","), Lit(10))
			}
			if basic.IsFloat() {
				s.Qual("github.com/555f/go-strings", "JoinFloat").Types(types.Convert(t.Value, qualFunc)).Call(id, Lit(","), Lit('f'), Lit(2), Lit(64))
			}
			if basic.IsString() {
				s.Qual("strings", "Join").Call(id, Lit(","))
			}
		}
	case *types.Named:
		var stingMethodFound bool
		for _, m := range t.Methods {
			if m.Name == "String" && len(m.Sig.Results) == 1 {
				if b, ok := m.Sig.Results[0].Type.(*types.Basic); ok && b.IsString() {
					stingMethodFound = true
				}
			}
		}
		if stingMethodFound {
			s.Add(id).Dot("String").Call()
			return s
		}
		switch t.Pkg.Path {
		case "gopkg.in/guregu/null.v4":
			switch t.Name {
			case "String":
				s.Add(id).Dot("String")
			case "Time":
				s.Add(id).Dot("Time").Dot("Format").Call(Lit(timeFormat))
			case "Int":
				s.Qual("strconv", "FormatInt").Call(Add(id).Dot("Int64"), Lit(10))
			case "Float":
				s.Qual("strconv", "FormatFloat").Call(Add(id).Dot("Float64"), LitRune('g'), Lit(-1), Lit(64))
			case "Bool":
				s.Qual("strconv", "FormatBool").Call(Add(id).Dot("Bool"))
			}
		case "time":
			if t.Name == "Time" {
				s.Add(id).Dot("Format").Call(Lit(timeFormat))
			}
		}
	}
	return s
}

func ParseValue(id, assignID Code, op string, t any, qualFunc types.QualFunc, errCalback func() Code) (s *Statement) {
	s = Custom(Options{Multi: true})
	switch t := t.(type) {
	case *types.Basic:
		if t.IsString() {
			s.Add(assignID).Op(op).Add(id)
		} else {
			s.List(assignID, Err()).Op(op).Do(parseFunc(id, t, qualFunc))
			s.Do(func(s *Statement) {
				s.Custom(Options{Multi: true}, errCalback())
			})
		}
	case *types.Named:
		if basic, ok := t.Type.(*types.Basic); ok {
			if basic.IsString() {
				s.Add(assignID).Op(op).Do(qualFunc(t.Pkg.Path, t.Name)).Call(id)
			} else {
				s.CustomFunc(Options{Multi: true}, func(g *Group) {
					g.List(Id("v"), Err()).Op(":=").Do(parseFunc(id, basic, qualFunc))
					g.Do(CheckErr(
						Return(Nil(), Err()),
					))
					g.Add(assignID).Op(op).Qual(t.Pkg.Path, t.Name).Call(Id("v"))
				})
			}
			return
		}

		switch t.Pkg.Path {
		case "net/url":
			switch t.Name {
			case "URL":
				s.List(assignID, Err()).Op(op).Qual("net/url", "Parse").Call(id)
				s.Do(func(s *Statement) {
					s.Custom(Options{Multi: true}, errCalback())
				})
			}
		case "time":
			switch t.Name {
			case "Time":
				s.List(assignID, Err()).Op(op).Qual("time", "Parse").Call(Qual("time", "RFC3339"), id)
				s.Do(func(s *Statement) {
					s.Custom(Options{Multi: true}, errCalback())
				})
			case "Duration":
				s.List(assignID, Err()).Op(op).Qual("time", "ParseDuration").Call(id)
				s.Do(func(s *Statement) {
					s.Custom(Options{Multi: true}, errCalback())
				})
			}
		case "github.com/google/uuid":
			switch t.Name {
			case "UUID":
				s.List(assignID, Err()).Op(op).Qual(t.Pkg.Path, "Parse").Call(id)
				s.Do(func(s *Statement) {
					s.Custom(Options{Multi: true}, errCalback())
				})
			}
		case "github.com/satori/go.uuid":
			switch t.Name {
			case "UUID":
				s.List(assignID, Err()).Op(op).Qual(t.Pkg.Path, "FromString").Call(id)
				s.Do(func(s *Statement) {
					s.Custom(Options{Multi: true}, errCalback())
				})
			}
		case "gopkg.in/guregu/null.v4":
			switch t.Name {
			case "String":
				s.Add(assignID).Op(op).Qual(t.Pkg.Path, "StringFrom").Call(id)
			case "Int":
				s.CustomFunc(Options{Multi: true}, func(g *Group) {
					g.Var().Id("val").Int64()
					g.List(Id("val"), Err()).Op("=").Qual("strconv", "ParseInt").Call(id, Lit(10), Lit(64))
					s.Do(func(s *Statement) {
						s.Custom(Options{Multi: true}, errCalback())
					})
					// g.Do(CheckErr(Return()))
					g.Add(assignID).Op(op).Qual(t.Pkg.Path, "IntFrom").Call(Id("val"))
				})
			case "Float":
				s.CustomFunc(Options{Multi: true}, func(g *Group) {
					g.Var().Id("val").Float64()
					g.List(Id("val"), Err()).Op("=").Qual("strconv", "ParseFloat").Call(id, Lit(10), Lit(64))
					s.Do(func(s *Statement) {
						s.Custom(Options{Multi: true}, errCalback())
					})
					// g.Do(CheckErr(Return()))
					g.Add(assignID).Op(op).Qual(t.Pkg.Path, "FloatFrom").Call(Id("val"))
				})
			case "Bool":
				s.CustomFunc(Options{Multi: true}, func(g *Group) {
					g.Var().Id("val").Bool()
					g.List(Id("val"), Err()).Op("=").Qual("strconv", "ParseBool").Call(id)
					s.Do(func(s *Statement) {
						s.Custom(Options{Multi: true}, errCalback())
					})
					// g.Do(CheckErr(Return()))
					g.Add(assignID).Op(op).Qual(t.Pkg.Path, "BoolFrom").Call(Id("val"))
				})
			case "Time":
				s.CustomFunc(Options{Multi: true}, func(g *Group) {
					g.Var().Id("val").Qual("time", "Time")
					g.List(Id("val"), Err()).Op("=").Qual("time", "Parse").Call(Qual("time", "RFC3339"), id)
					s.Do(func(s *Statement) {
						s.Custom(Options{Multi: true}, errCalback())
					})
					// g.Do(CheckErr(Return()))
					g.Add(assignID).Op(op).Qual(t.Pkg.Path, "TimeFrom").Call(Id("val"))
				})
			}
		}
	case *types.Slice:
		switch tv := t.Value.(type) {
		case *types.Basic:
			if tv.IsString() {
				s.List(assignID, Err()).Op(op).Qual("github.com/555f/go-strings", "Split").Call(id, Lit(";"))
				s.Do(func(s *Statement) {
					s.Custom(Options{Multi: true}, errCalback())
				})
			}
			if tv.IsNumeric() {
				s.List(assignID, Err()).Op(op).Qual("github.com/555f/go-strings", "SplitInt").Types(types.Convert(t.Value, qualFunc)).Call(id, Lit(";"), Lit(10), Lit(64))
				s.Do(func(s *Statement) {
					s.Custom(Options{Multi: true}, errCalback())
				})
			}
		}
	case *types.Map:
		switch tv := t.Value.(type) {
		case *types.Basic:
			if tv.IsSigned() {
				s.List(assignID, Err()).Op(op).Qual("github.com/555f/go-strings", "SplitKeyValInt").Types(types.Convert(t.Value, qualFunc)).Call(id, Lit(","), Lit("="), Lit(10), Lit(64))
				s.Do(func(s *Statement) {
					s.Custom(Options{Multi: true}, errCalback())
				})
			}
			if tv.IsUnsigned() {
				s.List(assignID, Err()).Op(op).Qual("github.com/555f/go-strings", "SplitKeyValUint").Types(types.Convert(t.Value, qualFunc)).Call(id, Lit(","), Lit("="), Lit(10), Lit(64))
				s.Do(func(s *Statement) {
					s.Custom(Options{Multi: true}, errCalback())
				})
			}
			if tv.IsFloat() {
				s.List(assignID, Err()).Op(op).Qual("github.com/555f/go-strings", "SplitKeyValFloat").Types(types.Convert(t.Value, qualFunc)).Call(id, Lit(","), Lit("="), Lit(64))
				s.Do(func(s *Statement) {
					s.Custom(Options{Multi: true}, errCalback())
				})
			}
			if tv.IsString() {
				s.List(assignID, Err()).Op(op).Qual("github.com/555f/go-strings", "SplitKeyValString").Types(types.Convert(t.Value, qualFunc)).Call(id, Lit(","), Lit("="))
				s.Do(func(s *Statement) {
					s.Custom(Options{Multi: true}, errCalback())
				})
			}
		}
	}
	return s
}

func ExtractFields(v any) []*types.StructFieldType {
	switch t := v.(type) {
	default:
		return nil
	case *types.Struct:
		return t.Fields
	case *types.Named:
		switch t.Pkg.Path {
		case "net/url", "time", "gopkg.in/guregu/null.v4", "github.com/google/uuid":
			return nil
		}
		return ExtractFields(t.Type)
	}
}

func WrapResponse(names []string, qualFunc types.QualFunc) func(g *Group) {
	return func(g *Group) {
		if len(names) > 0 {
			g.Id(strcase.ToCamel(names[0])).StructFunc(WrapResponse(names[1:], qualFunc)).Tag(map[string]string{"json": names[0]})
		}
	}
}
