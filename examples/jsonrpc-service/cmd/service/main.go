package main

import (
	"net/http"

	"github.com/f555/gg-examples/internal/interface/controller"

	"github.com/f555/gg-examples/internal/transport"
)

func main() {
	t := transport.NewJSONRPCServer(
		transport.ProfileController(
			new(controller.ProfileController),
		),
	)
	_ = http.ListenAndServe(":8080", t)
}
