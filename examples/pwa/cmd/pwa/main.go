package main

import (
	"log"
	"net/http"

	"github.com/555f/gg/examples/pwa/internal/component"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func main() {
	app.RouteWithRegexp("^/.*", &component.SideBarLayout{})

	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Name:        "Pass",
		Description: "Password manager",
		Styles: []string{
			"/web/style.css",
		},
	})
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal(err)
	}
}
