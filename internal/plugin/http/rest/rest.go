package rest

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/dave/jennifer/jen"
)

const (
	ctxPkg        = "context"
	httpPkg       = "net/http"
	jsonPkg       = "encoding/json"
	contextPkg    = "context"
	fmtPkg        = "fmt"
	ioPkg         = "io"
	netPkg        = "net"
	prometheusPkg = "github.com/prometheus/client_golang/prometheus"
	promautoPkg   = "github.com/prometheus/client_golang/prometheus/promauto"
	promhttpPkg   = "github.com/prometheus/client_golang/prometheus/promhttp"
	chiPkg        = "github.com/go-chi/chi/v5"
	jsonrpcPkg    = "github.com/555f/jsonrpc"
	echoPkg       = "github.com/labstack/echo/v4"
)

type HandlerStrategyBuilderFactory func() HandlerStrategy

type Qualifier interface {
	Qual(pkgPath, name string) func(s *jen.Statement)
}

type HandlerFuncBuilder interface {
	Build() jen.Code
}

type ServerControllerBuilder interface {
	Build() ServerControllerBuilder
}

type ServerBuilder interface {
	SetErrorWrapper(errorWrapper *options.ErrorWrapper) ServerBuilder
	Build() jen.Code
	Controller(iface options.Iface) ServerControllerBuilder
}

type ClientBuilder interface {
	SetErrorWrapper(errorWrapper *options.ErrorWrapper) ClientBuilder
	BuildTypes() ClientBuilder
	BuildConstruct(iface options.Iface) ClientBuilder
	BuildStruct(iface options.Iface) ClientBuilder
	Endpoint(iface options.Iface, ep options.Endpoint) ClientEndpointBuilder
	Build() jen.Code
}

type ClientEndpointBuilder interface {
	BuildReqStruct() ClientEndpointBuilder
	BuildSetters() ClientEndpointBuilder
	BuildReqMethod() ClientEndpointBuilder
	BuildMethod() ClientEndpointBuilder
	BuildExecuteMethod() ClientEndpointBuilder
}

type HandlerStrategy interface {
	ID() string
	ReqArgName() string
	RespType() (typ jen.Code)
	RespArgName() string
	LibType() (typ jen.Code)
	LibArgName() string
	Context() jen.Code
	QueryParams() (typ jen.Code)
	QueryParam(queryName string) (name string, typ jen.Code)
	PathParam(pathName string) (name string, typ jen.Code)
	HeaderParam(headerName string) (name string, typ jen.Code)
	BodyPathParam() (typ jen.Code)
	FormParam(formName string) (name string, typ jen.Code)
	MultipartFormParam(formName string) (name string, typ jen.Code)
	FormParams() (typ jen.Code, hasErr bool)
	MultipartFormParams(multipartMaxMemory int64) (typ jen.Code, hasErr bool)
	MiddlewareType() jen.Code
	HandlerFunc(method, pattern string, middlewares jen.Code, handlerFunc func(g *jen.Group)) (typ jen.Code)
	SetHeader(k, v jen.Code) (typ jen.Code)
	UsePathParams() bool
	WriteError(statusCode, data jen.Code) (typ jen.Code)
	WriteBody(data, contentType jen.Code, statusCode int) (typ jen.Code)
}
