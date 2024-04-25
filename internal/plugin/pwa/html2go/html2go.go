package html2go

import (
	"database/sql"
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
	codes = hg.recursiveParse(p.HtmlElements(), 0)

	return
}

func htmlToString(t antlr.Tree) (tagText string) {
	if r, ok := t.(antlr.RuleNode); ok {
		ctx := r.GetRuleContext().(antlr.ParserRuleContext)
		switch e := ctx.(type) {
		default:
			for _, n := range ctx.GetChildren() {
				tagText += htmlToString(n)
			}
		case *parser.HtmlChardataContext:
			if e.HTML_TEXT() != nil {
				return strings.TrimSpace(e.HTML_TEXT().GetText())
			}
			return ""
		case *parser.HtmlElementContext:
			tagName := e.TAG_NAME(0).GetText()

			tagText = "<" + tagName + " "
			for i, a := range e.AllHtmlAttribute() {
				if i > 0 {
					tagText += " "
				}
				tagText += a.TAG_NAME().GetText() + "=" + a.ATTVALUE_VALUE().GetText()
			}

			if len(e.GetChildren()) > 0 {
				tagText += ">"
				for _, n := range ctx.GetChildren() {
					tagText += htmlToString(n)
				}
				tagText += "</" + tagName + ">"
			} else {
				tagText += "/>"
			}
		}
	}
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
		case *parser.HtmlChardataContext:
			if e.HTML_TEXT() != nil {
				text := strings.TrimSpace(e.HTML_TEXT().GetText())
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
					codes = append(codes, jen.Qual(pkgApp, "Textf").CallFunc(func(g *jen.Group) {
						g.Lit(text)
						for _, id := range varIDs {
							g.Id("c").Dot(id)
						}
					}))
				} else {
					codes = append(codes, jen.Qual(pkgApp, "Text").Call(jen.Lit(text)))
				}
			}
		case *parser.HtmlElementContext:
			// fmt.Println(strings.Repeat(" ", nested), tagName)
			tagName := e.TAG_NAME(0).GetText()
			if tagName == "svg" {
				tagText := htmlToString(e)
				codes = append(codes, jen.Qual(pkgApp, "Raw").Call(jen.Lit(tagText)))
				return
			}

			tag := jen.Qual(pkgApp, strcase.ToCamel(tagName)).Call()

			var (
				callCodes   = make([]jen.Code, 0, 3000)
				vRange, vIf string
				cmpMethods  = make([]*types.Func, 0, 128)
			)

			if ss, ok := hg.structMap[tagName]; ok {

				cmpMethods = ss.Named.Methods

				tag = jen.Op("&").Do(hg.qual(ss.Named.Pkg.Path, ss.Named.Name))
				tag.Values()
				tag = jen.Call(tag)
			} else if p := hg.s.Type.Path(strcase.ToLowerCamel(tagName)); p != nil {

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
				var attrVal sql.NullString
				attrName := a.TAG_NAME().GetText()
				if a.ATTVALUE_VALUE() != nil {
					attrVal = sql.NullString{Valid: true, String: trimQuote(a.ATTVALUE_VALUE().GetText())}
				}

				if f, ok := cmpMethodMap[strcase.ToCamel(attrName)]; ok && f.Sig.Params.Len() > 0 {
					pt := f.Sig.Params[0]

					tt := pt.Type
					if t, ok := pt.Type.(*types.Slice); ok && pt.IsVariadic {
						tt = t.Value
					}
					var codeVal jen.Code
					switch t := tt.(type) {
					default:
						codeVal = jen.Id(attrVal.String)
					case *types.Basic:
						switch {
						default:
							codeVal = jen.Lit(attrVal.String)
						case t.IsBool():
							val, _ := strconv.ParseBool(attrVal.String)
							codeVal = jen.Lit(val)
						case t.IsFloat():
							val, _ := strconv.ParseFloat(attrVal.String, 64)
							codeVal = jen.Lit(val)
						case t.IsInteger():
							i, _ := strconv.ParseInt(attrVal.String, 10, 64)
							codeVal = jen.Lit(int(i))
						}
						tag.Dot(strcase.ToCamel(attrName)).Call(codeVal)
					}
					continue
				}

				if strings.HasPrefix(attrName, ":") {
					attrName := attrName[1:]
					if strings.HasPrefix(attrName, "on") {
						methodName := strcase.ToCamel(attrName[2:])
						tag.Dot("On" + methodName).Call(jen.Id("c").Dot(attrVal.String))
					} else {
						tag.Dot(strcase.ToCamel(attrName)).Call(jen.Id("c").Dot(attrVal.String))
					}
					continue
				}
				if attrName == "style" {
					attrVal.String = strings.Trim(attrVal.String, ";")
					attrVal.String = strings.TrimSpace(attrVal.String)

					styles := strings.Split(attrVal.String, ";")
					for _, style := range styles {
						parts := strings.Split(style, ":")
						key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
						tag.Dot("Style").Call(jen.Lit(key), jen.Lit(val))
					}
					continue
				}
				if key := extractKeyName(attrName, "aria-"); key != "" {
					tag.Dot("Aria").Call(jen.Lit(key), jen.Lit(attrVal.String))
					continue
				}
				if key := extractKeyName(attrName, "data-"); key != "" {
					tag.Dot("DataSet").Call(jen.Lit(key), jen.Lit(attrVal.String))
					continue
				}

				switch attrName {
				case "v-range":
					vRange = attrVal.String
					continue
				case "v-if":
					vIf = attrVal.String
					continue
				}
				switch attrName {
				default:
					attrName = strcase.ToCamel(attrName)
				case "autocomplete", "required":
					continue
				case "tabindex":
					attrName = "TabIndex"
					v, _ := strconv.ParseInt(attrVal.String, 10, 64)
					tag.Dot(attrName).Call(jen.Lit(int(v)))
					continue
				case "id":
					attrName = "ID"
				}
				tag.Dot(attrName).CallFunc(func(g *jen.Group) {
					if attrVal.Valid {
						g.Lit(attrVal.String)
					}
				})
			}
			for _, n := range ctx.GetChildren() {
				callCodes = append(callCodes, hg.recursiveParse(n, nested+1)...)
			}
			if len(callCodes) > 0 {
				tag.Dot("Body").Call(callCodes...)
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
