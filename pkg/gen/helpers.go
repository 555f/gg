package gen

import (
	"github.com/555f/gg/pkg/types"
	. "github.com/dave/jennifer/jen"
)

func formatFunc(id Code, t *types.Basic) func(s *Statement) {
	return func(s *Statement) {
		if t.IsSigned() {
			s.Qual("strconv", "FormatInt").CallFunc(func(g *Group) {
				if t.BitSize() == 64 && !t.IsInt() {
					g.Add(id)
				} else {
					g.Id("int64").Call(id)
				}
				g.Lit(10)
			})
		} else if t.IsUnsigned() {
			s.Qual("strconv", "FormatUint").CallFunc(func(g *Group) {
				if t.BitSize() == 64 && t.IsUint() {
					g.Add(id)
				} else {
					g.Id("uint64").Call(id)
				}
				g.Lit(10)
			})
		} else if t.IsFloat() {
			s.Qual("strconv", "FormatFloat").CallFunc(func(g *Group) {
				if t.BitSize() == 64 {
					g.Add(id)
				} else {
					g.Id("float64").Call(id)
				}
				g.LitRune('g')
				g.Lit(-1)
				g.Lit(t.BitSize())
			})
		} else if t.IsBool() {
			s.Qual("strconv", "FormatBool").Call(id)
		} else if t.IsString() {
			s.Add(id)
		}
	}
}

func parseFunc(id Code, t *types.Basic, qualFunc types.QualFunc) func(s *Statement) {
	return func(s *Statement) {
		if t.IsSigned() {
			s.Qual("github.com/555f/go-strings", "ParseInt").Types(types.Convert(t, qualFunc)).Call(id, Lit(10), Lit(t.BitSize()))
		}
		if t.IsUnsigned() {
			s.Qual("github.com/555f/go-strings", "ParseUint").Types(types.Convert(t, qualFunc)).Call(id, Lit(10), Lit(t.BitSize()))
		}
		if t.IsFloat() {
			s.Qual("github.com/555f/go-strings", "ParseFloat").Types(types.Convert(t, qualFunc)).Call(id, Lit(t.BitSize()))
		}
		if t.IsBool() {
			s.Qual("github.com/555f/go-strings", "ParseBool").Call(id)
		}
	}
}
