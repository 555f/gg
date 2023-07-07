package main

import (
	"net/http"

	"github.com/555f/gg/examples/jsonrpc-service/internal/interface/controller"
	"github.com/555f/gg/examples/jsonrpc-service/internal/server"
	"github.com/555f/jsonrpc"
)

func main() {
	s := jsonrpc.NewServer()
	server.SetupRoutesProfileController(new(controller.ProfileController), s)
	_ = http.ListenAndServe(":8080", s)
}
