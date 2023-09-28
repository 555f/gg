package gen

import (
	"github.com/dave/jennifer/jen"
)

var _ HandlerStrategy = &HandlerStrategyJSONRPC{}

type HandlerStrategyJSONRPC struct{}

// HandlerFunc implements HandlerStrategy.
func (s *HandlerStrategyJSONRPC) HandlerFunc(method string, endpoint jen.Code, middlewares jen.Code, bodyFunc ...jen.Code) (typ jen.Code) {
	return jen.Id(s.LibArgName()).Dot("Register").Params(
		jen.Lit(method),
		endpoint,
		jen.Func().Params(
			jen.Id("ctx").Qual(contextPkg, "Context"),
			jen.Id("r").Op("*").Qual(httpPkg, "Request"),
			jen.Id("params").Qual(jsonPkg, "RawMessage"),
		).Params(
			jen.Id("req").Any(),
			jen.Id("err").Error(),
		).BlockFunc(func(g *jen.Group) {
			g.Add(bodyFunc...)
			// 	if result != nil {
			// 		g.Id("err").Op("=").Id(s.RespArgName()).Dot("JSON").Call(jen.Lit(200), result)
			// 		g.Do(gen.CheckErr(
			// 			s.SetHeader(jen.Lit("content-type"), jen.Lit("text/plain")),
			// 			jen.Id(s.RespArgName()).Dot("Response").Call().Dot("WriteHeader").Call(jen.Lit(500)),
			// 			jen.Id(s.RespArgName()).Dot("Response").Call().Dot("Write").Call(jen.Index().Byte().Call(jen.Id("err").Dot("Error").Call())),
			// 			jen.Return(),
			// 		))
			// 	}
			// 	g.Return(jen.Nil())
		}),
		jen.Add(middlewares).Op("..."),
	)
}

// ID implements HandlerStrategy.
func (*HandlerStrategyJSONRPC) ID() string {
	return "default"
}

// LibArgName implements HandlerStrategy.
func (*HandlerStrategyJSONRPC) LibArgName() string {
	return "r"
}

// LibType implements HandlerStrategy.
func (*HandlerStrategyJSONRPC) LibType() (typ jen.Code) {
	return jen.Op("*").Qual(jsonrpcPkg, "Server")
}

// MiddlewareType implements HandlerStrategy.
func (*HandlerStrategyJSONRPC) MiddlewareType() jen.Code {
	return jen.Qual(jsonrpcPkg, "Option")
}

// ReqArgName implements HandlerStrategy.
func (*HandlerStrategyJSONRPC) ReqArgName() string {
	return "params"
}

// RespArgName implements HandlerStrategy.
func (*HandlerStrategyJSONRPC) RespArgName() string {
	return "r"
}

func NewHandlerStrategyJSONRPC() *HandlerStrategyJSONRPC {
	return &HandlerStrategyJSONRPC{}
}
