package pgen

import (
	"fmt"
	"io"
)

type tokenType string

const (
	packageToken    tokenType = "package"
	identifierToken tokenType = "identifier"
	keywordToken    tokenType = "keyword"
	operatorToken   tokenType = "operator"
	delimiterToken  tokenType = "delimiter"
	layoutToken     tokenType = "layout"
	literalToken    tokenType = "literal"
)

type token struct {
	typ     tokenType
	content interface{}
}

func (t token) render(f *File, w io.Writer, s *Statement) error {
	switch t.typ {
	case literalToken:
		var out string
		switch t.content.(type) {
		case bool, string, int, complex128:
			out = fmt.Sprintf("%#v", t.content)
		default:
			panic(fmt.Sprintf("unsupported type for literal: %T", t.content))
		}
		if _, err := w.Write([]byte(out)); err != nil {
			return err
		}
	case keywordToken, operatorToken, layoutToken, delimiterToken:
		if _, err := w.Write([]byte(fmt.Sprintf("%s", t.content))); err != nil {
			return err
		}
	case identifierToken:
		if _, err := w.Write([]byte(t.content.(string))); err != nil {
			return err
		}
	}
	return nil
}

func Empty() *Statement {
	return newStatement().Empty()
}

func (g *Group) Empty() *Statement {
	s := Empty()
	g.items = append(g.items, s)
	return s
}

func (s *Statement) Empty() *Statement {
	t := token{
		typ:     operatorToken,
		content: "",
	}
	*s = append(*s, t)
	return s
}

func Op(op string) *Statement {
	return newStatement().Op(op)
}

func (g *Group) Op(op string) *Statement {
	s := Op(op)
	g.items = append(g.items, s)
	return s
}

func (s *Statement) Op(op string) *Statement {
	t := token{
		typ:     operatorToken,
		content: op,
	}
	*s = append(*s, t)
	return s
}

func Dot(name string) *Statement {
	// notest
	return newStatement().Dot(name)
}

func (g *Group) Dot(name string) *Statement {
	// notest
	s := Dot(name)
	g.items = append(g.items, s)
	return s
}

func (s *Statement) Dot(name string) *Statement {
	d := token{
		typ:     delimiterToken,
		content: ".",
	}
	t := token{
		typ:     identifierToken,
		content: name,
	}
	*s = append(*s, d, t)
	return s
}

func Id(name string) *Statement {
	return newStatement().Id(name)
}

func (g *Group) Id(name string) *Statement {
	s := Id(name)
	g.items = append(g.items, s)
	return s
}

func (s *Statement) Id(name string) *Statement {
	t := token{
		typ:     identifierToken,
		content: name,
	}
	*s = append(*s, t)
	return s
}

func Option() *Statement {
	return newStatement().Option()
}

func (g *Group) Option() *Statement {
	s := Option()
	g.items = append(g.items, s)
	return s
}

func (s *Statement) Option() *Statement {
	return s.Id("option")
}

func Stream() *Statement {
	return newStatement().Stream()
}

func (g *Group) Stream() *Statement {
	s := Option()
	g.items = append(g.items, s)
	return s
}

func (s *Statement) Stream() *Statement {
	return s.Id("stream")
}

func Qual(path string) *Statement {
	return newStatement().Qual(path)
}

func (g *Group) Qual(path string) *Statement {
	s := Qual(path)
	g.items = append(g.items, s)
	return s
}

func (s *Statement) Qual(path string) *Statement {
	g := &Group{
		close: "",
		items: []Code{
			token{
				typ:     packageToken,
				content: path,
			},
			token{
				typ:     identifierToken,
				content: path,
			},
		},
		name:      "qual",
		open:      "",
		separator: "",
	}
	*s = append(*s, g)
	return s
}

func Line() *Statement {
	return newStatement().Line()
}

func (g *Group) Line() *Statement {
	s := Line()
	g.items = append(g.items, s)
	return s
}

func (s *Statement) Line() *Statement {
	t := token{
		typ:     layoutToken,
		content: "\n",
	}
	*s = append(*s, t)
	return s
}
