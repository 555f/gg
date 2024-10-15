package types

import (
	"github.com/dave/jennifer/jen"
)

type Construct struct {
	qual QualFunc

	c *Converter
}

func NewConstruct(qual QualFunc) *Construct {
	return &Construct{qual: qual, c: NewConvert(qual)}
}

func (c *Construct) Convert(t any) (s *jen.Statement) {
	s = new(jen.Statement)
	switch t := t.(type) {
	// case *Interface:
	// s.Interface()
	case *Map:
		// if t.IsPointer {
		// 	s.Op("*")
		// }
		s.Map(c.c.Convert(t.Key)).Add(c.c.Convert(t.Value)).Values()
	case *Array:
		// if t.IsPointer {
		// s.Op("*")
		// }
		s.Index(jen.Lit(t.Len)).Add(c.c.Convert(t.Value)).Values()
	case *Slice:
		// if t.IsPointer {
		// s.Op("*")
		// }
		s.Index().Add(c.c.Convert(t.Value)).Values()
	case *Var:
		s.Id(t.Name).Op(":=").Add(c.Convert(t.Type))
	// case Vars:
	// 	var params []jen.Code
	// 	for _, v := range t {
	// 		var st jen.Statement
	// 		st.Id(v.Name)
	// 		typ := v.Type
	// 		if s, ok := typ.(*Slice); ok {
	// 			if v.IsVariadic {
	// 				st.Op("...")
	// 			} else {
	// 				st.Index()
	// 			}
	// 			typ = s.Value
	// 		}
	// 		st.Add(c.Convert(typ))
	// 		params = append(params, &st)
	// 	}
	// 	s.Params(params...)
	// case *Sign:
	// 	s.Add(c.Convert(t.Params))
	// 	if len(t.Results) == 1 && t.Results[0].Name == "" {
	// 		s.Add(c.Convert(t.Results[0]))
	// 	} else {
	// 		s.Add(c.Convert(t.Results))
	// 	}
	case *Basic:
		// if t.IsPointer {
		// s.Op("*")
		// }
		s.Id(t.Name)
	case *Named:
		if t.IsPointer {
			s.Op("&")
		}
		if t.Pkg == nil {
			s.Id(t.Name)
		} else {
			s.Do(c.qual(t.Pkg.Path, t.Name))
		}
		s.Values()
		// case *Func:
		// 	s.Func().Id(t.Name)
		// 	if t.Sig != nil {
		// 		s.Add(c.Convert(t.Sig))
		// 	}
		// case *Chan:
		// 	if t.Dir == RecvOnly {
		// 		s.Op("<-")
		// 	}
		// 	s.Chan()
		// 	if t.Dir == SendOnly {
		// 		s.Op("->")
		// 	}
		// 	s.Add(c.Convert(t.Type))
	}

	return
}
