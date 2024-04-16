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

	Enable, isOk bool
	Children     app.UI

	Params url.Values
}

func (c *Settings) OnMount(ctx app.Context) {
	fmt.Println("component mounted")
	c.test = c.Params.Get("test")
}
