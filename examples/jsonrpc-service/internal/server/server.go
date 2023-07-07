package server

import (
	"bytes"
	"context"
	"encoding/json"
	controller "github.com/555f/gg/examples/jsonrpc-service/internal/usecase/controller"
	dto "github.com/555f/gg/examples/jsonrpc-service/pkg/dto"
	jsonrpc "github.com/555f/jsonrpc"
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
	w.WriteHeader(statusCode)
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

type profileControllerCreateReq struct {
	Token     string `json:"token"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address   string `json:"address"`
}
type profileControllerCreateResp struct {
	Profile *dto.Profile `json:"profile"`
}

func profileControllerCreateEndpoint(svc controller.ProfileController) func(ctx context.Context, request any) (any, error) {
	return func(ctx context.Context, request any) (any, error) {
		r := request.(*profileControllerCreateReq)
		profile, err := svc.Create(r.Token, r.FirstName, r.LastName, r.Address)
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
func SetupRoutesProfileController(svc controller.ProfileController, s *jsonrpc.Server) {
	s.Register("profile.create", profileControllerCreateEndpoint(svc), func(ctx context.Context, r *http.Request, params json.RawMessage) (result any, err error) {
		return profileControllerCreateReqDec(&pathParamsNoop{}, r, params)
	})
	s.Register("profile.delete", profileControllerRemoveEndpoint(svc), func(ctx context.Context, r *http.Request, params json.RawMessage) (result any, err error) {
		return profileControllerRemoveReqDec(&pathParamsNoop{}, r, params)
	})
}
