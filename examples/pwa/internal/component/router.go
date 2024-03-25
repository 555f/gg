package component

import (
	"net/url"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type RouterComponent interface {
	app.UI
}

// @gg:"pwa"
type Route struct {
	Path      string
	TestInt   int
	component RouterComponent
}

func (r *Route) Children(component RouterComponent) *Route {
	r.component = component
	return r
}

func (r *Route) Try(path string) (url.Values, bool) {
	p := make(url.Values)
	if r.Path == path {
		return p, true
	}
	var i, j int
	for i < len(path) {
		switch {
		case j >= len(r.Path):
			if r.Path != "/" && len(r.Path) > 0 && r.Path[len(r.Path)-1] == '/' {
				return p, true
			}
			return nil, false
		case r.Path[j] == ':':
			var name, val string
			var nextc byte
			name, nextc, j = match(r.Path, isAlnum, j+1)
			val, _, i = match(path, matchPart(nextc), i)
			escval, err := url.QueryUnescape(val)
			if err != nil {
				return nil, false
			}
			p.Add(name, escval)
		case path[i] == r.Path[j]:
			i++
			j++
		default:
			return nil, false
		}
	}
	if j != len(r.Path) {
		return nil, false
	}
	return p, true
}

func matchPart(b byte) func(byte) bool {
	return func(c byte) bool {
		return c != b && c != '/'
	}
}

func match(s string, f func(byte) bool, i int) (matched string, next byte, j int) {
	j = i
	for j < len(s) && f(s[j]) {
		j++
	}
	if j < len(s) {
		next = s[j]
	}
	return s[i:j], next, j
}

func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isAlnum(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}

// @gg:"pwa"
type RouterProvider struct {
	app.Compo

	currentURL *url.URL
	routes     []*Route
}

func (c *RouterProvider) OnNav(ctx app.Context) {
	c.currentURL = ctx.Page().URL()
}

func (c *RouterProvider) Children(elems ...any) *RouterProvider {
	for _, e := range elems {
		if r, ok := e.(*Route); ok {
			c.routes = append(c.routes, r)
		}
	}
	return c
}

func (c *RouterProvider) Render() app.UI {
	if c.currentURL != nil {
		for _, r := range c.routes {
			if _, ok := r.Try(c.currentURL.Path); ok {
				// r.component.Props(params, c.currentURL.Query())
				return r.component
			}
		}
	}
	return app.Div().Body(app.Text("Route not found"))
}
