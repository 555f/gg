package main

import (
	"flag"
	"log"
	"net/http"
)

var servePath = flag.String("dir", "./public", "path to serve")

type logHandler struct {
	handler http.Handler
}

func (lh *logHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.Printf("(example-server) received request %s\n", r.URL.RequestURI())
	lh.handler.ServeHTTP(rw, r)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	log.Printf("Starting HTTP server on :8081 serving path %q Ctrl + C to close and quit", *servePath)
	log.Fatal(http.ListenAndServe(":8081", &logHandler{
		handler: http.FileServer(http.Dir(*servePath))},
	))
}
