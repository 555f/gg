package component

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type item struct {
	children []string
}

// @gg:"goapp"
// @goapp-template:"./views/header.html"
// @goapp-tagname:"settings"
type Settings struct {
	app.Compo
	items []item
	title string
}

func (c *Settings) OnMount(ctx app.Context) {
	fmt.Println("component mounted")
}
