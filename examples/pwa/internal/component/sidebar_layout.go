package component

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// @gg:"pwa2"
// @pwa-view:"./views/sidebar_layout.html"
type SideBarLayout struct {
	app.Compo

	// Body         app.UI
	currentPath  string
	sideBarClose bool
}

func (c *SideBarLayout) classNames(classes C) string {
	return classNames(classes)
}

func (c *SideBarLayout) toggleSideBar(ctx app.Context, e app.Event) {
	ctx.Dispatch(func(ctx app.Context) {
		c.sideBarClose = !c.sideBarClose
	})
}

func (c *SideBarLayout) activePath(path string, activeClass, defaultCalss string) string {
	if c.currentPath == path {
		return activeClass
	}
	return defaultCalss
}

func (c *SideBarLayout) OnMount(ctx app.Context) {
	c.currentPath = ctx.Page().URL().Path
}
