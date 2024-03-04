package pgen

func Service(name string, methods ...Code) *Statement {
	return newStatement().Service(name, methods...)
}

func (g *Group) Service(name string, methods ...Code) *Statement {
	s := Service(name, methods...)
	g.items = append(g.items, s)
	return s
}

func (s *Statement) Service(name string, methods ...Code) *Statement {
	g := &Group{
		close:     "}",
		items:     methods,
		multi:     true,
		name:      "service",
		open:      "service " + name + " {",
		separator: "{}",
	}
	*s = append(*s, g)
	return s
}

func RPC() *Statement {
	return newStatement().RPC()
}

func (g *Group) RPC() *Statement {
	// notest
	s := newStatement().RPC()
	g.items = append(g.items, s)
	return s
}

func (s *Statement) RPC() *Statement {
	// notest
	t := token{
		content: "rpc",
		typ:     keywordToken,
	}
	*s = append(*s, t)
	return s
}

func Message() *Statement {
	return newStatement().Message()
}

func (g *Group) Message() *Statement {
	// notest
	s := newStatement().Message()
	g.items = append(g.items, s)
	return s
}

func (s *Statement) Message() *Statement {
	// notest
	t := token{
		content: "message",
		typ:     keywordToken,
	}
	*s = append(*s, t)
	return s
}

func (s *Statement) Request(request Code) *Statement {
	g := &Group{
		close:     ")",
		items:     []Code{request},
		name:      "params",
		open:      "(",
		separator: "",
	}
	*s = append(*s, g)
	return s
}

func (s *Statement) Returns(returns Code) *Statement {
	g := &Group{
		close:     ")",
		items:     []Code{returns},
		name:      "returns",
		open:      "returns (",
		separator: "",
	}
	*s = append(*s, g)
	return s
}

func (s *Statement) Values(values ...Code) *Statement {
	g := &Group{
		close:     "}",
		items:     values,
		name:      "values",
		open:      "{",
		separator: ";",
	}
	*s = append(*s, g)
	return s
}

func Lit(v interface{}) *Statement {
	return newStatement().Lit(v)
}

func (g *Group) Lit(v interface{}) *Statement {
	s := Lit(v)
	g.items = append(g.items, s)
	return s
}

func (s *Statement) Lit(v interface{}) *Statement {
	t := token{
		typ:     literalToken,
		content: v,
	}
	*s = append(*s, t)
	return s
}
