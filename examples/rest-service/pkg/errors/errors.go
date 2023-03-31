package errors

import "net/http"

// ServerError failed server request
// @gg:"http"
type ServerError struct{}

func (*ServerError) Error() string {
	return "server error"
}

func (*ServerError) Code() string {
	return "SERVER_ERROR"
}

func (*ServerError) StatusCode() int {
	return http.StatusInternalServerError
}
