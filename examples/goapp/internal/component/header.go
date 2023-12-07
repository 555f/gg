package component

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// @gg:"goapp"
// @goapp-template:"./header.html"
type Settings struct {
	app.Compo
}

func (c *Settings) OnMount(ctx app.Context) {
	fmt.Println("component mounted")
}
