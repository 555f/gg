package server

import (
	"bytes"
	"context"
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

type contentTypeInvalidError struct{}

func (*contentTypeInvalidError) Error() string {
	return "content type invalid"
}
func (*contentTypeInvalidError) StatusCode() int {
	return 400
}

type pathParams interface {
	Param(string) string
}
type pathParamsNoop struct{}

func (*pathParamsNoop) Param(string) string {
	return ""
}
func serverErrorEncoder(w http.ResponseWriter, err error) {
	var statusCode int = http.StatusInternalServerError
	h := w.Header()
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
				h.Add(k, v)
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
	data, jsonErr := json.Marshal(errorWrapper)
	if jsonErr != nil {
		_, _ = w.Write([]byte("unexpected marshal error"))
		return
	}
	h.Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write(data)
}
func httpHandler(ep func(ctx context.Context, request any) (any, error), reqDec func(pathParams pathParams, request *http.Request, params json.RawMessage) (result any, err error), respEnc func(result any) (any, error), pathParams pathParams) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var params any
		var result any
		var err error
		if reqDec != nil {
			var wb = make([]byte, 0, 10485760)
			buf := bytes.NewBuffer(wb)
			written, err := io.Copy(buf, r.Body)
			if err != nil {
				serverErrorEncoder(rw, err)
				return
			}
			params, err = reqDec(pathParams, r, wb[:written])
			if err != nil {
				serverErrorEncoder(rw, err)
				return
			}
		}
		result, err = ep(r.Context(), params)
		if err != nil {
			serverErrorEncoder(rw, err)
			return
		}
		if respEnc != nil {
			result, err = respEnc(result)
			if err != nil {
				serverErrorEncoder(rw, err)
				return
			}
			rw.Header().Set("Content-Type", "application/json")
			data, err := json.Marshal(result)
			if err != nil {
				serverErrorEncoder(rw, err)
				return
			}
			if _, err := rw.Write(data); err != nil {
				return
			}
		}
		statusCode := 200
		rw.WriteHeader(statusCode)
	}
}

type ProfileControllerOption func(*ProfileControllerOptions)
type ProfileControllerOptions struct {
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

type profileControllerCreateReq struct {
	XMLName   xml.Name `xml:"profile"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Address   string   `json:"address"`
	Zip       int      `json:"zip"`
}
type profileControllerCreateResp struct {
	Profile *dto.Profile `json:"profile"`
}

func profileControllerCreateEndpoint(svc controller.ProfileController) func(ctx context.Context, request any) (any, error) {
	return func(ctx context.Context, request any) (any, error) {
		r := request.(*profileControllerCreateReq)
		profile, err := svc.Create(r.FirstName, r.LastName, r.Address, r.Zip)
		if err != nil {
			return nil, err
		}
		return &profileControllerCreateResp{Profile: profile}, nil
	}
}
func profileControllerCreateRespEnc(result any) (any, error) {
	var wrapResult struct {
		Profile *dto.Profile `json:"profile"`
	}
	wrapResult.Profile = result.(*profileControllerCreateResp).Profile
	result = wrapResult
	return result, nil
}
func profileControllerCreateReqDec(pathParams pathParams, r *http.Request, params json.RawMessage) (result any, err error) {
	var param profileControllerCreateReq
	contentType := r.Header.Get("content-type")
	parts := strings.Split(contentType, ";")
	if len(parts) == 0 {
		return nil, err
	}
	contentType = parts[0]
	switch contentType {
	default:
		return nil, &contentTypeInvalidError{}
	case "application/json":
		err = json.Unmarshal(params, &param)
		if err != nil {
			return nil, err
		}
	case "application/xml":
		err = xml.Unmarshal(params, &param)
		if err != nil {
			return nil, err
		}
	case "application/x-www-form-urlencoded":
		err = r.ParseForm()
		if err != nil {
			return nil, err
		}
		param.FirstName = r.Form.Get("firstName")
		param.LastName = r.Form.Get("lastName")
		param.Address = r.Form.Get("address")
		param.Zip, err = gostrings.ParseInt[int](r.Form.Get("zip"), 10, 64)
		if err != nil {
			return nil, err
		}
	case "multipart/form-data":
		err = r.ParseMultipartForm(int64(67108864))
		if err != nil {
			return nil, err
		}
		param.FirstName = r.FormValue("firstName")
		param.LastName = r.FormValue("lastName")
		param.Address = r.FormValue("address")
		param.Zip, err = gostrings.ParseInt[int](r.FormValue("zip"), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return &param, nil
}

type profileControllerDownloadFileReq struct {
	Id string `json:"-"`
}
type profileControllerDownloadFileResp struct {
	Data string `json:"data"`
}

func profileControllerDownloadFileEndpoint(svc controller.ProfileController) func(ctx context.Context, request any) (any, error) {
	return func(ctx context.Context, request any) (any, error) {
		r := request.(*profileControllerDownloadFileReq)
		data, err := svc.DownloadFile(r.Id)
		if err != nil {
			return nil, err
		}
		return &profileControllerDownloadFileResp{Data: data}, nil
	}
}
func profileControllerDownloadFileRespEnc(result any) (any, error) {
	var wrapResult struct {
		Data string `json:"data"`
	}
	wrapResult.Data = result.(*profileControllerDownloadFileResp).Data
	result = wrapResult
	return result, nil
}
func profileControllerDownloadFileReqDec(pathParams pathParams, r *http.Request, _ json.RawMessage) (result any, err error) {
	var param profileControllerDownloadFileReq
	if s := pathParams.Param("id"); s != "" {
		param.Id = s
		if err != nil {
			return
		}
	}
	return &param, nil
}

type profileControllerRemoveReq struct {
	Id string `json:"id"`
}

func profileControllerRemoveEndpoint(svc controller.ProfileController) func(ctx context.Context, request any) (any, error) {
	return func(ctx context.Context, request any) (any, error) {
		r := request.(*profileControllerRemoveReq)
		err := svc.Remove(r.Id)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}
func profileControllerRemoveReqDec(pathParams pathParams, r *http.Request, params json.RawMessage) (result any, err error) {
	var param profileControllerRemoveReq
	contentType := r.Header.Get("content-type")
	parts := strings.Split(contentType, ";")
	if len(parts) == 0 {
		return nil, err
	}
	contentType = parts[0]
	switch contentType {
	default:
		return nil, &contentTypeInvalidError{}
	case "application/json":
		err = json.Unmarshal(params, &param)
		if err != nil {
			return nil, err
		}
	}
	return &param, nil
}

// SetupRoutesProfileController route init for service
func SetupRoutesProfileController(svc controller.ProfileController, s *v4.Echo, opts ...ProfileControllerOption) {
	o := &ProfileControllerOptions{}
	for _, opt := range opts {
		opt(o)
	}
	s.Add("POST", "/profiles", func(ctx v4.Context) error {
		httpHandler(profileControllerCreateEndpoint(svc), profileControllerCreateReqDec, profileControllerCreateRespEnc, ctx).ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}, append(o.middleware, o.middlewareCreate...)...)
	s.Add("GET", "/profiles/:id/file", func(ctx v4.Context) error {
		httpHandler(profileControllerDownloadFileEndpoint(svc), profileControllerDownloadFileReqDec, profileControllerDownloadFileRespEnc, ctx).ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}, append(o.middleware, o.middlewareDownloadFile...)...)
	s.Add("DELETE", "/profiles/{id}", func(ctx v4.Context) error {
		httpHandler(profileControllerRemoveEndpoint(svc), profileControllerRemoveReqDec, nil, ctx).ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}, append(o.middleware, o.middlewareRemove...)...)
}
