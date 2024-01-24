package errors

// ErrorWrapper
// @gg:"http"
type ErrorWrapper struct {
	// @http-error-interface:"Data() interface{}"
	Data interface{} `json:"data,omitempty"`
	// @http-error-interface:"Error() string"
	ErrorText string `json:"errorText"`
	// @http-error-interface:"Code() string"
	Code string `json:"code"`
}

// DefaultError
// @gg:"http"
type DefaultError struct {
	Data      interface{}
	ErrorText string
	Code      string
}

func (e *DefaultError) Error() string {
	return e.ErrorText
}
