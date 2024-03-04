package types

import (
	"github.com/dave/jennifer/jen"
)

type Converter struct {
	onlySign bool
	qual     QualFunc
}

type QualFunc func(pkgPath, name string) func(s *jen.Statement)

func Convert(t any, qual QualFunc) (s *jen.Statement) {
	return newConverter(qual).Convert(t)
}

func (c *Converter) OnlySign() *Converter {
	c.onlySign = true
	return c
}

func (c *Converter) Convert(t any) (s *jen.Statement) {
	s = new(jen.Statement)
	switch t := t.(type) {
	case *Interface:
		s.Interface()
	case *Map:
		if t.IsPointer {
			s.Op("*")
		}
		s.Map(c.Convert(t.Key)).Add(c.Convert(t.Value))
	case *Array:
		if t.IsPointer {
			s.Op("*")
		}
		s.Index(jen.Lit(t.Len)).Add(c.Convert(t.Value))
	case *Slice:
		if t.IsPointer {
			s.Op("*")
		}
		s.Index().Add(c.Convert(t.Value))
	case *Var:
		s.Id(t.Name).Add(c.Convert(t.Type))
	case Vars:
		var params []jen.Code
		for _, v := range t {
			var st jen.Statement
			if !c.onlySign {
				st.Id(v.Name)
			}
			typ := v.Type
			if s, ok := typ.(*Slice); ok {
				if v.IsVariadic {
					st.Op("...")
				} else {
					st.Index()
				}
				typ = s.Value
			}
			if v.IsPointer {
				st.Op("*")
			}
			st.Add(c.Convert(typ))
			params = append(params, &st)
		}
		s.Params(params...)
	case *Sign:
		s.Add(c.Convert(t.Params))
		if len(t.Results) == 1 && t.Results[0].Name == "" {
			s.Add(c.Convert(t.Results[0]))
		} else {
			s.Add(c.Convert(t.Results))
		}
	case *Basic:
		if t.IsPointer {
			s.Op("*")
		}
		s.Id(t.Name)
	case *Named:
		if t.IsPointer {
			s.Op("*")
		}
		if t.Pkg == nil {
			s.Id(t.Name)
		} else {
			s.Do(c.qual(t.Pkg.Path, t.Name))
		}
	case *Func:
		s.Func().Id(t.Name)
		if t.Sig != nil {
			s.Add(c.Convert(t.Sig))
		}
	case *Chan:
		if t.Dir == RecvOnly {
			s.Op("<-")
		}
		s.Chan()
		if t.Dir == SendOnly {
			s.Op("<-")
		}
		s.Add(c.Convert(t.Type))
	}

	return
}

func newConverter(qual QualFunc) *Converter {
	return &Converter{qual: qual}
}
