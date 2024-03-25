// Generated from /home/vitaly/Documents/work/my/gg/internal/plugin/pwa/html2go/HTMLParser.g4 by ANTLR 4.13.1
import org.antlr.v4.runtime.atn.*;
import org.antlr.v4.runtime.dfa.DFA;
import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.misc.*;
import org.antlr.v4.runtime.tree.*;
import java.util.List;
import java.util.Iterator;
import java.util.ArrayList;

@SuppressWarnings({"all", "warnings", "unchecked", "unused", "cast", "CheckReturnValue"})
public class HTMLParser extends Parser {
	static { RuntimeMetaData.checkVersion("4.13.1", RuntimeMetaData.VERSION); }

	protected static final DFA[] _decisionToDFA;
	protected static final PredictionContextCache _sharedContextCache =
		new PredictionContextCache();
	public static final int
		HTML_COMMENT=1, HTML_CONDITIONAL_COMMENT=2, XML=3, CDATA=4, DTD=5, SCRIPTLET=6, 
		SEA_WS=7, SCRIPT_OPEN=8, STYLE_OPEN=9, TAG_OPEN=10, HTML_TEXT=11, TAG_CLOSE=12, 
		TAG_SLASH_CLOSE=13, TAG_SLASH=14, TAG_EQUALS=15, TAG_NAME=16, TAG_WHITESPACE=17, 
		SCRIPT_BODY=18, SCRIPT_SHORT_BODY=19, STYLE_BODY=20, STYLE_SHORT_BODY=21, 
		ATTVALUE_VALUE=22, ATTRIBUTE=23;
	public static final int
		RULE_htmlDocument = 0, RULE_scriptletOrSeaWs = 1, RULE_htmlElements = 2, 
		RULE_htmlElement = 3, RULE_htmlContent = 4, RULE_htmlAttribute = 5, RULE_htmlChardata = 6, 
		RULE_htmlMisc = 7, RULE_htmlComment = 8, RULE_script = 9, RULE_style = 10;
	private static String[] makeRuleNames() {
		return new String[] {
			"htmlDocument", "scriptletOrSeaWs", "htmlElements", "htmlElement", "htmlContent", 
			"htmlAttribute", "htmlChardata", "htmlMisc", "htmlComment", "script", 
			"style"
		};
	}
	public static final String[] ruleNames = makeRuleNames();

	private static String[] makeLiteralNames() {
		return new String[] {
			null, null, null, null, null, null, null, null, null, null, "'<'", null, 
			"'>'", "'/>'", "'/'", "'='"
		};
	}
	private static final String[] _LITERAL_NAMES = makeLiteralNames();
	private static String[] makeSymbolicNames() {
		return new String[] {
			null, "HTML_COMMENT", "HTML_CONDITIONAL_COMMENT", "XML", "CDATA", "DTD", 
			"SCRIPTLET", "SEA_WS", "SCRIPT_OPEN", "STYLE_OPEN", "TAG_OPEN", "HTML_TEXT", 
			"TAG_CLOSE", "TAG_SLASH_CLOSE", "TAG_SLASH", "TAG_EQUALS", "TAG_NAME", 
			"TAG_WHITESPACE", "SCRIPT_BODY", "SCRIPT_SHORT_BODY", "STYLE_BODY", "STYLE_SHORT_BODY", 
			"ATTVALUE_VALUE", "ATTRIBUTE"
		};
	}
	private static final String[] _SYMBOLIC_NAMES = makeSymbolicNames();
	public static final Vocabulary VOCABULARY = new VocabularyImpl(_LITERAL_NAMES, _SYMBOLIC_NAMES);

	/**
	 * @deprecated Use {@link #VOCABULARY} instead.
	 */
	@Deprecated
	public static final String[] tokenNames;
	static {
		tokenNames = new String[_SYMBOLIC_NAMES.length];
		for (int i = 0; i < tokenNames.length; i++) {
			tokenNames[i] = VOCABULARY.getLiteralName(i);
			if (tokenNames[i] == null) {
				tokenNames[i] = VOCABULARY.getSymbolicName(i);
			}

			if (tokenNames[i] == null) {
				tokenNames[i] = "<INVALID>";
			}
		}
	}

	@Override
	@Deprecated
	public String[] getTokenNames() {
		return tokenNames;
	}

	@Override

	public Vocabulary getVocabulary() {
		return VOCABULARY;
	}

	@Override
	public String getGrammarFileName() { return "HTMLParser.g4"; }

	@Override
	public String[] getRuleNames() { return ruleNames; }

	@Override
	public String getSerializedATN() { return _serializedATN; }

	@Override
	public ATN getATN() { return _ATN; }

	public HTMLParser(TokenStream input) {
		super(input);
		_interp = new ParserATNSimulator(this,_ATN,_decisionToDFA,_sharedContextCache);
	}

	@SuppressWarnings("CheckReturnValue")
	public static class HtmlDocumentContext extends ParserRuleContext {
		public List<ScriptletOrSeaWsContext> scriptletOrSeaWs() {
			return getRuleContexts(ScriptletOrSeaWsContext.class);
		}
		public ScriptletOrSeaWsContext scriptletOrSeaWs(int i) {
			return getRuleContext(ScriptletOrSeaWsContext.class,i);
		}
		public TerminalNode XML() { return getToken(HTMLParser.XML, 0); }
		public TerminalNode DTD() { return getToken(HTMLParser.DTD, 0); }
		public List<HtmlElementsContext> htmlElements() {
			return getRuleContexts(HtmlElementsContext.class);
		}
		public HtmlElementsContext htmlElements(int i) {
			return getRuleContext(HtmlElementsContext.class,i);
		}
		public HtmlDocumentContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_htmlDocument; }
	}

	public final HtmlDocumentContext htmlDocument() throws RecognitionException {
		HtmlDocumentContext _localctx = new HtmlDocumentContext(_ctx, getState());
		enterRule(_localctx, 0, RULE_htmlDocument);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(25);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,0,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(22);
					scriptletOrSeaWs();
					}
					} 
				}
				setState(27);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,0,_ctx);
			}
			setState(29);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==XML) {
				{
				setState(28);
				match(XML);
				}
			}

			setState(34);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,2,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(31);
					scriptletOrSeaWs();
					}
					} 
				}
				setState(36);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,2,_ctx);
			}
			setState(38);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==DTD) {
				{
				setState(37);
				match(DTD);
				}
			}

			setState(43);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,4,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(40);
					scriptletOrSeaWs();
					}
					} 
				}
				setState(45);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,4,_ctx);
			}
			setState(49);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 1990L) != 0)) {
				{
				{
				setState(46);
				htmlElements();
				}
				}
				setState(51);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ScriptletOrSeaWsContext extends ParserRuleContext {
		public TerminalNode SCRIPTLET() { return getToken(HTMLParser.SCRIPTLET, 0); }
		public TerminalNode SEA_WS() { return getToken(HTMLParser.SEA_WS, 0); }
		public ScriptletOrSeaWsContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_scriptletOrSeaWs; }
	}

	public final ScriptletOrSeaWsContext scriptletOrSeaWs() throws RecognitionException {
		ScriptletOrSeaWsContext _localctx = new ScriptletOrSeaWsContext(_ctx, getState());
		enterRule(_localctx, 2, RULE_scriptletOrSeaWs);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(52);
			_la = _input.LA(1);
			if ( !(_la==SCRIPTLET || _la==SEA_WS) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class HtmlElementsContext extends ParserRuleContext {
		public HtmlElementContext htmlElement() {
			return getRuleContext(HtmlElementContext.class,0);
		}
		public List<HtmlMiscContext> htmlMisc() {
			return getRuleContexts(HtmlMiscContext.class);
		}
		public HtmlMiscContext htmlMisc(int i) {
			return getRuleContext(HtmlMiscContext.class,i);
		}
		public HtmlElementsContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_htmlElements; }
	}

	public final HtmlElementsContext htmlElements() throws RecognitionException {
		HtmlElementsContext _localctx = new HtmlElementsContext(_ctx, getState());
		enterRule(_localctx, 4, RULE_htmlElements);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(57);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 134L) != 0)) {
				{
				{
				setState(54);
				htmlMisc();
				}
				}
				setState(59);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(60);
			htmlElement();
			setState(64);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,7,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(61);
					htmlMisc();
					}
					} 
				}
				setState(66);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,7,_ctx);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class HtmlElementContext extends ParserRuleContext {
		public List<TerminalNode> TAG_OPEN() { return getTokens(HTMLParser.TAG_OPEN); }
		public TerminalNode TAG_OPEN(int i) {
			return getToken(HTMLParser.TAG_OPEN, i);
		}
		public List<TerminalNode> TAG_NAME() { return getTokens(HTMLParser.TAG_NAME); }
		public TerminalNode TAG_NAME(int i) {
			return getToken(HTMLParser.TAG_NAME, i);
		}
		public List<TerminalNode> TAG_CLOSE() { return getTokens(HTMLParser.TAG_CLOSE); }
		public TerminalNode TAG_CLOSE(int i) {
			return getToken(HTMLParser.TAG_CLOSE, i);
		}
		public TerminalNode TAG_SLASH_CLOSE() { return getToken(HTMLParser.TAG_SLASH_CLOSE, 0); }
		public List<HtmlAttributeContext> htmlAttribute() {
			return getRuleContexts(HtmlAttributeContext.class);
		}
		public HtmlAttributeContext htmlAttribute(int i) {
			return getRuleContext(HtmlAttributeContext.class,i);
		}
		public HtmlContentContext htmlContent() {
			return getRuleContext(HtmlContentContext.class,0);
		}
		public TerminalNode TAG_SLASH() { return getToken(HTMLParser.TAG_SLASH, 0); }
		public TerminalNode SCRIPTLET() { return getToken(HTMLParser.SCRIPTLET, 0); }
		public ScriptContext script() {
			return getRuleContext(ScriptContext.class,0);
		}
		public StyleContext style() {
			return getRuleContext(StyleContext.class,0);
		}
		public HtmlElementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_htmlElement; }
	}

	public final HtmlElementContext htmlElement() throws RecognitionException {
		HtmlElementContext _localctx = new HtmlElementContext(_ctx, getState());
		enterRule(_localctx, 6, RULE_htmlElement);
		int _la;
		try {
			setState(90);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case TAG_OPEN:
				enterOuterAlt(_localctx, 1);
				{
				setState(67);
				match(TAG_OPEN);
				setState(68);
				match(TAG_NAME);
				setState(72);
				_errHandler.sync(this);
				_la = _input.LA(1);
				while (_la==TAG_NAME) {
					{
					{
					setState(69);
					htmlAttribute();
					}
					}
					setState(74);
					_errHandler.sync(this);
					_la = _input.LA(1);
				}
				setState(85);
				_errHandler.sync(this);
				switch (_input.LA(1)) {
				case TAG_CLOSE:
					{
					setState(75);
					match(TAG_CLOSE);
					setState(82);
					_errHandler.sync(this);
					switch ( getInterpreter().adaptivePredict(_input,9,_ctx) ) {
					case 1:
						{
						setState(76);
						htmlContent();
						setState(77);
						match(TAG_OPEN);
						setState(78);
						match(TAG_SLASH);
						setState(79);
						match(TAG_NAME);
						setState(80);
						match(TAG_CLOSE);
						}
						break;
					}
					}
					break;
				case TAG_SLASH_CLOSE:
					{
					setState(84);
					match(TAG_SLASH_CLOSE);
					}
					break;
				default:
					throw new NoViableAltException(this);
				}
				}
				break;
			case SCRIPTLET:
				enterOuterAlt(_localctx, 2);
				{
				setState(87);
				match(SCRIPTLET);
				}
				break;
			case SCRIPT_OPEN:
				enterOuterAlt(_localctx, 3);
				{
				setState(88);
				script();
				}
				break;
			case STYLE_OPEN:
				enterOuterAlt(_localctx, 4);
				{
				setState(89);
				style();
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class HtmlContentContext extends ParserRuleContext {
		public List<HtmlChardataContext> htmlChardata() {
			return getRuleContexts(HtmlChardataContext.class);
		}
		public HtmlChardataContext htmlChardata(int i) {
			return getRuleContext(HtmlChardataContext.class,i);
		}
		public List<HtmlElementContext> htmlElement() {
			return getRuleContexts(HtmlElementContext.class);
		}
		public HtmlElementContext htmlElement(int i) {
			return getRuleContext(HtmlElementContext.class,i);
		}
		public List<TerminalNode> CDATA() { return getTokens(HTMLParser.CDATA); }
		public TerminalNode CDATA(int i) {
			return getToken(HTMLParser.CDATA, i);
		}
		public List<HtmlCommentContext> htmlComment() {
			return getRuleContexts(HtmlCommentContext.class);
		}
		public HtmlCommentContext htmlComment(int i) {
			return getRuleContext(HtmlCommentContext.class,i);
		}
		public HtmlContentContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_htmlContent; }
	}

	public final HtmlContentContext htmlContent() throws RecognitionException {
		HtmlContentContext _localctx = new HtmlContentContext(_ctx, getState());
		enterRule(_localctx, 8, RULE_htmlContent);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(93);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==SEA_WS || _la==HTML_TEXT) {
				{
				setState(92);
				htmlChardata();
				}
			}

			setState(105);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,15,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(98);
					_errHandler.sync(this);
					switch (_input.LA(1)) {
					case SCRIPTLET:
					case SCRIPT_OPEN:
					case STYLE_OPEN:
					case TAG_OPEN:
						{
						setState(95);
						htmlElement();
						}
						break;
					case CDATA:
						{
						setState(96);
						match(CDATA);
						}
						break;
					case HTML_COMMENT:
					case HTML_CONDITIONAL_COMMENT:
						{
						setState(97);
						htmlComment();
						}
						break;
					default:
						throw new NoViableAltException(this);
					}
					setState(101);
					_errHandler.sync(this);
					_la = _input.LA(1);
					if (_la==SEA_WS || _la==HTML_TEXT) {
						{
						setState(100);
						htmlChardata();
						}
					}

					}
					} 
				}
				setState(107);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,15,_ctx);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class HtmlAttributeContext extends ParserRuleContext {
		public TerminalNode TAG_NAME() { return getToken(HTMLParser.TAG_NAME, 0); }
		public TerminalNode TAG_EQUALS() { return getToken(HTMLParser.TAG_EQUALS, 0); }
		public TerminalNode ATTVALUE_VALUE() { return getToken(HTMLParser.ATTVALUE_VALUE, 0); }
		public HtmlAttributeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_htmlAttribute; }
	}

	public final HtmlAttributeContext htmlAttribute() throws RecognitionException {
		HtmlAttributeContext _localctx = new HtmlAttributeContext(_ctx, getState());
		enterRule(_localctx, 10, RULE_htmlAttribute);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(108);
			match(TAG_NAME);
			setState(111);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==TAG_EQUALS) {
				{
				setState(109);
				match(TAG_EQUALS);
				setState(110);
				match(ATTVALUE_VALUE);
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class HtmlChardataContext extends ParserRuleContext {
		public TerminalNode HTML_TEXT() { return getToken(HTMLParser.HTML_TEXT, 0); }
		public TerminalNode SEA_WS() { return getToken(HTMLParser.SEA_WS, 0); }
		public HtmlChardataContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_htmlChardata; }
	}

	public final HtmlChardataContext htmlChardata() throws RecognitionException {
		HtmlChardataContext _localctx = new HtmlChardataContext(_ctx, getState());
		enterRule(_localctx, 12, RULE_htmlChardata);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(113);
			_la = _input.LA(1);
			if ( !(_la==SEA_WS || _la==HTML_TEXT) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class HtmlMiscContext extends ParserRuleContext {
		public HtmlCommentContext htmlComment() {
			return getRuleContext(HtmlCommentContext.class,0);
		}
		public TerminalNode SEA_WS() { return getToken(HTMLParser.SEA_WS, 0); }
		public HtmlMiscContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_htmlMisc; }
	}

	public final HtmlMiscContext htmlMisc() throws RecognitionException {
		HtmlMiscContext _localctx = new HtmlMiscContext(_ctx, getState());
		enterRule(_localctx, 14, RULE_htmlMisc);
		try {
			setState(117);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case HTML_COMMENT:
			case HTML_CONDITIONAL_COMMENT:
				enterOuterAlt(_localctx, 1);
				{
				setState(115);
				htmlComment();
				}
				break;
			case SEA_WS:
				enterOuterAlt(_localctx, 2);
				{
				setState(116);
				match(SEA_WS);
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class HtmlCommentContext extends ParserRuleContext {
		public TerminalNode HTML_COMMENT() { return getToken(HTMLParser.HTML_COMMENT, 0); }
		public TerminalNode HTML_CONDITIONAL_COMMENT() { return getToken(HTMLParser.HTML_CONDITIONAL_COMMENT, 0); }
		public HtmlCommentContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_htmlComment; }
	}

	public final HtmlCommentContext htmlComment() throws RecognitionException {
		HtmlCommentContext _localctx = new HtmlCommentContext(_ctx, getState());
		enterRule(_localctx, 16, RULE_htmlComment);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(119);
			_la = _input.LA(1);
			if ( !(_la==HTML_COMMENT || _la==HTML_CONDITIONAL_COMMENT) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ScriptContext extends ParserRuleContext {
		public TerminalNode SCRIPT_OPEN() { return getToken(HTMLParser.SCRIPT_OPEN, 0); }
		public TerminalNode SCRIPT_BODY() { return getToken(HTMLParser.SCRIPT_BODY, 0); }
		public TerminalNode SCRIPT_SHORT_BODY() { return getToken(HTMLParser.SCRIPT_SHORT_BODY, 0); }
		public ScriptContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_script; }
	}

	public final ScriptContext script() throws RecognitionException {
		ScriptContext _localctx = new ScriptContext(_ctx, getState());
		enterRule(_localctx, 18, RULE_script);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(121);
			match(SCRIPT_OPEN);
			setState(122);
			_la = _input.LA(1);
			if ( !(_la==SCRIPT_BODY || _la==SCRIPT_SHORT_BODY) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class StyleContext extends ParserRuleContext {
		public TerminalNode STYLE_OPEN() { return getToken(HTMLParser.STYLE_OPEN, 0); }
		public TerminalNode STYLE_BODY() { return getToken(HTMLParser.STYLE_BODY, 0); }
		public TerminalNode STYLE_SHORT_BODY() { return getToken(HTMLParser.STYLE_SHORT_BODY, 0); }
		public StyleContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_style; }
	}

	public final StyleContext style() throws RecognitionException {
		StyleContext _localctx = new StyleContext(_ctx, getState());
		enterRule(_localctx, 20, RULE_style);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(124);
			match(STYLE_OPEN);
			setState(125);
			_la = _input.LA(1);
			if ( !(_la==STYLE_BODY || _la==STYLE_SHORT_BODY) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static final String _serializedATN =
		"\u0004\u0001\u0017\u0080\u0002\u0000\u0007\u0000\u0002\u0001\u0007\u0001"+
		"\u0002\u0002\u0007\u0002\u0002\u0003\u0007\u0003\u0002\u0004\u0007\u0004"+
		"\u0002\u0005\u0007\u0005\u0002\u0006\u0007\u0006\u0002\u0007\u0007\u0007"+
		"\u0002\b\u0007\b\u0002\t\u0007\t\u0002\n\u0007\n\u0001\u0000\u0005\u0000"+
		"\u0018\b\u0000\n\u0000\f\u0000\u001b\t\u0000\u0001\u0000\u0003\u0000\u001e"+
		"\b\u0000\u0001\u0000\u0005\u0000!\b\u0000\n\u0000\f\u0000$\t\u0000\u0001"+
		"\u0000\u0003\u0000\'\b\u0000\u0001\u0000\u0005\u0000*\b\u0000\n\u0000"+
		"\f\u0000-\t\u0000\u0001\u0000\u0005\u00000\b\u0000\n\u0000\f\u00003\t"+
		"\u0000\u0001\u0001\u0001\u0001\u0001\u0002\u0005\u00028\b\u0002\n\u0002"+
		"\f\u0002;\t\u0002\u0001\u0002\u0001\u0002\u0005\u0002?\b\u0002\n\u0002"+
		"\f\u0002B\t\u0002\u0001\u0003\u0001\u0003\u0001\u0003\u0005\u0003G\b\u0003"+
		"\n\u0003\f\u0003J\t\u0003\u0001\u0003\u0001\u0003\u0001\u0003\u0001\u0003"+
		"\u0001\u0003\u0001\u0003\u0001\u0003\u0003\u0003S\b\u0003\u0001\u0003"+
		"\u0003\u0003V\b\u0003\u0001\u0003\u0001\u0003\u0001\u0003\u0003\u0003"+
		"[\b\u0003\u0001\u0004\u0003\u0004^\b\u0004\u0001\u0004\u0001\u0004\u0001"+
		"\u0004\u0003\u0004c\b\u0004\u0001\u0004\u0003\u0004f\b\u0004\u0005\u0004"+
		"h\b\u0004\n\u0004\f\u0004k\t\u0004\u0001\u0005\u0001\u0005\u0001\u0005"+
		"\u0003\u0005p\b\u0005\u0001\u0006\u0001\u0006\u0001\u0007\u0001\u0007"+
		"\u0003\u0007v\b\u0007\u0001\b\u0001\b\u0001\t\u0001\t\u0001\t\u0001\n"+
		"\u0001\n\u0001\n\u0001\n\u0000\u0000\u000b\u0000\u0002\u0004\u0006\b\n"+
		"\f\u000e\u0010\u0012\u0014\u0000\u0005\u0001\u0000\u0006\u0007\u0002\u0000"+
		"\u0007\u0007\u000b\u000b\u0001\u0000\u0001\u0002\u0001\u0000\u0012\u0013"+
		"\u0001\u0000\u0014\u0015\u0089\u0000\u0019\u0001\u0000\u0000\u0000\u0002"+
		"4\u0001\u0000\u0000\u0000\u00049\u0001\u0000\u0000\u0000\u0006Z\u0001"+
		"\u0000\u0000\u0000\b]\u0001\u0000\u0000\u0000\nl\u0001\u0000\u0000\u0000"+
		"\fq\u0001\u0000\u0000\u0000\u000eu\u0001\u0000\u0000\u0000\u0010w\u0001"+
		"\u0000\u0000\u0000\u0012y\u0001\u0000\u0000\u0000\u0014|\u0001\u0000\u0000"+
		"\u0000\u0016\u0018\u0003\u0002\u0001\u0000\u0017\u0016\u0001\u0000\u0000"+
		"\u0000\u0018\u001b\u0001\u0000\u0000\u0000\u0019\u0017\u0001\u0000\u0000"+
		"\u0000\u0019\u001a\u0001\u0000\u0000\u0000\u001a\u001d\u0001\u0000\u0000"+
		"\u0000\u001b\u0019\u0001\u0000\u0000\u0000\u001c\u001e\u0005\u0003\u0000"+
		"\u0000\u001d\u001c\u0001\u0000\u0000\u0000\u001d\u001e\u0001\u0000\u0000"+
		"\u0000\u001e\"\u0001\u0000\u0000\u0000\u001f!\u0003\u0002\u0001\u0000"+
		" \u001f\u0001\u0000\u0000\u0000!$\u0001\u0000\u0000\u0000\" \u0001\u0000"+
		"\u0000\u0000\"#\u0001\u0000\u0000\u0000#&\u0001\u0000\u0000\u0000$\"\u0001"+
		"\u0000\u0000\u0000%\'\u0005\u0005\u0000\u0000&%\u0001\u0000\u0000\u0000"+
		"&\'\u0001\u0000\u0000\u0000\'+\u0001\u0000\u0000\u0000(*\u0003\u0002\u0001"+
		"\u0000)(\u0001\u0000\u0000\u0000*-\u0001\u0000\u0000\u0000+)\u0001\u0000"+
		"\u0000\u0000+,\u0001\u0000\u0000\u0000,1\u0001\u0000\u0000\u0000-+\u0001"+
		"\u0000\u0000\u0000.0\u0003\u0004\u0002\u0000/.\u0001\u0000\u0000\u0000"+
		"03\u0001\u0000\u0000\u00001/\u0001\u0000\u0000\u000012\u0001\u0000\u0000"+
		"\u00002\u0001\u0001\u0000\u0000\u000031\u0001\u0000\u0000\u000045\u0007"+
		"\u0000\u0000\u00005\u0003\u0001\u0000\u0000\u000068\u0003\u000e\u0007"+
		"\u000076\u0001\u0000\u0000\u00008;\u0001\u0000\u0000\u000097\u0001\u0000"+
		"\u0000\u00009:\u0001\u0000\u0000\u0000:<\u0001\u0000\u0000\u0000;9\u0001"+
		"\u0000\u0000\u0000<@\u0003\u0006\u0003\u0000=?\u0003\u000e\u0007\u0000"+
		">=\u0001\u0000\u0000\u0000?B\u0001\u0000\u0000\u0000@>\u0001\u0000\u0000"+
		"\u0000@A\u0001\u0000\u0000\u0000A\u0005\u0001\u0000\u0000\u0000B@\u0001"+
		"\u0000\u0000\u0000CD\u0005\n\u0000\u0000DH\u0005\u0010\u0000\u0000EG\u0003"+
		"\n\u0005\u0000FE\u0001\u0000\u0000\u0000GJ\u0001\u0000\u0000\u0000HF\u0001"+
		"\u0000\u0000\u0000HI\u0001\u0000\u0000\u0000IU\u0001\u0000\u0000\u0000"+
		"JH\u0001\u0000\u0000\u0000KR\u0005\f\u0000\u0000LM\u0003\b\u0004\u0000"+
		"MN\u0005\n\u0000\u0000NO\u0005\u000e\u0000\u0000OP\u0005\u0010\u0000\u0000"+
		"PQ\u0005\f\u0000\u0000QS\u0001\u0000\u0000\u0000RL\u0001\u0000\u0000\u0000"+
		"RS\u0001\u0000\u0000\u0000SV\u0001\u0000\u0000\u0000TV\u0005\r\u0000\u0000"+
		"UK\u0001\u0000\u0000\u0000UT\u0001\u0000\u0000\u0000V[\u0001\u0000\u0000"+
		"\u0000W[\u0005\u0006\u0000\u0000X[\u0003\u0012\t\u0000Y[\u0003\u0014\n"+
		"\u0000ZC\u0001\u0000\u0000\u0000ZW\u0001\u0000\u0000\u0000ZX\u0001\u0000"+
		"\u0000\u0000ZY\u0001\u0000\u0000\u0000[\u0007\u0001\u0000\u0000\u0000"+
		"\\^\u0003\f\u0006\u0000]\\\u0001\u0000\u0000\u0000]^\u0001\u0000\u0000"+
		"\u0000^i\u0001\u0000\u0000\u0000_c\u0003\u0006\u0003\u0000`c\u0005\u0004"+
		"\u0000\u0000ac\u0003\u0010\b\u0000b_\u0001\u0000\u0000\u0000b`\u0001\u0000"+
		"\u0000\u0000ba\u0001\u0000\u0000\u0000ce\u0001\u0000\u0000\u0000df\u0003"+
		"\f\u0006\u0000ed\u0001\u0000\u0000\u0000ef\u0001\u0000\u0000\u0000fh\u0001"+
		"\u0000\u0000\u0000gb\u0001\u0000\u0000\u0000hk\u0001\u0000\u0000\u0000"+
		"ig\u0001\u0000\u0000\u0000ij\u0001\u0000\u0000\u0000j\t\u0001\u0000\u0000"+
		"\u0000ki\u0001\u0000\u0000\u0000lo\u0005\u0010\u0000\u0000mn\u0005\u000f"+
		"\u0000\u0000np\u0005\u0016\u0000\u0000om\u0001\u0000\u0000\u0000op\u0001"+
		"\u0000\u0000\u0000p\u000b\u0001\u0000\u0000\u0000qr\u0007\u0001\u0000"+
		"\u0000r\r\u0001\u0000\u0000\u0000sv\u0003\u0010\b\u0000tv\u0005\u0007"+
		"\u0000\u0000us\u0001\u0000\u0000\u0000ut\u0001\u0000\u0000\u0000v\u000f"+
		"\u0001\u0000\u0000\u0000wx\u0007\u0002\u0000\u0000x\u0011\u0001\u0000"+
		"\u0000\u0000yz\u0005\b\u0000\u0000z{\u0007\u0003\u0000\u0000{\u0013\u0001"+
		"\u0000\u0000\u0000|}\u0005\t\u0000\u0000}~\u0007\u0004\u0000\u0000~\u0015"+
		"\u0001\u0000\u0000\u0000\u0012\u0019\u001d\"&+19@HRUZ]beiou";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}