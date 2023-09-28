package gen

import (
	"github.com/555f/gg/internal/plugin/jsonrpc/options"

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
	jsonrpcPkg    = "github.com/555f/jsonrpc"
)

type HandlerStrategyBuilderFactory func() HandlerStrategy

type Qualifier interface {
	Qual(pkgPath, name string) func(s *jen.Statement)
}

type HandlerFuncBuilder interface {
	Build() jen.Code
}

type ServerControllerBuilder interface {
	Endpoint(ep options.Endpoint) ServerEndpointBuilder
	BuildHandlers() ServerControllerBuilder
}

type ServerEndpointBuilder interface {
	BuildReqStruct() ServerEndpointBuilder
	BuildRespStruct() ServerEndpointBuilder
	BuildReqDec() ServerEndpointBuilder
	// BuildRespEnc() ServerEndpointBuilder
	Build()
}

type ServerBuilder interface {
	Build() jen.Code
	BuildTypes() ServerBuilder
	// Controller(iface options.Iface) ServerControllerBuilder
}

type ClientBuilder interface {
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

type ExampleBuilder interface {
	Build() []byte
}

type HandlerStrategy interface {
	ID() string
	// ReqType() (typ jen.Code)
	ReqArgName() string
	// RespType() (typ jen.Code)
	RespArgName() string
	LibType() (typ jen.Code)
	LibArgName() string
	// QueryParams() (typ jen.Code)
	// QueryParam(queryName string) (name string, typ jen.Code)
	// PathParam(pathName string) (name string, typ jen.Code)
	// HeaderParam(headerName string) (name string, typ jen.Code)
	// BodyPathParam() (typ jen.Code)
	// FormParam(formName string) (name string, typ jen.Code)
	// MultipartFormParam(formName string) (name string, typ jen.Code)
	// FormParams() (typ jen.Code)
	// MultipartFormParams(multipartMaxMemory int64) (typ jen.Code)
	MiddlewareType() jen.Code
	HandlerFunc(method string, endpoint, middlewares jen.Code, bodyFunc ...jen.Code) (typ jen.Code)
	// SetHeader(k, v jen.Code) (typ jen.Code)
	// UsePathParams() bool
	// WriteError(statusCode, data jen.Code) (typ jen.Code)
}
