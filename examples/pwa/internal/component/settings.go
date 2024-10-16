package component

import (
	"fmt"
	"net/url"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type item struct {
	children []string
}

// @gg:"pwa"
// @pwa-view:"./views/settings.html"
type Settings struct {
	app.Compo
	items       []item
	title, test string
	num         int

	itemMap map[string]string

	isOk     bool
	Children app.UI
	enable   bool
	Params   url.Values
}

func (c *Settings) Enable(enable bool) *Settings {
	c.enable = enable
	return c
}

func (c *Settings) OnMount(ctx app.Context) {
	fmt.Println("component mounted")
	c.test = c.Params.Get("test")
}
