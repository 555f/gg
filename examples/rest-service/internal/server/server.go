package server

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	errors1 "errors"
	helpers "github.com/555f/helpers"
	controller "github.com/f555/gg-examples/internal/usecase/controller"
	dto "github.com/f555/gg-examples/pkg/dto"
	errors "github.com/f555/gg-examples/pkg/errors"
	prometheus "github.com/prometheus/client_golang/prometheus"
	promauto "github.com/prometheus/client_golang/prometheus/promauto"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func encodeJSONResponse(ctx context.Context, w http.ResponseWriter, response any) {
	statusCode := 200
	var data bytes.Buffer
	if response != nil {
		if v, ok := response.(statusCoder); ok {
			statusCode = v.StatusCode()
		}
		if v, ok := response.(contenter); ok {
			w.Header().Set("Content-Type", v.ContentType())
		} else {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		}
		if v, ok := response.(dataer); ok {
			data = v.Data()
		} else {
			if err := json.NewEncoder(&data).Encode(response); err != nil {
				return
			}
		}
		if v, ok := response.(headerer); ok {
			for key, values := range v.Headers() {
				for _, val := range values {
					w.Header().Add(key, val)
				}
			}
		}
		if v, ok := response.(cookier); ok {
			for _, c := range v.Cookies() {
				http.SetCookie(w, &c)
			}
		}
	} else {
		statusCode = 204
	}
	w.WriteHeader(statusCode)
	if _, err := w.Write(data.Bytes()); err != nil {
		panic(err)
	}
}
func serverErrorEncoder(ctx context.Context, w http.ResponseWriter, err error) {
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

type Regex struct {
	Pattern string
	Params  []string
}
type route struct {
	method  string
	regex   *regexp.Regexp
	pattern string
	params  []string
	handler http.Handler
	before  []ServerBeforeFunc
	after   []ServerAfterFunc
}
type Server struct {
	staticRoutes map[string]*route
	regexRoutes  map[string][]*route
}
type RouteOption func(r *route)

func After(after ...ServerAfterFunc) RouteOption {
	return func(r *route) {
		r.after = append(r.after, after...)
	}
}
func Before(before ...ServerBeforeFunc) RouteOption {
	return func(r *route) {
		r.before = append(r.before, before...)
	}
}
func Middleware(middleware ...ServerMiddlewareFunc) RouteOption {
	return func(r *route) {
		r.handler = serverMiddlewareChain(middleware)(r.handler)
	}
}
func (s *Server) AddRoute(method string, pattern any, handler http.Handler, opts ...RouteOption) {
	r := &route{method: method, handler: handler}
	for _, opt := range opts {
		opt(r)
	}
	switch t := pattern.(type) {
	default:
		panic("pattern must be string or Regex type")
	case string:
		s.staticRoutes[method+t] = r
	case Regex:
		r.regex = regexp.MustCompile(t.Pattern)
		r.params = t.Params
		s.regexRoutes[method] = append(s.regexRoutes[method], r)
	}
}
func (s *Server) handleRoute(route *route, rw http.ResponseWriter, r *http.Request) {
	for _, before := range route.before {
		r = r.WithContext(before(r.Context(), r))
	}
	route.handler.ServeHTTP(rw, r)
	for _, after := range route.after {
		r = r.WithContext(after(r.Context(), rw))
	}
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Path
	if route, ok := s.staticRoutes[r.Method+requestPath]; ok {
		s.handleRoute(route, w, r)
		return
	}
	if routes, ok := s.regexRoutes[r.Method]; ok {
		for _, route := range routes {
			if !route.regex.MatchString(requestPath) {
				continue
			}
			matches := route.regex.FindStringSubmatch(requestPath)
			if len(matches[0]) != len(requestPath) {
				continue
			}
			if len(route.params) > 0 {
				values := r.URL.Query()
				for i, match := range matches[1:] {
					values.Add(route.params[i], match)
				}
				r.URL.RawQuery = values.Encode()
			}
			s.handleRoute(route, w, r)
			return
		}
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
func NewRESTServer(opts ...ServerOption) *Server {
	s := &Server{staticRoutes: make(map[string]*route), regexRoutes: make(map[string][]*route)}
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
func profileControllerCreateReqDecode(ctx context.Context, r *http.Request) (result any, err error) {
	param := new(profileControllerCreateRequest)
	contentType := r.Header.Get("content-type")
	parts := strings.Split(contentType, ";")
	if len(parts) == 0 {
		return nil, errors1.New("invalid content type")
	}
	contentType = parts[0]
	switch contentType {
	default:
		return nil, errors1.New("invalid content type")
	case "application/json":
		var body struct {
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Address   string `json:"address"`
			Zip       int    `json:"zip"`
		}
		var data []byte
		data, err = io.ReadAll(r.Body)
		if err != nil {
			return
		}
		err = json.Unmarshal(data, &body)
		if err != nil {
			return
		}
		param.firstName = body.FirstName
		param.lastName = body.LastName
		param.address = body.Address
		param.zip = body.Zip
	case "application/xml":
		var body struct {
			XMLName   xml.Name `xml:"profile"`
			FirstName string   `xml:"firstName"`
			LastName  string   `xml:"lastName"`
			Address   string   `xml:"address"`
			Zip       int      `xml:"zip"`
		}
		var data []byte
		data, err = io.ReadAll(r.Body)
		if err != nil {
			return
		}
		err = xml.Unmarshal(data, &body)
		if err != nil {
			return
		}
		param.firstName = body.FirstName
		param.lastName = body.LastName
		param.address = body.Address
		param.zip = body.Zip
	case "application/x-www-form-urlencoded":
		err = r.ParseForm()
		if err != nil {
			return
		}
		param.firstName = r.Form.Get("firstName")
		param.lastName = r.Form.Get("lastName")
		param.address = r.Form.Get("address")
		param.zip, err = helpers.ParseInt[int](r.Form.Get("zip"), 10, 64)
		if err != nil {
			return
		}
	case "multipart/form-data":
		err = r.ParseMultipartForm(int64(67108864))
		if err != nil {
			return
		}
		param.firstName = r.FormValue("firstName")
		param.lastName = r.FormValue("lastName")
		param.address = r.FormValue("address")
		param.zip, err = helpers.ParseInt[int](r.FormValue("zip"), 10, 64)
		if err != nil {
			return
		}
	}
	return param, nil
}
func profileControllerDownloadFileReqDecode(ctx context.Context, r *http.Request) (result any, err error) {
	param := new(profileControllerDownloadFileRequest)
	q := r.URL.Query()
	if s := q.Get("id"); s != "" {
		param.id = s
		if err != nil {
			return
		}
	}
	return param, nil
}
func profileControllerRemoveReqDecode(ctx context.Context, r *http.Request) (result any, err error) {
	param := new(profileControllerRemoveRequest)
	q := r.URL.Query()
	if s := q.Get("id"); s != "" {
		param.id = s
		if err != nil {
			return
		}
	}
	return param, nil
}

// ProfileController apply service
func ProfileController(svc controller.ProfileController, opts ...ProfileControllerServerOption) ServerOption {
	return func(s *Server) {
		profileController := &ProfileControllerOptions{}
		for _, opt := range opts {
			opt(profileController)
		}
		createOpts := []RouteOption{Before(append(profileController.before, profileController.create.before...)...), After(append(profileController.after, profileController.create.after...)...), Middleware(append(profileController.serverMiddleware, profileController.create.serverMiddleware...)...)}
		s.AddRoute("POST", "/profiles", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			reqData, err := profileControllerCreateReqDecode(r.Context(), r)
			if err != nil {
				serverErrorEncoder(r.Context(), rw, err)
				return
			}
			resp, err := middlewareChain(append(profileController.middleware, profileController.create.middleware...))(profileControllerCreateEndpoint(svc))(r.Context(), reqData)
			if err != nil {
				serverErrorEncoder(r.Context(), rw, err)
				return
			}
			encodeJSONResponse(r.Context(), rw, resp)
		}), createOpts...)
		downloadFileOpts := []RouteOption{Before(append(profileController.before, profileController.downloadFile.before...)...), After(append(profileController.after, profileController.downloadFile.after...)...), Middleware(append(profileController.serverMiddleware, profileController.downloadFile.serverMiddleware...)...)}
		s.AddRoute("GET", Regex{Pattern: "/profiles/([^/]+)/file", Params: []string{"id"}}, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			reqData, err := profileControllerDownloadFileReqDecode(r.Context(), r)
			if err != nil {
				serverErrorEncoder(r.Context(), rw, err)
				return
			}
			resp, err := middlewareChain(append(profileController.middleware, profileController.downloadFile.middleware...))(profileControllerDownloadFileEndpoint(svc))(r.Context(), reqData)
			if err != nil {
				serverErrorEncoder(r.Context(), rw, err)
				return
			}
			encodeJSONResponse(r.Context(), rw, resp)
		}), downloadFileOpts...)
		removeOpts := []RouteOption{Before(append(profileController.before, profileController.remove.before...)...), After(append(profileController.after, profileController.remove.after...)...), Middleware(append(profileController.serverMiddleware, profileController.remove.serverMiddleware...)...)}
		s.AddRoute("DELETE", Regex{Pattern: "/profiles/([^/]+)", Params: []string{"id"}}, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			reqData, err := profileControllerRemoveReqDecode(r.Context(), r)
			if err != nil {
				serverErrorEncoder(r.Context(), rw, err)
				return
			}
			resp, err := middlewareChain(append(profileController.middleware, profileController.remove.middleware...))(profileControllerRemoveEndpoint(svc))(r.Context(), reqData)
			if err != nil {
				serverErrorEncoder(r.Context(), rw, err)
				return
			}
			encodeJSONResponse(r.Context(), rw, resp)
		}), removeOpts...)
	}
}

func ProfileControllerMetricMiddleware(registry prometheus.Registerer, buckets []float64) ProfileControllerServerOption {
	return ProfileControllerApplyOptions(ProfileControllerCreateServerMiddlewareFunc(baseMetricMiddleware(registry, buckets, "ProfileController.Create")), ProfileControllerDownloadFileServerMiddlewareFunc(baseMetricMiddleware(registry, buckets, "ProfileController.DownloadFile")), ProfileControllerRemoveServerMiddlewareFunc(baseMetricMiddleware(registry, buckets, "ProfileController.Remove")))
}

type ProfileControllerServerOption func(*ProfileControllerOptions)
type ProfileControllerOptions struct {
	serverOptions
	create       serverOptions
	downloadFile serverOptions
	remove       serverOptions
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
func ProfileControllerDownloadFileBefore(before ...ServerBeforeFunc) ProfileControllerServerOption {
	return func(o *ProfileControllerOptions) {
		o.downloadFile.before = append(o.downloadFile.before, before...)
	}
}
func ProfileControllerDownloadFileAfter(after ...ServerAfterFunc) ProfileControllerServerOption {
	return func(o *ProfileControllerOptions) {
		o.downloadFile.after = append(o.downloadFile.after, after...)
	}
}
func ProfileControllerDownloadFileMiddleware(middleware ...EndpointMiddleware) ProfileControllerServerOption {
	return func(o *ProfileControllerOptions) {
		o.downloadFile.middleware = append(o.downloadFile.middleware, middleware...)
	}
}
func ProfileControllerDownloadFileServerMiddlewareFunc(serverMiddleware ...ServerMiddlewareFunc) ProfileControllerServerOption {
	return func(o *ProfileControllerOptions) {
		o.downloadFile.serverMiddleware = append(o.downloadFile.serverMiddleware, serverMiddleware...)
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
	zip       int
}

type profileControllerCreateResponse struct {
	Profile *dto.Profile `json:"profile"`
}

func profileControllerCreateEndpoint(s controller.ProfileController) func(ctx context.Context, request any) (any, error) {
	return func(ctx context.Context, request any) (any, error) {
		r := request.(*profileControllerCreateRequest)
		profile, err := s.Create(r.firstName, r.lastName, r.address, r.zip)
		return profileControllerCreateResponse{Profile: profile}, err
	}
}

type profileControllerDownloadFileRequest struct {
	id string
}

type profileControllerDownloadFileResponse struct {
	Data string `json:"data"`
}

func profileControllerDownloadFileEndpoint(s controller.ProfileController) func(ctx context.Context, request any) (any, error) {
	return func(ctx context.Context, request any) (any, error) {
		r := request.(*profileControllerDownloadFileRequest)
		data, err := s.DownloadFile(r.id)
		return profileControllerDownloadFileResponse{Data: data}, err
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
