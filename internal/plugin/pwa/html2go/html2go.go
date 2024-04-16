package html2go

import (
	"bytes"
	"fmt"

	"strconv"
	"strings"

	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/antlr4-go/antlr/v4"
	"github.com/dave/jennifer/jen"

	"github.com/555f/gg/internal/plugin/pwa/html2go/parser"
	"github.com/555f/gg/internal/plugin/pwa/varparser"
)

const pkgApp = "github.com/maxence-charriere/go-app/v9/pkg/app"

type tagType int

const (
	unsetTagType tagType = iota
	rangeTagType
	ifTagType
)

type htmlListener struct {
	*parser.BaseHTMLParserListener

	name      string
	structMap map[string]*gg.Struct
	buf       bytes.Buffer

	classBuf bytes.Buffer

	index     int
	prevIndex int

	tagType      tagType
	tagTypeValue string

	prevVars []varparser.Var
}

func (s *htmlListener) EnterHtmlMisc(ctx *parser.HtmlMiscContext) {}

func (s *htmlListener) EnterHtmlChardata(ctx *parser.HtmlChardataContext) {
	// n := ctx.HTML_TEXT()
	// if n == nil {
	// 	return
	// }

	// s.index++

	// if s.index == s.prevIndex {
	// 	fmt.Fprintf(&s.buf, ",")
	// }
	// text := strings.TrimSpace(n.GetText())
	// vars := varparser.Parse(text)

	// for _, v := range vars {
	// 	st := s.structMap[s.name]
	// 	if f := st.Type.Path(v.ID); f != nil {
	// 		if b, ok := f.Var.Type.(*types.Basic); ok {
	// 			var fmtf string
	// 			switch {
	// 			default:
	// 				fmtf = "%%s"
	// 			case b.IsNumeric():
	// 				fmtf = "%%d"
	// 			case b.IsFloat():
	// 				fmtf = "%%f"
	// 			}
	// 			text = text[0:v.Pos.Start] + fmtf + text[v.Pos.Finish:]
	// 			s.prevVars = append(s.prevVars, v)
	// 		}
	// 	}
	// }
	// fmt.Fprintf(&s.buf, "app.Text")
	// if len(s.prevVars) > 0 {
	// 	fmt.Fprintf(&s.buf, "f")
	// }

	// fmt.Fprintf(&s.buf, "("+strconv.Quote(text))
}

func (s *htmlListener) ExitHtmlChardata(ctx *parser.HtmlChardataContext) {
	// n := ctx.HTML_TEXT()
	// if n == nil {
	// 	return
	// }
	// if len(s.prevVars) > 0 {
	// 	for _, v := range s.prevVars {
	// 		fmt.Fprintf(&s.buf, ",")
	// 		fmt.Fprintf(&s.buf, "c."+v.ID)
	// 	}
	// 	s.prevVars = []varparser.Var{}
	// }

	// fmt.Fprintf(&s.buf, ")")

	// s.prevIndex = s.index

	// s.index--
}

func (s *htmlListener) EnterHtmlElement(ctx *parser.HtmlElementContext) {
	// s.index++
	// tagName := strcase.ToCamel(ctx.TAG_NAME(0).GetText())
	// if s.index == s.prevIndex {
	// 	fmt.Fprintf(&s.buf, ",")
	// }

	// fmt.Println(tagName)

	// switch s.tagType {
	// default:
	// 	fmt.Fprintf(&s.buf, "app.%s().Body(", tagName)
	// 	// case ifTagType:
	// 	// 	fmt.Fprintf(&s.buf, "app.%s().Body(app.If(c.%s, ", tagName, s.tagTypeValue)
	// 	// case rangeTagType:
	// 	// 	idx := strings.Index(s.tagTypeValue, "in")
	// 	// 	if idx > -1 {
	// 	// 		vRangeKey := strings.TrimSpace(s.tagTypeValue[:idx])
	// 	// 		vRangeSlice := strings.TrimSpace(s.tagTypeValue[idx+2:])
	// 	// 		fmt.Fprintf(&s.buf, "app.%s().Body(app.Range(c.%s).Slice(func (%s int) app.UI { return ", tagName, vRangeSlice, vRangeKey)
	// 	// 	}
	// }
}

func (s *htmlListener) ExitHtmlElement(ctx *parser.HtmlElementContext) {
	// switch s.tagType {
	// default:
	// 	fmt.Fprintf(&s.buf, ")")
	// 	// case ifTagType:
	// 	// 	fmt.Fprint(&s.buf, "))")
	// 	// case rangeTagType:
	// 	// 	fmt.Fprint(&s.buf, "}))")
	// }

	// if s.classBuf.Len() > 0 {
	// 	fmt.Fprint(&s.buf, s.classBuf.String())
	// 	s.classBuf.Reset()
	// }

	// s.tagType = unsetTagType
	// s.tagTypeValue = ""

	// s.prevIndex = s.index

	// s.index--
}

func (s *htmlListener) EnterHtmlAttribute(ctx *parser.HtmlAttributeContext) {
	attrName := ctx.TAG_NAME().GetText()
	switch attrName {
	default:
		fmt.Fprintf(&s.classBuf, ".%s(", strcase.ToCamel(attrName))
	case "v-if":
		s.tagType = ifTagType
		s.tagTypeValue, _ = strconv.Unquote(ctx.ATTVALUE_VALUE().GetText())
	case "v-range":
		s.tagType = rangeTagType
		s.tagTypeValue, _ = strconv.Unquote(ctx.ATTVALUE_VALUE().GetText())
	}
}

func (s *htmlListener) ExitHtmlAttribute(ctx *parser.HtmlAttributeContext) {
	if s.tagType == unsetTagType {
		fmt.Fprintf(&s.classBuf, ")")
	}
}

var _ antlr.ErrorListener = &errorListener{}

type errorListener struct{}

// ReportAmbiguity implements antlr.ErrorListener.
func (e *errorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex int, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
	// fmt.Println(startIndex, stopIndex)
}

// ReportAttemptingFullContext implements antlr.ErrorListener.
func (e *errorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex int, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
	// fmt.Println(startIndex, stopIndex)
}

// ReportContextSensitivity implements antlr.ErrorListener.
func (e *errorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex int, stopIndex int, prediction int, configs *antlr.ATNConfigSet) {
	fmt.Println("ReportContextSensitivity", startIndex, stopIndex, prediction)
}

// SyntaxError implements antlr.ErrorListener.
func (*errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line int, column int, msg string, e antlr.RecognitionException) {
	fmt.Println("SyntaxError", line, column, msg)
}

type HTML2Go struct {
	s         *gg.Struct
	structMap map[string]*gg.Struct
	qual      func(path, name string) func(s *jen.Statement)
}

func (hg *HTML2Go) Parse(html string) (codes []jen.Code, err error) {
	is := antlr.NewInputStream(html)
	lexer := parser.NewHTMLLexer(is)

	// lexer.AddErrorListener(&errorListener{})

	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewHTMLParser(stream)

	// p.AddErrorListener(&errorListener{})
	codes = hg.recursiveParse(p.HtmlDocument(), 0)

	return
}

func (hg *HTML2Go) recursiveParse(t antlr.Tree, nested int) (codes []jen.Code) {
	if r, ok := t.(antlr.RuleNode); ok {
		ctx := r.GetRuleContext().(antlr.ParserRuleContext)
		switch e := ctx.(type) {
		default:
			for _, n := range ctx.GetChildren() {
				codes = append(codes, hg.recursiveParse(n, nested+1)...)
			}
		case *parser.HtmlElementContext:
			tagName := e.TAG_NAME(0).GetText()

			// fmt.Println(strings.Repeat(" ", nested), tagName)

			if tagName == "svg" {
				codes = append(codes, jen.Qual(pkgApp, "Raw").Call(jen.Lit(r.GetText())))
				return
			}

			tag := jen.Qual(pkgApp, strcase.ToCamel(tagName)).Call()

			var (
				callCodes   []jen.Code
				vRange, vIf string
				cmpMethods  []*types.Func
				isFoundComp bool
			)

			if ss, ok := hg.structMap[tagName]; ok {
				isFoundComp = true
				cmpMethods = ss.Named.Methods

				tag = jen.Op("&").Do(hg.qual(ss.Named.Pkg.Path, ss.Named.Name))
				tag.Values()
				tag = jen.Call(tag)
			} else if p := hg.s.Type.Path(strcase.ToLowerCamel(tagName)); p != nil {
				isFoundComp = true
				switch tt := p.Var.Type.(type) {
				case *types.Named:
					if iface := tt.Interface(); iface != nil {
						cmpMethods = iface.Methods
					} else {
						cmpMethods = tt.Methods
					}
				}
				tag = jen.Id("c").Dot(strcase.ToLowerCamel(tagName))
			}

			cmpMethodMap := make(map[string]*types.Func, len(cmpMethods))
			for _, m := range cmpMethods {
				cmpMethodMap[m.Name] = m
			}

			for _, a := range e.AllHtmlAttribute() {
				attrName := a.TAG_NAME().GetText()
				attrVal := trimQuote(a.ATTVALUE_VALUE().GetText())

				if f, ok := cmpMethodMap[strcase.ToCamel(attrName)]; ok && f.Sig.Params.Len() > 0 {
					pt := f.Sig.Params[0]

					tt := pt.Type
					if t, ok := pt.Type.(*types.Slice); ok && pt.IsVariadic {
						tt = t.Value
					}

					switch t := tt.(type) {
					case *types.Basic:
						var val any
						switch {
						default:
							val = attrVal
						case t.IsBool():
							val, _ = strconv.ParseBool(attrVal)
						case t.IsFloat():
							val, _ = strconv.ParseFloat(attrVal, 64)
						case t.IsInteger():
							i, _ := strconv.ParseInt(attrVal, 10, 64)
							val = int(i)
						}
						tag.Dot(strcase.ToCamel(attrName)).Call(jen.Lit(val))
					}
					continue
				}

				if strings.HasPrefix(attrName, ":") {
					attrName := attrName[1:]
					if strings.HasPrefix(attrName, "on") {
						methodName := strcase.ToCamel(attrName[2:])
						tag.Dot("On" + methodName).Call(jen.Id("c").Dot(attrVal))
					} else {
						tag.Dot(strcase.ToCamel(attrName)).Call(jen.Id("c").Dot(attrVal))
					}
					continue
				}

				if name, key := extractKeyValue(attrName, "aria-"); name != "" {
					tag.Dot(name).Call(jen.Lit(key), jen.Lit(attrVal))
					continue
				}
				if name, key := extractKeyValue(attrName, "data-"); name != "" {
					tag.Dot(name).Call(jen.Lit(key), jen.Lit(attrVal))
					continue
				}

				switch attrName {
				case "v-range":
					vRange = attrVal
					continue
				case "v-if":
					vIf = attrVal
					continue
				}
				switch attrName {
				default:
					attrName = strcase.ToCamel(attrName)
				case "tabindex":
					attrName = "TabIndex"
					v, _ := strconv.ParseInt(attrVal, 10, 64)
					tag.Dot(attrName).Call(jen.Lit(int(v)))
					continue
				case "id":
					attrName = "ID"
				}
				tag.Dot(attrName).Call(jen.Lit(attrVal))
			}

			if htmlContent := e.HtmlContent(); htmlContent != nil {
				for _, c := range e.HtmlContent().AllHtmlChardata() {
					if c.HTML_TEXT() != nil {
						text := strings.TrimSpace(c.HTML_TEXT().GetText())
						vars := varparser.Parse(text)

						var varIDs []string
						for _, v := range vars {
							if f := hg.s.Type.Path(v.ID); f != nil {
								t := f.Var.Type
								switch tt := t.(type) {
								case *types.Slice:
									t = tt.Value
								case *types.Array:
									t = tt.Value
								}

								if b, ok := t.(*types.Basic); ok {
									var fmtf string
									switch {
									default:
										fmtf = "%s"
									case b.IsNumeric():
										fmtf = "%d"
									case b.IsFloat():
										fmtf = "%f"
									}
									text = text[0:v.Pos.Start] + fmtf + text[v.Pos.Finish:]
									varIDs = append(varIDs, v.ID)
								}
							}
						}
						if len(varIDs) > 0 {
							callCodes = append(callCodes, jen.Qual(pkgApp, "Textf").CallFunc(func(g *jen.Group) {
								g.Lit(text)
								for _, id := range varIDs {
									g.Id("c").Dot(id)
								}
							}))
						} else {
							callCodes = append(callCodes, jen.Qual(pkgApp, "Text").Call(jen.Lit(text)))
						}
					}
				}
			}
			for _, n := range ctx.GetChildren() {
				callCodes = append(callCodes, hg.recursiveParse(n, nested+1)...)
			}

			if tagName == "routerprovider" {
				fmt.Println(ctx.GetChildren()[0].(*antlr.TerminalNodeImpl).String(), isFoundComp, len(callCodes))
			}

			if len(callCodes) > 0 {

				// if isFoundComp {
				// tag.Dot("Children").Call(callCodes...)
				// } else {
				tag.Dot("Body").Call(callCodes...)
				// }
			}
			if vIf != "" {
				tag = jen.Qual(pkgApp, "If").Call(jen.Id("c").Dot(vIf), tag)
			}
			if vRange != "" {
				idx := strings.Index(vRange, "in")
				if idx > -1 {
					key := strings.TrimSpace(vRange[:idx])
					slice := strings.TrimSpace(vRange[idx+2:])
					if f := hg.s.Type.Path(slice); f != nil {
						switch f.Var.Type.(type) {
						case *types.Array, *types.Slice:
							tag = jen.Qual(pkgApp, "Range").Call(jen.Id("c").Dot(slice)).Dot("Slice").Call(
								jen.Func().Params(jen.Id(key).Int()).Qual(pkgApp, "UI").Block(
									jen.Return(tag),
								),
							)
						case *types.Map:
							tag = jen.Qual(pkgApp, "Range").Call(jen.Id("c").Dot(slice)).Dot("Map").Call(
								jen.Func().Params(jen.Id(key).String()).Qual(pkgApp, "UI").Block(
									jen.Return(tag),
								),
							)
						}
					}
				}
			}
			codes = append(codes, tag)
		}
	}
	return
}

func NewHTML2Go(name string, qual func(path, name string) func(s *jen.Statement), structMap map[string]*gg.Struct) *HTML2Go {
	return &HTML2Go{
		s:         structMap[name],
		qual:      qual,
		structMap: structMap,
	}
}
