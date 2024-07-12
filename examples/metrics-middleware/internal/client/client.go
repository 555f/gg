package client

// @gg:"middleware"
// @gg:"metrics"
type FooClient interface {
	BarMethod(test string) (n int, err error)
}
