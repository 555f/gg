package component

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type SVGElem struct {
	app.HTMLElem
}

func SVG() SVGElem {
	return SVGElem{HTMLElem: app.Elem("svg").XMLNS("http://www.w3.org/2000/svg")}
}

func (e SVGElem) Class(v ...string) SVGElem {
	e.HTMLElem.Class(v...)
	return e
}

func (e SVGElem) ViewBox(s string) SVGElem {
	e.Attr("view-box", s)
	return e
}

func (e SVGElem) Fill(s string) SVGElem {
	e.Attr("fill", s)
	return e
}

func (e SVGElem) StrokeWidth(s string) SVGElem {
	e.Attr("stroke-width", s)
	return e
}

func (e SVGElem) Stroke(s string) SVGElem {
	e.Attr("stroke", s)
	return e
}

type SVGPathElem struct {
	app.HTMLElemSelfClosing
}

func SVGPath() SVGPathElem {
	return SVGPathElem{HTMLElemSelfClosing: app.ElemSelfClosing("path").XMLNS("http://www.w3.org/2000/svg")}
}

func (e SVGPathElem) Class(v ...string) SVGPathElem {
	e.HTMLElemSelfClosing.Class(v...)
	return e
}

func (e SVGPathElem) FillRule(s string) SVGPathElem {
	e.Attr("fill-rule", s)
	return e
}

func (e SVGPathElem) D(s string) SVGPathElem {
	e.Attr("d", s)
	return e
}

func (e SVGPathElem) ClipRule(s string) SVGPathElem {
	e.Attr("clip-rule", s)
	return e
}

func (e SVGPathElem) StrokeLinecap(s string) SVGPathElem {
	e.Attr("stroke-linecap", s)
	return e
}

func (e SVGPathElem) StrokeLinejoin(s string) SVGPathElem {
	e.Attr("stroke-linejoin", s)
	return e
}
