// Code generated from HTMLParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // HTMLParser

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type HTMLParser struct {
	*antlr.BaseParser
}

var HTMLParserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func htmlparserParserInit() {
	staticData := &HTMLParserParserStaticData
	staticData.LiteralNames = []string{
		"", "", "", "", "", "", "", "", "", "", "'<'", "", "'>'", "'/>'", "'/'",
		"'='",
	}
	staticData.SymbolicNames = []string{
		"", "HTML_COMMENT", "HTML_CONDITIONAL_COMMENT", "XML", "CDATA", "DTD",
		"SCRIPTLET", "SEA_WS", "SCRIPT_OPEN", "STYLE_OPEN", "TAG_OPEN", "HTML_TEXT",
		"TAG_CLOSE", "TAG_SLASH_CLOSE", "TAG_SLASH", "TAG_EQUALS", "TAG_NAME",
		"TAG_WHITESPACE", "SCRIPT_BODY", "SCRIPT_SHORT_BODY", "STYLE_BODY",
		"STYLE_SHORT_BODY", "ATTVALUE_VALUE", "ATTRIBUTE",
	}
	staticData.RuleNames = []string{
		"htmlDocument", "scriptletOrSeaWs", "htmlElements", "htmlElement", "htmlContent",
		"htmlAttribute", "htmlChardata", "htmlMisc", "htmlComment", "script",
		"style",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 23, 128, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 1, 0, 5, 0, 24, 8, 0, 10, 0, 12, 0, 27, 9, 0, 1, 0, 3, 0, 30, 8, 0,
		1, 0, 5, 0, 33, 8, 0, 10, 0, 12, 0, 36, 9, 0, 1, 0, 3, 0, 39, 8, 0, 1,
		0, 5, 0, 42, 8, 0, 10, 0, 12, 0, 45, 9, 0, 1, 0, 5, 0, 48, 8, 0, 10, 0,
		12, 0, 51, 9, 0, 1, 1, 1, 1, 1, 2, 5, 2, 56, 8, 2, 10, 2, 12, 2, 59, 9,
		2, 1, 2, 1, 2, 5, 2, 63, 8, 2, 10, 2, 12, 2, 66, 9, 2, 1, 3, 1, 3, 1, 3,
		5, 3, 71, 8, 3, 10, 3, 12, 3, 74, 9, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1,
		3, 1, 3, 3, 3, 83, 8, 3, 1, 3, 3, 3, 86, 8, 3, 1, 3, 1, 3, 1, 3, 3, 3,
		91, 8, 3, 1, 4, 3, 4, 94, 8, 4, 1, 4, 1, 4, 1, 4, 3, 4, 99, 8, 4, 1, 4,
		3, 4, 102, 8, 4, 5, 4, 104, 8, 4, 10, 4, 12, 4, 107, 9, 4, 1, 5, 1, 5,
		1, 5, 3, 5, 112, 8, 5, 1, 6, 1, 6, 1, 7, 1, 7, 3, 7, 118, 8, 7, 1, 8, 1,
		8, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 1, 10, 0, 0, 11, 0, 2, 4, 6,
		8, 10, 12, 14, 16, 18, 20, 0, 5, 1, 0, 6, 7, 2, 0, 7, 7, 11, 11, 1, 0,
		1, 2, 1, 0, 18, 19, 1, 0, 20, 21, 137, 0, 25, 1, 0, 0, 0, 2, 52, 1, 0,
		0, 0, 4, 57, 1, 0, 0, 0, 6, 90, 1, 0, 0, 0, 8, 93, 1, 0, 0, 0, 10, 108,
		1, 0, 0, 0, 12, 113, 1, 0, 0, 0, 14, 117, 1, 0, 0, 0, 16, 119, 1, 0, 0,
		0, 18, 121, 1, 0, 0, 0, 20, 124, 1, 0, 0, 0, 22, 24, 3, 2, 1, 0, 23, 22,
		1, 0, 0, 0, 24, 27, 1, 0, 0, 0, 25, 23, 1, 0, 0, 0, 25, 26, 1, 0, 0, 0,
		26, 29, 1, 0, 0, 0, 27, 25, 1, 0, 0, 0, 28, 30, 5, 3, 0, 0, 29, 28, 1,
		0, 0, 0, 29, 30, 1, 0, 0, 0, 30, 34, 1, 0, 0, 0, 31, 33, 3, 2, 1, 0, 32,
		31, 1, 0, 0, 0, 33, 36, 1, 0, 0, 0, 34, 32, 1, 0, 0, 0, 34, 35, 1, 0, 0,
		0, 35, 38, 1, 0, 0, 0, 36, 34, 1, 0, 0, 0, 37, 39, 5, 5, 0, 0, 38, 37,
		1, 0, 0, 0, 38, 39, 1, 0, 0, 0, 39, 43, 1, 0, 0, 0, 40, 42, 3, 2, 1, 0,
		41, 40, 1, 0, 0, 0, 42, 45, 1, 0, 0, 0, 43, 41, 1, 0, 0, 0, 43, 44, 1,
		0, 0, 0, 44, 49, 1, 0, 0, 0, 45, 43, 1, 0, 0, 0, 46, 48, 3, 4, 2, 0, 47,
		46, 1, 0, 0, 0, 48, 51, 1, 0, 0, 0, 49, 47, 1, 0, 0, 0, 49, 50, 1, 0, 0,
		0, 50, 1, 1, 0, 0, 0, 51, 49, 1, 0, 0, 0, 52, 53, 7, 0, 0, 0, 53, 3, 1,
		0, 0, 0, 54, 56, 3, 14, 7, 0, 55, 54, 1, 0, 0, 0, 56, 59, 1, 0, 0, 0, 57,
		55, 1, 0, 0, 0, 57, 58, 1, 0, 0, 0, 58, 60, 1, 0, 0, 0, 59, 57, 1, 0, 0,
		0, 60, 64, 3, 6, 3, 0, 61, 63, 3, 14, 7, 0, 62, 61, 1, 0, 0, 0, 63, 66,
		1, 0, 0, 0, 64, 62, 1, 0, 0, 0, 64, 65, 1, 0, 0, 0, 65, 5, 1, 0, 0, 0,
		66, 64, 1, 0, 0, 0, 67, 68, 5, 10, 0, 0, 68, 72, 5, 16, 0, 0, 69, 71, 3,
		10, 5, 0, 70, 69, 1, 0, 0, 0, 71, 74, 1, 0, 0, 0, 72, 70, 1, 0, 0, 0, 72,
		73, 1, 0, 0, 0, 73, 85, 1, 0, 0, 0, 74, 72, 1, 0, 0, 0, 75, 82, 5, 12,
		0, 0, 76, 77, 3, 8, 4, 0, 77, 78, 5, 10, 0, 0, 78, 79, 5, 14, 0, 0, 79,
		80, 5, 16, 0, 0, 80, 81, 5, 12, 0, 0, 81, 83, 1, 0, 0, 0, 82, 76, 1, 0,
		0, 0, 82, 83, 1, 0, 0, 0, 83, 86, 1, 0, 0, 0, 84, 86, 5, 13, 0, 0, 85,
		75, 1, 0, 0, 0, 85, 84, 1, 0, 0, 0, 86, 91, 1, 0, 0, 0, 87, 91, 5, 6, 0,
		0, 88, 91, 3, 18, 9, 0, 89, 91, 3, 20, 10, 0, 90, 67, 1, 0, 0, 0, 90, 87,
		1, 0, 0, 0, 90, 88, 1, 0, 0, 0, 90, 89, 1, 0, 0, 0, 91, 7, 1, 0, 0, 0,
		92, 94, 3, 12, 6, 0, 93, 92, 1, 0, 0, 0, 93, 94, 1, 0, 0, 0, 94, 105, 1,
		0, 0, 0, 95, 99, 3, 6, 3, 0, 96, 99, 5, 4, 0, 0, 97, 99, 3, 16, 8, 0, 98,
		95, 1, 0, 0, 0, 98, 96, 1, 0, 0, 0, 98, 97, 1, 0, 0, 0, 99, 101, 1, 0,
		0, 0, 100, 102, 3, 12, 6, 0, 101, 100, 1, 0, 0, 0, 101, 102, 1, 0, 0, 0,
		102, 104, 1, 0, 0, 0, 103, 98, 1, 0, 0, 0, 104, 107, 1, 0, 0, 0, 105, 103,
		1, 0, 0, 0, 105, 106, 1, 0, 0, 0, 106, 9, 1, 0, 0, 0, 107, 105, 1, 0, 0,
		0, 108, 111, 5, 16, 0, 0, 109, 110, 5, 15, 0, 0, 110, 112, 5, 22, 0, 0,
		111, 109, 1, 0, 0, 0, 111, 112, 1, 0, 0, 0, 112, 11, 1, 0, 0, 0, 113, 114,
		7, 1, 0, 0, 114, 13, 1, 0, 0, 0, 115, 118, 3, 16, 8, 0, 116, 118, 5, 7,
		0, 0, 117, 115, 1, 0, 0, 0, 117, 116, 1, 0, 0, 0, 118, 15, 1, 0, 0, 0,
		119, 120, 7, 2, 0, 0, 120, 17, 1, 0, 0, 0, 121, 122, 5, 8, 0, 0, 122, 123,
		7, 3, 0, 0, 123, 19, 1, 0, 0, 0, 124, 125, 5, 9, 0, 0, 125, 126, 7, 4,
		0, 0, 126, 21, 1, 0, 0, 0, 18, 25, 29, 34, 38, 43, 49, 57, 64, 72, 82,
		85, 90, 93, 98, 101, 105, 111, 117,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// HTMLParserInit initializes any static state used to implement HTMLParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewHTMLParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func HTMLParserInit() {
	staticData := &HTMLParserParserStaticData
	staticData.once.Do(htmlparserParserInit)
}

// NewHTMLParser produces a new parser instance for the optional input antlr.TokenStream.
func NewHTMLParser(input antlr.TokenStream) *HTMLParser {
	HTMLParserInit()
	this := new(HTMLParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &HTMLParserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "HTMLParser.g4"

	return this
}

// HTMLParser tokens.
const (
	HTMLParserEOF                      = antlr.TokenEOF
	HTMLParserHTML_COMMENT             = 1
	HTMLParserHTML_CONDITIONAL_COMMENT = 2
	HTMLParserXML                      = 3
	HTMLParserCDATA                    = 4
	HTMLParserDTD                      = 5
	HTMLParserSCRIPTLET                = 6
	HTMLParserSEA_WS                   = 7
	HTMLParserSCRIPT_OPEN              = 8
	HTMLParserSTYLE_OPEN               = 9
	HTMLParserTAG_OPEN                 = 10
	HTMLParserHTML_TEXT                = 11
	HTMLParserTAG_CLOSE                = 12
	HTMLParserTAG_SLASH_CLOSE          = 13
	HTMLParserTAG_SLASH                = 14
	HTMLParserTAG_EQUALS               = 15
	HTMLParserTAG_NAME                 = 16
	HTMLParserTAG_WHITESPACE           = 17
	HTMLParserSCRIPT_BODY              = 18
	HTMLParserSCRIPT_SHORT_BODY        = 19
	HTMLParserSTYLE_BODY               = 20
	HTMLParserSTYLE_SHORT_BODY         = 21
	HTMLParserATTVALUE_VALUE           = 22
	HTMLParserATTRIBUTE                = 23
)

// HTMLParser rules.
const (
	HTMLParserRULE_htmlDocument     = 0
	HTMLParserRULE_scriptletOrSeaWs = 1
	HTMLParserRULE_htmlElements     = 2
	HTMLParserRULE_htmlElement      = 3
	HTMLParserRULE_htmlContent      = 4
	HTMLParserRULE_htmlAttribute    = 5
	HTMLParserRULE_htmlChardata     = 6
	HTMLParserRULE_htmlMisc         = 7
	HTMLParserRULE_htmlComment      = 8
	HTMLParserRULE_script           = 9
	HTMLParserRULE_style            = 10
)

// IHtmlDocumentContext is an interface to support dynamic dispatch.
type IHtmlDocumentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllScriptletOrSeaWs() []IScriptletOrSeaWsContext
	ScriptletOrSeaWs(i int) IScriptletOrSeaWsContext
	XML() antlr.TerminalNode
	DTD() antlr.TerminalNode
	AllHtmlElements() []IHtmlElementsContext
	HtmlElements(i int) IHtmlElementsContext

	// IsHtmlDocumentContext differentiates from other interfaces.
	IsHtmlDocumentContext()
}

type HtmlDocumentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHtmlDocumentContext() *HtmlDocumentContext {
	var p = new(HtmlDocumentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_htmlDocument
	return p
}

func InitEmptyHtmlDocumentContext(p *HtmlDocumentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_htmlDocument
}

func (*HtmlDocumentContext) IsHtmlDocumentContext() {}

func NewHtmlDocumentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HtmlDocumentContext {
	var p = new(HtmlDocumentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = HTMLParserRULE_htmlDocument

	return p
}

func (s *HtmlDocumentContext) GetParser() antlr.Parser { return s.parser }

func (s *HtmlDocumentContext) AllScriptletOrSeaWs() []IScriptletOrSeaWsContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IScriptletOrSeaWsContext); ok {
			len++
		}
	}

	tst := make([]IScriptletOrSeaWsContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IScriptletOrSeaWsContext); ok {
			tst[i] = t.(IScriptletOrSeaWsContext)
			i++
		}
	}

	return tst
}

func (s *HtmlDocumentContext) ScriptletOrSeaWs(i int) IScriptletOrSeaWsContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IScriptletOrSeaWsContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IScriptletOrSeaWsContext)
}

func (s *HtmlDocumentContext) XML() antlr.TerminalNode {
	return s.GetToken(HTMLParserXML, 0)
}

func (s *HtmlDocumentContext) DTD() antlr.TerminalNode {
	return s.GetToken(HTMLParserDTD, 0)
}

func (s *HtmlDocumentContext) AllHtmlElements() []IHtmlElementsContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IHtmlElementsContext); ok {
			len++
		}
	}

	tst := make([]IHtmlElementsContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IHtmlElementsContext); ok {
			tst[i] = t.(IHtmlElementsContext)
			i++
		}
	}

	return tst
}

func (s *HtmlDocumentContext) HtmlElements(i int) IHtmlElementsContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHtmlElementsContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHtmlElementsContext)
}

func (s *HtmlDocumentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HtmlDocumentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HtmlDocumentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.EnterHtmlDocument(s)
	}
}

func (s *HtmlDocumentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.ExitHtmlDocument(s)
	}
}

func (p *HTMLParser) HtmlDocument() (localctx IHtmlDocumentContext) {
	localctx = NewHtmlDocumentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, HTMLParserRULE_htmlDocument)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(25)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 0, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(22)
				p.ScriptletOrSeaWs()
			}

		}
		p.SetState(27)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 0, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	p.SetState(29)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == HTMLParserXML {
		{
			p.SetState(28)
			p.Match(HTMLParserXML)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	p.SetState(34)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(31)
				p.ScriptletOrSeaWs()
			}

		}
		p.SetState(36)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	p.SetState(38)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == HTMLParserDTD {
		{
			p.SetState(37)
			p.Match(HTMLParserDTD)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	p.SetState(43)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 4, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(40)
				p.ScriptletOrSeaWs()
			}

		}
		p.SetState(45)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 4, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	p.SetState(49)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1990) != 0 {
		{
			p.SetState(46)
			p.HtmlElements()
		}

		p.SetState(51)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IScriptletOrSeaWsContext is an interface to support dynamic dispatch.
type IScriptletOrSeaWsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SCRIPTLET() antlr.TerminalNode
	SEA_WS() antlr.TerminalNode

	// IsScriptletOrSeaWsContext differentiates from other interfaces.
	IsScriptletOrSeaWsContext()
}

type ScriptletOrSeaWsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyScriptletOrSeaWsContext() *ScriptletOrSeaWsContext {
	var p = new(ScriptletOrSeaWsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_scriptletOrSeaWs
	return p
}

func InitEmptyScriptletOrSeaWsContext(p *ScriptletOrSeaWsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_scriptletOrSeaWs
}

func (*ScriptletOrSeaWsContext) IsScriptletOrSeaWsContext() {}

func NewScriptletOrSeaWsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ScriptletOrSeaWsContext {
	var p = new(ScriptletOrSeaWsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = HTMLParserRULE_scriptletOrSeaWs

	return p
}

func (s *ScriptletOrSeaWsContext) GetParser() antlr.Parser { return s.parser }

func (s *ScriptletOrSeaWsContext) SCRIPTLET() antlr.TerminalNode {
	return s.GetToken(HTMLParserSCRIPTLET, 0)
}

func (s *ScriptletOrSeaWsContext) SEA_WS() antlr.TerminalNode {
	return s.GetToken(HTMLParserSEA_WS, 0)
}

func (s *ScriptletOrSeaWsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ScriptletOrSeaWsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ScriptletOrSeaWsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.EnterScriptletOrSeaWs(s)
	}
}

func (s *ScriptletOrSeaWsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.ExitScriptletOrSeaWs(s)
	}
}

func (p *HTMLParser) ScriptletOrSeaWs() (localctx IScriptletOrSeaWsContext) {
	localctx = NewScriptletOrSeaWsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, HTMLParserRULE_scriptletOrSeaWs)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(52)
		_la = p.GetTokenStream().LA(1)

		if !(_la == HTMLParserSCRIPTLET || _la == HTMLParserSEA_WS) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IHtmlElementsContext is an interface to support dynamic dispatch.
type IHtmlElementsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	HtmlElement() IHtmlElementContext
	AllHtmlMisc() []IHtmlMiscContext
	HtmlMisc(i int) IHtmlMiscContext

	// IsHtmlElementsContext differentiates from other interfaces.
	IsHtmlElementsContext()
}

type HtmlElementsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHtmlElementsContext() *HtmlElementsContext {
	var p = new(HtmlElementsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_htmlElements
	return p
}

func InitEmptyHtmlElementsContext(p *HtmlElementsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_htmlElements
}

func (*HtmlElementsContext) IsHtmlElementsContext() {}

func NewHtmlElementsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HtmlElementsContext {
	var p = new(HtmlElementsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = HTMLParserRULE_htmlElements

	return p
}

func (s *HtmlElementsContext) GetParser() antlr.Parser { return s.parser }

func (s *HtmlElementsContext) HtmlElement() IHtmlElementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHtmlElementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHtmlElementContext)
}

func (s *HtmlElementsContext) AllHtmlMisc() []IHtmlMiscContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IHtmlMiscContext); ok {
			len++
		}
	}

	tst := make([]IHtmlMiscContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IHtmlMiscContext); ok {
			tst[i] = t.(IHtmlMiscContext)
			i++
		}
	}

	return tst
}

func (s *HtmlElementsContext) HtmlMisc(i int) IHtmlMiscContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHtmlMiscContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHtmlMiscContext)
}

func (s *HtmlElementsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HtmlElementsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HtmlElementsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.EnterHtmlElements(s)
	}
}

func (s *HtmlElementsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.ExitHtmlElements(s)
	}
}

func (p *HTMLParser) HtmlElements() (localctx IHtmlElementsContext) {
	localctx = NewHtmlElementsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, HTMLParserRULE_htmlElements)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(57)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&134) != 0 {
		{
			p.SetState(54)
			p.HtmlMisc()
		}

		p.SetState(59)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(60)
		p.HtmlElement()
	}
	p.SetState(64)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 7, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(61)
				p.HtmlMisc()
			}

		}
		p.SetState(66)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 7, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IHtmlElementContext is an interface to support dynamic dispatch.
type IHtmlElementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTAG_OPEN() []antlr.TerminalNode
	TAG_OPEN(i int) antlr.TerminalNode
	AllTAG_NAME() []antlr.TerminalNode
	TAG_NAME(i int) antlr.TerminalNode
	AllTAG_CLOSE() []antlr.TerminalNode
	TAG_CLOSE(i int) antlr.TerminalNode
	TAG_SLASH_CLOSE() antlr.TerminalNode
	AllHtmlAttribute() []IHtmlAttributeContext
	HtmlAttribute(i int) IHtmlAttributeContext
	HtmlContent() IHtmlContentContext
	TAG_SLASH() antlr.TerminalNode
	SCRIPTLET() antlr.TerminalNode
	Script() IScriptContext
	Style() IStyleContext

	// IsHtmlElementContext differentiates from other interfaces.
	IsHtmlElementContext()
}

type HtmlElementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHtmlElementContext() *HtmlElementContext {
	var p = new(HtmlElementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_htmlElement
	return p
}

func InitEmptyHtmlElementContext(p *HtmlElementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_htmlElement
}

func (*HtmlElementContext) IsHtmlElementContext() {}

func NewHtmlElementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HtmlElementContext {
	var p = new(HtmlElementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = HTMLParserRULE_htmlElement

	return p
}

func (s *HtmlElementContext) GetParser() antlr.Parser { return s.parser }

func (s *HtmlElementContext) AllTAG_OPEN() []antlr.TerminalNode {
	return s.GetTokens(HTMLParserTAG_OPEN)
}

func (s *HtmlElementContext) TAG_OPEN(i int) antlr.TerminalNode {
	return s.GetToken(HTMLParserTAG_OPEN, i)
}

func (s *HtmlElementContext) AllTAG_NAME() []antlr.TerminalNode {
	return s.GetTokens(HTMLParserTAG_NAME)
}

func (s *HtmlElementContext) TAG_NAME(i int) antlr.TerminalNode {
	return s.GetToken(HTMLParserTAG_NAME, i)
}

func (s *HtmlElementContext) AllTAG_CLOSE() []antlr.TerminalNode {
	return s.GetTokens(HTMLParserTAG_CLOSE)
}

func (s *HtmlElementContext) TAG_CLOSE(i int) antlr.TerminalNode {
	return s.GetToken(HTMLParserTAG_CLOSE, i)
}

func (s *HtmlElementContext) TAG_SLASH_CLOSE() antlr.TerminalNode {
	return s.GetToken(HTMLParserTAG_SLASH_CLOSE, 0)
}

func (s *HtmlElementContext) AllHtmlAttribute() []IHtmlAttributeContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IHtmlAttributeContext); ok {
			len++
		}
	}

	tst := make([]IHtmlAttributeContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IHtmlAttributeContext); ok {
			tst[i] = t.(IHtmlAttributeContext)
			i++
		}
	}

	return tst
}

func (s *HtmlElementContext) HtmlAttribute(i int) IHtmlAttributeContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHtmlAttributeContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHtmlAttributeContext)
}

func (s *HtmlElementContext) HtmlContent() IHtmlContentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHtmlContentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHtmlContentContext)
}

func (s *HtmlElementContext) TAG_SLASH() antlr.TerminalNode {
	return s.GetToken(HTMLParserTAG_SLASH, 0)
}

func (s *HtmlElementContext) SCRIPTLET() antlr.TerminalNode {
	return s.GetToken(HTMLParserSCRIPTLET, 0)
}

func (s *HtmlElementContext) Script() IScriptContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IScriptContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IScriptContext)
}

func (s *HtmlElementContext) Style() IStyleContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStyleContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStyleContext)
}

func (s *HtmlElementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HtmlElementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HtmlElementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.EnterHtmlElement(s)
	}
}

func (s *HtmlElementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.ExitHtmlElement(s)
	}
}

func (p *HTMLParser) HtmlElement() (localctx IHtmlElementContext) {
	localctx = NewHtmlElementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, HTMLParserRULE_htmlElement)
	var _la int

	p.SetState(90)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case HTMLParserTAG_OPEN:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(67)
			p.Match(HTMLParserTAG_OPEN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(68)
			p.Match(HTMLParserTAG_NAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(72)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == HTMLParserTAG_NAME {
			{
				p.SetState(69)
				p.HtmlAttribute()
			}

			p.SetState(74)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		p.SetState(85)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case HTMLParserTAG_CLOSE:
			{
				p.SetState(75)
				p.Match(HTMLParserTAG_CLOSE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(82)
			p.GetErrorHandler().Sync(p)

			if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 9, p.GetParserRuleContext()) == 1 {
				{
					p.SetState(76)
					p.HtmlContent()
				}
				{
					p.SetState(77)
					p.Match(HTMLParserTAG_OPEN)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(78)
					p.Match(HTMLParserTAG_SLASH)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(79)
					p.Match(HTMLParserTAG_NAME)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(80)
					p.Match(HTMLParserTAG_CLOSE)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			} else if p.HasError() { // JIM
				goto errorExit
			}

		case HTMLParserTAG_SLASH_CLOSE:
			{
				p.SetState(84)
				p.Match(HTMLParserTAG_SLASH_CLOSE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

	case HTMLParserSCRIPTLET:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(87)
			p.Match(HTMLParserSCRIPTLET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case HTMLParserSCRIPT_OPEN:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(88)
			p.Script()
		}

	case HTMLParserSTYLE_OPEN:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(89)
			p.Style()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IHtmlContentContext is an interface to support dynamic dispatch.
type IHtmlContentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllHtmlChardata() []IHtmlChardataContext
	HtmlChardata(i int) IHtmlChardataContext
	AllHtmlElement() []IHtmlElementContext
	HtmlElement(i int) IHtmlElementContext
	AllCDATA() []antlr.TerminalNode
	CDATA(i int) antlr.TerminalNode
	AllHtmlComment() []IHtmlCommentContext
	HtmlComment(i int) IHtmlCommentContext

	// IsHtmlContentContext differentiates from other interfaces.
	IsHtmlContentContext()
}

type HtmlContentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHtmlContentContext() *HtmlContentContext {
	var p = new(HtmlContentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_htmlContent
	return p
}

func InitEmptyHtmlContentContext(p *HtmlContentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_htmlContent
}

func (*HtmlContentContext) IsHtmlContentContext() {}

func NewHtmlContentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HtmlContentContext {
	var p = new(HtmlContentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = HTMLParserRULE_htmlContent

	return p
}

func (s *HtmlContentContext) GetParser() antlr.Parser { return s.parser }

func (s *HtmlContentContext) AllHtmlChardata() []IHtmlChardataContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IHtmlChardataContext); ok {
			len++
		}
	}

	tst := make([]IHtmlChardataContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IHtmlChardataContext); ok {
			tst[i] = t.(IHtmlChardataContext)
			i++
		}
	}

	return tst
}

func (s *HtmlContentContext) HtmlChardata(i int) IHtmlChardataContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHtmlChardataContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHtmlChardataContext)
}

func (s *HtmlContentContext) AllHtmlElement() []IHtmlElementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IHtmlElementContext); ok {
			len++
		}
	}

	tst := make([]IHtmlElementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IHtmlElementContext); ok {
			tst[i] = t.(IHtmlElementContext)
			i++
		}
	}

	return tst
}

func (s *HtmlContentContext) HtmlElement(i int) IHtmlElementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHtmlElementContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHtmlElementContext)
}

func (s *HtmlContentContext) AllCDATA() []antlr.TerminalNode {
	return s.GetTokens(HTMLParserCDATA)
}

func (s *HtmlContentContext) CDATA(i int) antlr.TerminalNode {
	return s.GetToken(HTMLParserCDATA, i)
}

func (s *HtmlContentContext) AllHtmlComment() []IHtmlCommentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IHtmlCommentContext); ok {
			len++
		}
	}

	tst := make([]IHtmlCommentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IHtmlCommentContext); ok {
			tst[i] = t.(IHtmlCommentContext)
			i++
		}
	}

	return tst
}

func (s *HtmlContentContext) HtmlComment(i int) IHtmlCommentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHtmlCommentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHtmlCommentContext)
}

func (s *HtmlContentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HtmlContentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HtmlContentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.EnterHtmlContent(s)
	}
}

func (s *HtmlContentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.ExitHtmlContent(s)
	}
}

func (p *HTMLParser) HtmlContent() (localctx IHtmlContentContext) {
	localctx = NewHtmlContentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, HTMLParserRULE_htmlContent)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(93)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == HTMLParserSEA_WS || _la == HTMLParserHTML_TEXT {
		{
			p.SetState(92)
			p.HtmlChardata()
		}

	}
	p.SetState(105)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 15, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(98)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetTokenStream().LA(1) {
			case HTMLParserSCRIPTLET, HTMLParserSCRIPT_OPEN, HTMLParserSTYLE_OPEN, HTMLParserTAG_OPEN:
				{
					p.SetState(95)
					p.HtmlElement()
				}

			case HTMLParserCDATA:
				{
					p.SetState(96)
					p.Match(HTMLParserCDATA)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			case HTMLParserHTML_COMMENT, HTMLParserHTML_CONDITIONAL_COMMENT:
				{
					p.SetState(97)
					p.HtmlComment()
				}

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}
			p.SetState(101)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			if _la == HTMLParserSEA_WS || _la == HTMLParserHTML_TEXT {
				{
					p.SetState(100)
					p.HtmlChardata()
				}

			}

		}
		p.SetState(107)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 15, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IHtmlAttributeContext is an interface to support dynamic dispatch.
type IHtmlAttributeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TAG_NAME() antlr.TerminalNode
	TAG_EQUALS() antlr.TerminalNode
	ATTVALUE_VALUE() antlr.TerminalNode

	// IsHtmlAttributeContext differentiates from other interfaces.
	IsHtmlAttributeContext()
}

type HtmlAttributeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHtmlAttributeContext() *HtmlAttributeContext {
	var p = new(HtmlAttributeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_htmlAttribute
	return p
}

func InitEmptyHtmlAttributeContext(p *HtmlAttributeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_htmlAttribute
}

func (*HtmlAttributeContext) IsHtmlAttributeContext() {}

func NewHtmlAttributeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HtmlAttributeContext {
	var p = new(HtmlAttributeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = HTMLParserRULE_htmlAttribute

	return p
}

func (s *HtmlAttributeContext) GetParser() antlr.Parser { return s.parser }

func (s *HtmlAttributeContext) TAG_NAME() antlr.TerminalNode {
	return s.GetToken(HTMLParserTAG_NAME, 0)
}

func (s *HtmlAttributeContext) TAG_EQUALS() antlr.TerminalNode {
	return s.GetToken(HTMLParserTAG_EQUALS, 0)
}

func (s *HtmlAttributeContext) ATTVALUE_VALUE() antlr.TerminalNode {
	return s.GetToken(HTMLParserATTVALUE_VALUE, 0)
}

func (s *HtmlAttributeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HtmlAttributeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HtmlAttributeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.EnterHtmlAttribute(s)
	}
}

func (s *HtmlAttributeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.ExitHtmlAttribute(s)
	}
}

func (p *HTMLParser) HtmlAttribute() (localctx IHtmlAttributeContext) {
	localctx = NewHtmlAttributeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, HTMLParserRULE_htmlAttribute)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(108)
		p.Match(HTMLParserTAG_NAME)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(111)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == HTMLParserTAG_EQUALS {
		{
			p.SetState(109)
			p.Match(HTMLParserTAG_EQUALS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(110)
			p.Match(HTMLParserATTVALUE_VALUE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IHtmlChardataContext is an interface to support dynamic dispatch.
type IHtmlChardataContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	HTML_TEXT() antlr.TerminalNode
	SEA_WS() antlr.TerminalNode

	// IsHtmlChardataContext differentiates from other interfaces.
	IsHtmlChardataContext()
}

type HtmlChardataContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHtmlChardataContext() *HtmlChardataContext {
	var p = new(HtmlChardataContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_htmlChardata
	return p
}

func InitEmptyHtmlChardataContext(p *HtmlChardataContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_htmlChardata
}

func (*HtmlChardataContext) IsHtmlChardataContext() {}

func NewHtmlChardataContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HtmlChardataContext {
	var p = new(HtmlChardataContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = HTMLParserRULE_htmlChardata

	return p
}

func (s *HtmlChardataContext) GetParser() antlr.Parser { return s.parser }

func (s *HtmlChardataContext) HTML_TEXT() antlr.TerminalNode {
	return s.GetToken(HTMLParserHTML_TEXT, 0)
}

func (s *HtmlChardataContext) SEA_WS() antlr.TerminalNode {
	return s.GetToken(HTMLParserSEA_WS, 0)
}

func (s *HtmlChardataContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HtmlChardataContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HtmlChardataContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.EnterHtmlChardata(s)
	}
}

func (s *HtmlChardataContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.ExitHtmlChardata(s)
	}
}

func (p *HTMLParser) HtmlChardata() (localctx IHtmlChardataContext) {
	localctx = NewHtmlChardataContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, HTMLParserRULE_htmlChardata)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(113)
		_la = p.GetTokenStream().LA(1)

		if !(_la == HTMLParserSEA_WS || _la == HTMLParserHTML_TEXT) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IHtmlMiscContext is an interface to support dynamic dispatch.
type IHtmlMiscContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	HtmlComment() IHtmlCommentContext
	SEA_WS() antlr.TerminalNode

	// IsHtmlMiscContext differentiates from other interfaces.
	IsHtmlMiscContext()
}

type HtmlMiscContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHtmlMiscContext() *HtmlMiscContext {
	var p = new(HtmlMiscContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_htmlMisc
	return p
}

func InitEmptyHtmlMiscContext(p *HtmlMiscContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_htmlMisc
}

func (*HtmlMiscContext) IsHtmlMiscContext() {}

func NewHtmlMiscContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HtmlMiscContext {
	var p = new(HtmlMiscContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = HTMLParserRULE_htmlMisc

	return p
}

func (s *HtmlMiscContext) GetParser() antlr.Parser { return s.parser }

func (s *HtmlMiscContext) HtmlComment() IHtmlCommentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHtmlCommentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHtmlCommentContext)
}

func (s *HtmlMiscContext) SEA_WS() antlr.TerminalNode {
	return s.GetToken(HTMLParserSEA_WS, 0)
}

func (s *HtmlMiscContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HtmlMiscContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HtmlMiscContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.EnterHtmlMisc(s)
	}
}

func (s *HtmlMiscContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.ExitHtmlMisc(s)
	}
}

func (p *HTMLParser) HtmlMisc() (localctx IHtmlMiscContext) {
	localctx = NewHtmlMiscContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, HTMLParserRULE_htmlMisc)
	p.SetState(117)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case HTMLParserHTML_COMMENT, HTMLParserHTML_CONDITIONAL_COMMENT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(115)
			p.HtmlComment()
		}

	case HTMLParserSEA_WS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(116)
			p.Match(HTMLParserSEA_WS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IHtmlCommentContext is an interface to support dynamic dispatch.
type IHtmlCommentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	HTML_COMMENT() antlr.TerminalNode
	HTML_CONDITIONAL_COMMENT() antlr.TerminalNode

	// IsHtmlCommentContext differentiates from other interfaces.
	IsHtmlCommentContext()
}

type HtmlCommentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHtmlCommentContext() *HtmlCommentContext {
	var p = new(HtmlCommentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_htmlComment
	return p
}

func InitEmptyHtmlCommentContext(p *HtmlCommentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_htmlComment
}

func (*HtmlCommentContext) IsHtmlCommentContext() {}

func NewHtmlCommentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HtmlCommentContext {
	var p = new(HtmlCommentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = HTMLParserRULE_htmlComment

	return p
}

func (s *HtmlCommentContext) GetParser() antlr.Parser { return s.parser }

func (s *HtmlCommentContext) HTML_COMMENT() antlr.TerminalNode {
	return s.GetToken(HTMLParserHTML_COMMENT, 0)
}

func (s *HtmlCommentContext) HTML_CONDITIONAL_COMMENT() antlr.TerminalNode {
	return s.GetToken(HTMLParserHTML_CONDITIONAL_COMMENT, 0)
}

func (s *HtmlCommentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HtmlCommentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HtmlCommentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.EnterHtmlComment(s)
	}
}

func (s *HtmlCommentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.ExitHtmlComment(s)
	}
}

func (p *HTMLParser) HtmlComment() (localctx IHtmlCommentContext) {
	localctx = NewHtmlCommentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, HTMLParserRULE_htmlComment)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(119)
		_la = p.GetTokenStream().LA(1)

		if !(_la == HTMLParserHTML_COMMENT || _la == HTMLParserHTML_CONDITIONAL_COMMENT) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IScriptContext is an interface to support dynamic dispatch.
type IScriptContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SCRIPT_OPEN() antlr.TerminalNode
	SCRIPT_BODY() antlr.TerminalNode
	SCRIPT_SHORT_BODY() antlr.TerminalNode

	// IsScriptContext differentiates from other interfaces.
	IsScriptContext()
}

type ScriptContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyScriptContext() *ScriptContext {
	var p = new(ScriptContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_script
	return p
}

func InitEmptyScriptContext(p *ScriptContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_script
}

func (*ScriptContext) IsScriptContext() {}

func NewScriptContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ScriptContext {
	var p = new(ScriptContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = HTMLParserRULE_script

	return p
}

func (s *ScriptContext) GetParser() antlr.Parser { return s.parser }

func (s *ScriptContext) SCRIPT_OPEN() antlr.TerminalNode {
	return s.GetToken(HTMLParserSCRIPT_OPEN, 0)
}

func (s *ScriptContext) SCRIPT_BODY() antlr.TerminalNode {
	return s.GetToken(HTMLParserSCRIPT_BODY, 0)
}

func (s *ScriptContext) SCRIPT_SHORT_BODY() antlr.TerminalNode {
	return s.GetToken(HTMLParserSCRIPT_SHORT_BODY, 0)
}

func (s *ScriptContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ScriptContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ScriptContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.EnterScript(s)
	}
}

func (s *ScriptContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.ExitScript(s)
	}
}

func (p *HTMLParser) Script() (localctx IScriptContext) {
	localctx = NewScriptContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, HTMLParserRULE_script)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(121)
		p.Match(HTMLParserSCRIPT_OPEN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(122)
		_la = p.GetTokenStream().LA(1)

		if !(_la == HTMLParserSCRIPT_BODY || _la == HTMLParserSCRIPT_SHORT_BODY) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IStyleContext is an interface to support dynamic dispatch.
type IStyleContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STYLE_OPEN() antlr.TerminalNode
	STYLE_BODY() antlr.TerminalNode
	STYLE_SHORT_BODY() antlr.TerminalNode

	// IsStyleContext differentiates from other interfaces.
	IsStyleContext()
}

type StyleContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStyleContext() *StyleContext {
	var p = new(StyleContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_style
	return p
}

func InitEmptyStyleContext(p *StyleContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = HTMLParserRULE_style
}

func (*StyleContext) IsStyleContext() {}

func NewStyleContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StyleContext {
	var p = new(StyleContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = HTMLParserRULE_style

	return p
}

func (s *StyleContext) GetParser() antlr.Parser { return s.parser }

func (s *StyleContext) STYLE_OPEN() antlr.TerminalNode {
	return s.GetToken(HTMLParserSTYLE_OPEN, 0)
}

func (s *StyleContext) STYLE_BODY() antlr.TerminalNode {
	return s.GetToken(HTMLParserSTYLE_BODY, 0)
}

func (s *StyleContext) STYLE_SHORT_BODY() antlr.TerminalNode {
	return s.GetToken(HTMLParserSTYLE_SHORT_BODY, 0)
}

func (s *StyleContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StyleContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StyleContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.EnterStyle(s)
	}
}

func (s *StyleContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HTMLParserListener); ok {
		listenerT.ExitStyle(s)
	}
}

func (p *HTMLParser) Style() (localctx IStyleContext) {
	localctx = NewStyleContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, HTMLParserRULE_style)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(124)
		p.Match(HTMLParserSTYLE_OPEN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(125)
		_la = p.GetTokenStream().LA(1)

		if !(_la == HTMLParserSTYLE_BODY || _la == HTMLParserSTYLE_SHORT_BODY) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}
