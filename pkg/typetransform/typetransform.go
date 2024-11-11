package typetransform

import (
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

func init() {
	AddTransformer(func() Transformer {
		return new(StringTypeParse)
	})
	AddTransformer(func() Transformer {
		return new(FloatTypeParse)
	})
	AddTransformer(func() Transformer {
		return new(IntTypeParse)
	})
	AddTransformer(func() Transformer {
		return new(BoolTypeParse)
	})
	AddTransformer(func() Transformer {
		return new(TimeTypeParse)
	})
	AddTransformer(func() Transformer {
		return new(URLTypeParse)
	})
	AddTransformer(func() Transformer {
		return new(NullTypeParse)
	})
	AddTransformer(func() Transformer {
		return new(SliceTypeParse)
	})
	AddTransformer(func() Transformer {
		return new(MapTypeParse)
	})
	AddTransformer(func() Transformer {
		return new(GoogleUUIDTypeParse)
	})
	AddTransformer(func() Transformer {
		return new(SatoriUUIDTypeParse)
	})
}

type Transformer interface {
	Parser
	Formatter
}

type Parser interface {
	Parse(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (parseCode []jen.Code, paramID jen.Code, hasError bool)
	Support(t any) bool
}

type Formatter interface {
	Format(valueID, assignID jen.Code, op string, t any, qualFn types.QualFunc) (formatCode []jen.Code, paramID jen.Code, hasError bool)
	Support(t any) bool
}

func AddTransformer(f func() Transformer) {
	parseFactories = append(parseFactories, parserFactory{
		factory: func() Parser {
			return f()
		},
		support: f().Support,
	})
	formatFactories = append(formatFactories, formatFactory{
		factory: func() Formatter {
			return f()
		},
		support: f().Support,
	})
}

func AddParse(f func() Parser) {
	parseFactories = append(parseFactories, parserFactory{
		factory: f,
		support: f().Support,
	})
}

func AddFormat(typeName string, f func() Formatter) {
	formatFactories = append(formatFactories, formatFactory{
		factory: f,
		support: f().Support,
	})
}

type Transform struct {
	valueID         jen.Code
	assignID        jen.Code
	op              string
	t               any
	qualFn          types.QualFunc
	errorCheck      *jen.Statement
	errorStatements []jen.Code
}

func (tr *Transform) parse(t any) (code []jen.Code, paramID jen.Code, hasError bool) {
	for _, pf := range parseFactories {
		if pf.support(t) {
			return pf.factory().Parse(tr.valueID, tr.assignID, tr.op, tr.t, tr.qualFn)
		}
	}
	return new(DefaultTypeParse).Parse(tr.valueID, tr.assignID, tr.op, tr.t, tr.qualFn, tr.errorStatements)
}

func (tr *Transform) format(t any) (code []jen.Code, paramID jen.Code, hasError bool) {
	for _, pf := range formatFactories {
		if pf.support(t) {
			return pf.factory().Format(tr.valueID, tr.assignID, tr.op, tr.t, tr.qualFn)
		}
	}
	return
}

func (tr *Transform) Parse() (jen.Code, jen.Code, bool) {
	if tr.assignID == nil {
		panic("assignID is not set")
	}
	if tr.valueID == nil {
		panic("valueID is not set")
	}
	parseCode, paramID, hasError := tr.parse(tr.t)
	code := jen.NewFile("")
	for _, c := range parseCode {
		code.Add(c)
	}
	if hasError {
		code.Add(tr.errorCheck.Block(tr.errorStatements...))
	}
	return code, paramID, hasError
}

func (tr *Transform) Format() (jen.Code, jen.Code, bool) {
	if tr.valueID == nil {
		panic("valueID is not set")
	}
	formatCode, paramID, hasError := tr.format(tr.t)
	code := jen.NewFile("")
	for _, c := range formatCode {
		code.Add(c)
	}
	if hasError {
		code.Add(tr.errorCheck.Block(tr.errorStatements...))
	}
	return code, paramID, hasError
}

func (tr *Transform) SetValueID(id jen.Code) *Transform {
	tr.valueID = id
	return tr
}

func (tr *Transform) SetAssignID(id jen.Code) *Transform {
	tr.assignID = id
	return tr
}

func (tr *Transform) SetQualFunc(qualFn types.QualFunc) *Transform {
	tr.qualFn = qualFn
	return tr
}

func (tr *Transform) SetOp(op string) *Transform {
	tr.op = op
	return tr
}

func (tr *Transform) SetErrStatements(errStatements ...jen.Code) *Transform {
	tr.errorStatements = errStatements
	return tr
}

func For(t any) *Transform {
	return &Transform{
		t:  t,
		op: ":=",
		qualFn: func(pkgPath, name string) func(s *jen.Statement) {
			return func(s *jen.Statement) {
				s.Qual(pkgPath, name)
			}
		},
		errorCheck:      jen.If(jen.Err().Op("!=").Nil()),
		errorStatements: []jen.Code{jen.Return(jen.Nil(), jen.Err())},
	}
}

var formatFactories []formatFactory
var parseFactories []parserFactory

type parserFactory struct {
	factory func() Parser
	support func(t any) bool
}

type formatFactory struct {
	factory func() Formatter
	support func(t any) bool
}
