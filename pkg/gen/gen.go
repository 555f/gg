package gen

import (
	"github.com/555f/gg/internal/plugin/http/options"
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
		switch t.Pkg.Path {
		case "time":
			if t.Name == "Time" {
				s.Add(id).Dot("Format").Call(Lit(timeFormat))
			}
		}
	}
	return s
}

func ParseValue(id, assignID Code, op string, t any, qualFunc types.QualFunc) (s *Statement) {
	s = new(Statement)
	switch t := t.(type) {
	case *types.Basic:
		if t.IsString() {
			s.Add(assignID).Op(op).Add(id)
		} else {
			s.List(assignID, Err()).Op(op).Do(parseFunc(id, t, qualFunc))
		}
	case *types.Named:
		switch t.Pkg.Path {
		case "net/url":
			switch t.Name {
			case "URL":
				s.List(assignID, Err()).Op(op).Qual("net/url", "Parse").Call(id)
			}
		case "time":
			switch t.Name {
			case "Time":
				s.List(assignID, Err()).Op(op).Qual("time", "Parse").Call(Qual("time", "RFC3339"), id)
			case "Duration":
				s.List(assignID, Err()).Op(op).Qual("time", "ParseDuration").Call(id)
			}
		case "gopkg.in/guregu/null.v4":
			switch t.Name {
			case "String":
				s.Add(assignID).Op(op).Qual(t.Pkg.Path, "StringFrom").Call(id)
			case "Int":
				s.CustomFunc(Options{Multi: true}, func(g *Group) {
					g.Var().Id("val").Int64()
					g.List(Id("val"), Err()).Op("=").Qual("strconv", "ParseInt").Call(id, Lit(10), Lit(64))
					g.Do(CheckErr(Return()))
					g.Add(assignID).Op(op).Qual(t.Pkg.Path, "IntFrom").Call(Id("val"))
				})
			case "Float":
				s.CustomFunc(Options{Multi: true}, func(g *Group) {
					g.Var().Id("val").Float64()
					g.List(Id("val"), Err()).Op("=").Qual("strconv", "ParseFloat").Call(id, Lit(10), Lit(64))
					g.Do(CheckErr(Return()))
					g.Add(assignID).Op(op).Qual(t.Pkg.Path, "FloatFrom").Call(Id("val"))
				})
			case "Bool":
				s.CustomFunc(Options{Multi: true}, func(g *Group) {
					g.Var().Id("val").Bool()
					g.List(Id("val"), Err()).Op("=").Qual("strconv", "ParseBool").Call(id)
					g.Do(CheckErr(Return()))
					g.Add(assignID).Op(op).Qual(t.Pkg.Path, "BoolFrom").Call(Id("val"))
				})
			case "Time":
				s.CustomFunc(Options{Multi: true}, func(g *Group) {
					g.Var().Id("val").Qual("time", "Time")
					g.List(Id("val"), Err()).Op("=").Qual("time", "Parse").Call(Qual("time", "RFC3339"), id)
					g.Do(CheckErr(Return()))
					g.Add(assignID).Op(op).Qual(t.Pkg.Path, "TimeFrom").Call(Id("val"))
				})
			}
		}
	case *types.Slice:
		switch tv := t.Value.(type) {
		case *types.Basic:
			if tv.IsString() {
				s.List(assignID, Err()).Op(op).Qual("github.com/555f/go-strings", "Split").Call(id, Lit(";"))
			}
			if tv.IsNumeric() {
				s.List(assignID, Err()).Op(op).Qual("github.com/555f/go-strings", "SplitInt").Types(types.Convert(t.Value, qualFunc)).Call(id, Lit(";"), Lit(10), Lit(64))
			}
		}
	case *types.Map:
		switch tv := t.Value.(type) {
		case *types.Basic:
			if tv.IsSigned() {
				s.List(assignID, Err()).Op(op).Qual("github.com/555f/go-strings", "SplitKeyValInt").Types(types.Convert(t.Value, qualFunc)).Call(id, Lit(","), Lit("="), Lit(10), Lit(64))
			}
			if tv.IsUnsigned() {
				s.List(assignID, Err()).Op(op).Qual("github.com/555f/go-strings", "SplitKeyValUint").Types(types.Convert(t.Value, qualFunc)).Call(id, Lit(","), Lit("="), Lit(10), Lit(64))
			}
			if tv.IsFloat() {
				s.List(assignID, Err()).Op(op).Qual("github.com/555f/go-strings", "SplitKeyValFloat").Types(types.Convert(t.Value, qualFunc)).Call(id, Lit(","), Lit("="), Lit(64))
			}
			if tv.IsString() {
				s.List(assignID, Err()).Op(op).Qual("github.com/555f/go-strings", "SplitKeyValString").Types(types.Convert(t.Value, qualFunc)).Call(id, Lit(","), Lit("="))
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
		case "net/url", "time", "gopkg.in/guregu/null.v4":
			return nil
		}
		return ExtractFields(t.Type)
	}
}

func WrapResponse(names []string, results []*options.EndpointResult, qualFunc types.QualFunc) func(g *Group) {
	return func(g *Group) {
		if len(names) > 0 {
			g.Id(strcase.ToCamel(names[0])).StructFunc(WrapResponse(names[1:], results, qualFunc)).Tag(map[string]string{"json": names[0]})
		} else {
			for _, result := range results {
				g.Id(result.FldNameExport).Add(types.Convert(result.Type, qualFunc)).Tag(map[string]string{"json": result.Name})
			}
		}
	}
}
