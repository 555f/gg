// Code generated by GG version dev. DO NOT EDIT.

//go:build !gg
// +build !gg

package server

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	controller "github.com/555f/gg/examples/rest-service/internal/usecase/controller"
	dto "github.com/555f/gg/examples/rest-service/pkg/dto"
	errors "github.com/555f/gg/examples/rest-service/pkg/errors"
	gostrings "github.com/555f/go-strings"
	v4 "github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strings"
)

func echoDefaultErrorEncoder(ctx v4.Context, err error) {
	var statusCode int = http.StatusInternalServerError
	if e, ok := err.(interface {
		StatusCode() int
	}); ok {
		statusCode = e.StatusCode()
	}
	if headerer, ok := err.(interface {
		Headers() http.Header
	}); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				ctx.Response().Header().Add(k, v)
			}
		}
	}
	errorWrapper := errors.ErrorWrapper{}
	if e, ok := err.(interface {
		Data() interface{}
	}); ok {
		errorWrapper.Data = e.Data()
	}
	if e, ok := err.(interface {
		Error() string
	}); ok {
		errorWrapper.ErrorText = e.Error()
	}
	if e, ok := err.(interface {
		Code() string
	}); ok {
		errorWrapper.Code = e.Code()
	}

	err = ctx.JSON(statusCode, err)
	if err != nil {
		ctx.Response().Header().Add("content-type", "text/plain")
		ctx.Response().WriteHeader(500)
		ctx.Response().Write([]byte(err.Error()))
	}
}
func encodeBody(rw http.ResponseWriter, data any) {}

type contentTypeInvalidError struct{}

func (*contentTypeInvalidError) Error() string {
	return "Unsupported Media Type"
}
func (*contentTypeInvalidError) StatusCode() int {
	return 415
}
func decodeBody(r *http.Request, data any) {
	var bodyData = make([]byte, 0, 10485760)
	buf := bytes.NewBuffer(bodyData)
	written, err := io.Copy(buf, ctx.Request().Body)
	switch r.Header.Get("Content-Type") {
	case "application/xml":
		err = xml.Unmarshal(bodyData[:written], data)
	}
}

type ProfileControllerOption func(*ProfileControllerOptions)
type ProfileControllerOptions struct {
	errorEncoder           func(ctx v4.Context, err error)
	middleware             []v4.MiddlewareFunc
	middlewareCreate       []v4.MiddlewareFunc
	middlewareDownloadFile []v4.MiddlewareFunc
	middlewareRemove       []v4.MiddlewareFunc
}

func ProfileControllerMiddleware(middleware ...v4.MiddlewareFunc) ProfileControllerOption {
	return func(o *ProfileControllerOptions) {
		o.middleware = append(o.middleware, middleware...)
	}
}
func ProfileControllerWithErrorEncoder(errorEncoder func(ctx v4.Context, err error)) ProfileControllerOption {
	return func(o *ProfileControllerOptions) {
		o.errorEncoder = errorEncoder
	}
}
func ProfileControllerCreateMiddleware(middleware ...v4.MiddlewareFunc) ProfileControllerOption {
	return func(o *ProfileControllerOptions) {
		o.middlewareCreate = append(o.middlewareCreate, middleware...)
	}
}
func ProfileControllerDownloadFileMiddleware(middleware ...v4.MiddlewareFunc) ProfileControllerOption {
	return func(o *ProfileControllerOptions) {
		o.middlewareDownloadFile = append(o.middlewareDownloadFile, middleware...)
	}
}
func ProfileControllerRemoveMiddleware(middleware ...v4.MiddlewareFunc) ProfileControllerOption {
	return func(o *ProfileControllerOptions) {
		o.middlewareRemove = append(o.middlewareRemove, middleware...)
	}
}
func SetupRoutesProfileController(svc controller.ProfileController, e *v4.Echo, opts ...ProfileControllerOption) {
	o := &ProfileControllerOptions{errorEncoder: echoDefaultErrorEncoder}
	for _, opt := range opts {
		opt(o)
	}
	e.Add("POST", "/profiles", func(ctx v4.Context) (_ error) {
		var err error
		var req struct {
			XMLName   xml.Name `xml:"profile"`
			FirstName string   `json:"firstName"`
			LastName  string   `json:"lastName"`
			Address   string   `json:"address"`
			Zip       int      `json:"zip"`
		}
		contentTypeHeaderParam := ctx.Request().Header.Get("content-type")
		parts := strings.Split(contentTypeHeaderParam, ";")
		if len(parts) > 0 {
			contentTypeHeaderParam = parts[0]
		}
		var bodyData = make([]byte, 0, 10485760)
		buf := bytes.NewBuffer(bodyData)
		written, err := io.Copy(buf, ctx.Request().Body)
		if err != nil {
			o.errorEncoder(ctx, err)
			return
		}
		switch contentTypeHeaderParam {
		default:
			o.errorEncoder(ctx, &contentTypeInvalidError{})
			return
		case "application/json":
			err = json.Unmarshal(bodyData[:written], &req)
			if err != nil {
				o.errorEncoder(ctx, err)
				return
			}
		case "application/xml":
			err = xml.Unmarshal(bodyData[:written], &req)
			if err != nil {
				o.errorEncoder(ctx, err)
				return
			}
		case "application/x-www-form-urlencoded":
			f, err := ctx.FormParams()
			if err != nil {
				o.errorEncoder(ctx, err)
				return
			}
			req.FirstName = f.Get("firstName")
			req.LastName = f.Get("lastName")
			req.Address = f.Get("address")
			req.Zip, err = gostrings.ParseInt[int](f.Get("zip"), 10, 64)
			if err != nil {
				o.errorEncoder(ctx, err)
				return
			}
		case "multipart/form-data":
			f, err := ctx.FormParams()
			req.FirstName = f.Get("firstName")
			req.LastName = f.Get("lastName")
			req.Address = f.Get("address")
			req.Zip, err = gostrings.ParseInt[int](f.Get("zip"), 10, 64)
			if err != nil {
				o.errorEncoder(ctx, err)
				return
			}
		}
		profile, err := svc.Create(req.FirstName, req.LastName, req.Address, req.Zip)
		if err != nil {
			o.errorEncoder(ctx, err)
			return
		}
		var resp struct {
			Profile *dto.Profile `json:"profile"`
		}
		resp.Profile = profile
		acceptHeaderParam := ctx.Request().Header.Get("accept")
		switch acceptHeaderParam {
		case "application/json":
			err = json.Marshal(result)
			if err != nil {
				o.errorEncoder(ctx, err)
				return
			}
		case "application/xml":
			err = xml.Marshal(result)
			if err != nil {
				o.errorEncoder(ctx, err)
				return
			}
		}
		return
	})
	e.Add("GET", "/profiles/:id/file", func(ctx v4.Context) (_ error) {
		var err error
		idPathParam := ctx.Param("id")
		var paramId string
		if idPathParam != "" {
			paramId = idPathParam
		}
		data, err := svc.DownloadFile(paramId)
		if err != nil {
			o.errorEncoder(ctx, err)
			return
		}
		var resp struct {
			Data string `json:"data"`
		}
		resp.Data = data
		acceptHeaderParam := ctx.Request().Header.Get("accept")
		switch acceptHeaderParam {
		case "application/json":
			err = json.Marshal(result)
			if err != nil {
				o.errorEncoder(ctx, err)
				return
			}
		}
		return
	})
	e.Add("DELETE", "/profiles/{id}", func(ctx v4.Context) (_ error) {
		var req struct {
			Id string `json:"id"`
		}
		contentTypeHeaderParam := ctx.Request().Header.Get("content-type")
		parts := strings.Split(contentTypeHeaderParam, ";")
		if len(parts) > 0 {
			contentTypeHeaderParam = parts[0]
		}
		var bodyData = make([]byte, 0, 10485760)
		buf := bytes.NewBuffer(bodyData)
		written, err := io.Copy(buf, ctx.Request().Body)
		if err != nil {
			o.errorEncoder(ctx, err)
			return
		}
		switch contentTypeHeaderParam {
		default:
			o.errorEncoder(ctx, &contentTypeInvalidError{})
			return
		case "application/json":
			err = json.Unmarshal(bodyData[:written], &req)
			if err != nil {
				o.errorEncoder(ctx, err)
				return
			}
		}
		err = svc.Remove(req.Id)
		if err != nil {
			o.errorEncoder(ctx, err)
			return
		}
		return
	})
}
