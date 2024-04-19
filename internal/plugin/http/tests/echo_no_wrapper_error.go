package tests

// @gg:"http"
// @http-type:"echo"
// @http-client
type EchoClientNoWrapperErrorController interface {
	// @http-method:"GET"
	// @http-path:"/foo/:a"
	Foo(a string) (err error)
}
