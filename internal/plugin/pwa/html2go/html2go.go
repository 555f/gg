package html2go

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/antlr4-go/antlr/v4"

	"github.com/555f/gg/internal/plugin/pwa/html2go/parser"
)

type htmlListener struct {
	*parser.BaseHTMLParserListener

	buf       bytes.Buffer
	index     int
	prevIndex int
}

func (s *htmlListener) EnterHtmlMisc(ctx *parser.HtmlMiscContext) {}

func (s *htmlListener) EnterHtmlChardata(ctx *parser.HtmlChardataContext) {
	n := ctx.HTML_TEXT()
	if n == nil {
		return
	}

	s.index++

	if s.index == s.prevIndex {
		fmt.Fprintf(&s.buf, ",")
	}
	// text := strings.TrimSpace(n.GetText())
	// strings.Index()

	fmt.Fprintf(&s.buf, "app.Text("+strconv.Quote(ctx.HTML_TEXT().GetText()))
}

func (s *htmlListener) ExitHtmlChardata(ctx *parser.HtmlChardataContext) {
	n := ctx.HTML_TEXT()
	if n == nil {
		return
	}

	fmt.Fprintf(&s.buf, ")")

	s.prevIndex = s.index

	s.index--
}

func (s *htmlListener) EnterHtmlElement(ctx *parser.HtmlElementContext) {
	s.index++
	tagName := strcase.ToCamel(ctx.TAG_NAME(0).GetText())
	if s.index == s.prevIndex {
		fmt.Fprintf(&s.buf, ",")
	}
	fmt.Fprintf(&s.buf, "app."+tagName+"()")

	fmt.Fprintf(&s.buf, ".Body(")
}

func (s *htmlListener) ExitHtmlElement(ctx *parser.HtmlElementContext) {
	fmt.Fprintf(&s.buf, ")")

	s.prevIndex = s.index

	s.index--
}

func (s *htmlListener) EnterHtmlContent(ctx *parser.HtmlContentContext) {}

type HTML2Go struct {
	name      string
	structMap map[string]*gg.Struct
}

func (hg *HTML2Go) Parse(html string) (code string, err error) {
	is := antlr.NewInputStream(html)
	lexer := parser.NewHTMLLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewHTMLParser(stream)
	l := &htmlListener{}

	antlr.ParseTreeWalkerDefault.Walk(l, p.HtmlDocument())
	code = fmt.Sprintf("func RenderView() app.UI {\nreturn %s\n}\n", l.buf.String())
	return
}

func NewHTML2Go(name string, structMap map[string]*gg.Struct) *HTML2Go {
	return &HTML2Go{
		name:      name,
		structMap: structMap,
	}
}
