package component

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// @gg:"goapp"
// @goapp-template:"./views/command_palette.html"
type CommandPalette struct {
	app.Compo
}

func (c *CommandPalette) OnMount(ctx app.Context) {
	fmt.Println("component mounted")
}
