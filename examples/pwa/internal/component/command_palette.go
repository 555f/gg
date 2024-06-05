package component

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// @gg:"pwa"
// @pwa-view:"./views/command_palette.html"
type CommandPalette struct {
	app.Compo
}

func (c *CommandPalette) Props(props ...any) {
	fmt.Println(props...)
}

func (c *CommandPalette) OnMount(ctx app.Context) {
	fmt.Println("component mounted")
}
