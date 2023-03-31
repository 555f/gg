package server

import (
	"bytes"
	"context"
	"encoding/json"
	controller "github.com/f555/gg-examples/internal/usecase/controller"
	dto "github.com/f555/gg-examples/pkg/dto"
	prometheus "github.com/prometheus/client_golang/prometheus"
	promauto "github.com/prometheus/client_golang/prometheus/promauto"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

const jsonRPCParseError int = -32700
const jsonRPCInvalidRequestError int = -32600
const jsonRPCMethodNotFoundError int = -32601
const jsonRPCInvalidParamsError int = -32602
const jsonRPCInternalError int = -32603

type JSONRPCHandlerFunc func(ctx context.Context, rw http.ResponseWriter, r *http.Request, prams json.RawMessage) (any, error)
type jsonRPCRoute struct {
	handler JSONRPCHandlerFunc
	before  []ServerBeforeFunc
	after   []ServerAfterFunc
}
type Server struct {
	routes map[string]jsonRPCRoute
}
type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
type jsonRPCRequestData struct {
	requests []jsonRPCRequest
	isBatch  bool
}

func (r *jsonRPCRequestData) UnmarshalJSON(b []byte) error {
	if bytes.HasPrefix(b, []byte("[")) {
		r.isBatch = true
		return json.Unmarshal(b, &r.requests)
	}
	var req jsonRPCRequest
	if err := json.Unmarshal(b, &req); err != nil {
		return err
	}
	r.requests = append(r.requests, req)
	return nil
}

type jsonRPCRequest struct {
	ID      any             `json:"id"`
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}
type jsonRPCResponse struct {
	ID      any             `json:"id"`
	Version string          `json:"jsonrpc"`
	Error   *jsonRPCError   `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}

func (s *Server) addRoute(method string, handler JSONRPCHandlerFunc, before []ServerBeforeFunc, after []ServerAfterFunc) {
	s.routes[method] = jsonRPCRoute{handler: handler, before: before, after: after}
}
func (s *Server) makeErrorResponse(id any, code int, message string) jsonRPCResponse {
	return jsonRPCResponse{ID: id, Version: "2.0", Error: &jsonRPCError{Code: code, Message: message}}
}
func (s *Server) handleRoute(route jsonRPCRoute, ctx context.Context, w http.ResponseWriter, r *http.Request, prams json.RawMessage) (resp any, err error) {
	for _, before := range route.before {
		ctx = before(ctx, r)
	}
	resp, err = route.handler(ctx, w, r, prams)
	if err != nil {
		return nil, err
	}
	for _, after := range route.after {
		ctx = after(ctx, w)
	}
	return
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestData jsonRPCRequestData
	var responses []jsonRPCResponse
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		responses = append(responses, s.makeErrorResponse(nil, jsonRPCParseError, err.Error()))
	} else {
		for _, req := range requestData.requests {
			if route, ok := s.routes[req.Method]; ok {
				resp, err := s.handleRoute(route, ctx, w, r, req.Params)
				if err != nil {
					responses = append(responses, s.makeErrorResponse(req.ID, jsonRPCInternalError, err.Error()))
					continue
				}
				result, err := json.Marshal(resp)
				if err != nil {
					responses = append(responses, s.makeErrorResponse(req.ID, jsonRPCInternalError, err.Error()))
					continue
				}
				responses = append(responses, jsonRPCResponse{ID: req.ID, Version: "2.0", Result: result})
			}
		}
	}
	var data any
	if requestData.isBatch {
		data = responses
	} else {
		data = responses[0]
	}
	_ = json.NewEncoder(w).Encode(data)
}
func NewJSONRPCServer(opts ...ServerOption) *Server {
	s := &Server{routes: make(map[string]jsonRPCRoute)}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type statusCoder interface {
	StatusCode() int
}
type dataer interface {
	Data() bytes.Buffer
}
type contenter interface {
	ContentType() string
}
type headerer interface {
	Headers() http.Header
}
type cookier interface {
	Cookies() []http.Cookie
}
type serverOptions struct {
	before           []ServerBeforeFunc
	after            []ServerAfterFunc
	serverMiddleware []ServerMiddlewareFunc
	middleware       []EndpointMiddleware
}
type ServerErrorEncoder func(ctx context.Context, err error, w http.ResponseWriter)
type ServerBeforeFunc func(context.Context, *http.Request) context.Context
type ServerAfterFunc func(context.Context, http.ResponseWriter) context.Context
type ServerMiddlewareFunc func(http.Handler) http.Handler
type ServerOption func(*Server)
type Endpoint = func(ctx context.Context, request interface{}) (response interface{}, err error)
type EndpointMiddleware = func(Endpoint) Endpoint

func middlewareChain(middlewares []EndpointMiddleware) EndpointMiddleware {
	return func(next Endpoint) Endpoint {
		if len(middlewares) == 0 {
			return next
		}
		outer := middlewares[0]
		others := middlewares[1:]
		for i := len(others) - 1; i >= 0; i-- {
			next = others[i](next)
		}
		return outer(next)
	}
}
func serverMiddlewareChain(serverMiddlewares []ServerMiddlewareFunc) ServerMiddlewareFunc {
	return func(next http.Handler) http.Handler {
		if len(serverMiddlewares) == 0 {
			return next
		}
		outer := serverMiddlewares[0]
		others := serverMiddlewares[1:]
		for i := len(others) - 1; i >= 0; i-- {
			next = others[i](next)
		}
		return outer(next)
	}
}
func baseMetricMiddleware(registry prometheus.Registerer, buckets []float64, handlerName string) func(http.Handler) http.Handler {
	reg := prometheus.WrapRegistererWith(prometheus.Labels{"handler": handlerName}, registry)
	requestsTotal := promauto.With(reg).NewCounterVec(prometheus.CounterOpts{Name: "http_requests_total", Help: "Tracks the number of HTTP requests."}, []string{"method", "code"})
	requestDuration := promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "Tracks the latencies for HTTP requests."}, []string{"method", "code"})
	requestSize := promauto.With(reg).NewSummaryVec(prometheus.SummaryOpts{Name: "http_request_size_bytes", Help: "Tracks the size of HTTP requests."}, []string{"method", "code"})
	responseSize := promauto.With(reg).NewSummaryVec(prometheus.SummaryOpts{Name: "http_response_size_bytes", Help: "Tracks the size of HTTP responses."}, []string{"method", "code"})
	return func(next http.Handler) http.Handler {
		return promhttp.InstrumentHandlerCounter(requestsTotal, promhttp.InstrumentHandlerDuration(requestDuration, promhttp.InstrumentHandlerRequestSize(requestSize, promhttp.InstrumentHandlerRequestSize(responseSize, next))))
	}
}
func profileControllerCreateReqDecode(ctx context.Context, r *http.Request, params json.RawMessage) (result any, err error) {
	body := new(profileControllerCreateRequest)
	err = json.Unmarshal(params, &body)
	if err != nil {
		return
	}
	return body, nil
}
func profileControllerRemoveReqDecode(ctx context.Context, r *http.Request, params json.RawMessage) (result any, err error) {
	body := new(profileControllerRemoveRequest)
	err = json.Unmarshal(params, &body)
	if err != nil {
		return
	}
	return body, nil
}

// ProfileController apply service
func ProfileController(svc controller.ProfileController, opts ...ProfileControllerServerOption) ServerOption {
	return func(s *Server) {
		profileController := &ProfileControllerOptions{}
		for _, opt := range opts {
			opt(profileController)
		}
		s.addRoute("profile.create", func(ctx context.Context, w http.ResponseWriter, r *http.Request, params json.RawMessage) (any, error) {
			reqData, err := profileControllerCreateReqDecode(ctx, r, params)
			if err != nil {
				return nil, err
			}
			resp, err := middlewareChain(append(profileController.middleware, profileController.create.middleware...))(profileControllerCreateEndpoint(svc))(ctx, reqData)
			if err != nil {
				return nil, err
			}
			return resp, nil
		}, append(profileController.before, profileController.create.before...), append(profileController.after, profileController.create.after...))
		s.addRoute("profile.delete", func(ctx context.Context, w http.ResponseWriter, r *http.Request, params json.RawMessage) (any, error) {
			reqData, err := profileControllerRemoveReqDecode(ctx, r, params)
			if err != nil {
				return nil, err
			}
			resp, err := middlewareChain(append(profileController.middleware, profileController.remove.middleware...))(profileControllerRemoveEndpoint(svc))(ctx, reqData)
			if err != nil {
				return nil, err
			}
			return resp, nil
		}, append(profileController.before, profileController.remove.before...), append(profileController.after, profileController.remove.after...))
	}
}

type ProfileControllerServerOption func(*ProfileControllerOptions)
type ProfileControllerOptions struct {
	serverOptions
	create serverOptions
	remove serverOptions
}

func ProfileControllerApplyOptions(options ...ProfileControllerServerOption) ProfileControllerServerOption {
	return func(o *ProfileControllerOptions) {
		for _, opt := range options {
			opt(o)
		}
	}
}
func ProfileControllerBefore(before ...ServerBeforeFunc) ProfileControllerServerOption {
	return func(o *ProfileControllerOptions) {
		o.before = append(o.before, before...)
	}
}
func ProfileControllerAfter(after ...ServerAfterFunc) ProfileControllerServerOption {
	return func(o *ProfileControllerOptions) {
		o.after = append(o.after, after...)
	}
}
func ProfileControllerMiddleware(middleware ...EndpointMiddleware) ProfileControllerServerOption {
	return func(o *ProfileControllerOptions) {
		o.middleware = append(o.middleware, middleware...)
	}
}
func ProfileControllerServerMiddleware(serverMiddleware ...ServerMiddlewareFunc) ProfileControllerServerOption {
	return func(o *ProfileControllerOptions) {
		o.serverMiddleware = append(o.serverMiddleware, serverMiddleware...)
	}
}
func ProfileControllerCreateBefore(before ...ServerBeforeFunc) ProfileControllerServerOption {
	return func(o *ProfileControllerOptions) {
		o.create.before = append(o.create.before, before...)
	}
}
func ProfileControllerCreateAfter(after ...ServerAfterFunc) ProfileControllerServerOption {
	return func(o *ProfileControllerOptions) {
		o.create.after = append(o.create.after, after...)
	}
}
func ProfileControllerCreateMiddleware(middleware ...EndpointMiddleware) ProfileControllerServerOption {
	return func(o *ProfileControllerOptions) {
		o.create.middleware = append(o.create.middleware, middleware...)
	}
}
func ProfileControllerCreateServerMiddlewareFunc(serverMiddleware ...ServerMiddlewareFunc) ProfileControllerServerOption {
	return func(o *ProfileControllerOptions) {
		o.create.serverMiddleware = append(o.create.serverMiddleware, serverMiddleware...)
	}
}
func ProfileControllerRemoveBefore(before ...ServerBeforeFunc) ProfileControllerServerOption {
	return func(o *ProfileControllerOptions) {
		o.remove.before = append(o.remove.before, before...)
	}
}
func ProfileControllerRemoveAfter(after ...ServerAfterFunc) ProfileControllerServerOption {
	return func(o *ProfileControllerOptions) {
		o.remove.after = append(o.remove.after, after...)
	}
}
func ProfileControllerRemoveMiddleware(middleware ...EndpointMiddleware) ProfileControllerServerOption {
	return func(o *ProfileControllerOptions) {
		o.remove.middleware = append(o.remove.middleware, middleware...)
	}
}
func ProfileControllerRemoveServerMiddlewareFunc(serverMiddleware ...ServerMiddlewareFunc) ProfileControllerServerOption {
	return func(o *ProfileControllerOptions) {
		o.remove.serverMiddleware = append(o.remove.serverMiddleware, serverMiddleware...)
	}
}

type profileControllerCreateRequest struct {
	firstName string
	lastName  string
	address   string
}

type profileControllerCreateResponse struct {
	Profile *dto.Profile `json:"profile"`
}

func profileControllerCreateEndpoint(s controller.ProfileController) func(ctx context.Context, request any) (any, error) {
	return func(ctx context.Context, request any) (any, error) {
		r := request.(*profileControllerCreateRequest)
		profile, err := s.Create(r.firstName, r.lastName, r.address)
		return profileControllerCreateResponse{Profile: profile}, err
	}
}

type profileControllerRemoveRequest struct {
	id string
}

func profileControllerRemoveEndpoint(s controller.ProfileController) func(ctx context.Context, request any) (any, error) {
	return func(ctx context.Context, request any) (any, error) {
		r := request.(*profileControllerRemoveRequest)
		err := s.Remove(r.id)
		return nil, err
	}
}
