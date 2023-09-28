package rest

import (
	"github.com/555f/gg/pkg/gen"
	"github.com/dave/jennifer/jen"
)

var _ HandlerStrategy = &HandlerStrategyEcho{}

type HandlerStrategyEcho struct {
}

func (s *HandlerStrategyEcho) ID() string {
	return "echo"
}

func (s *HandlerStrategyEcho) QueryParams() (typ jen.Code) {
	typ = jen.Id("q").Op(":=").Id(s.ReqArgName()).Dot("QueryParams").Call()
	return
}

func (s *HandlerStrategyEcho) QueryParam(queryName string) (name string, typ jen.Code) {
	name = normalizeVarName(queryName) + "QueryParam"
	typ = jen.Id(name).Op(":=").Id("q").Dot("Get").Call(jen.Lit(queryName))
	return
}

func (s *HandlerStrategyEcho) PathParam(pathName string) (name string, typ jen.Code) {
	name = normalizeVarName(pathName) + "PathParam"
	typ = jen.Id(name).Op(":=").Id(s.ReqArgName()).Dot("Param").Call(jen.Lit(pathName))
	return
}

func (s *HandlerStrategyEcho) HeaderParam(headerName string) (name string, typ jen.Code) {
	name = normalizeVarName(headerName) + "HeaderParam"
	typ = jen.Id(name).Op(":=").Id(s.ReqArgName()).Dot("Request").Call().Dot("Header").Dot("Get").Call(jen.Lit(headerName))
	return
}

func (s *HandlerStrategyEcho) BodyPathParam() (typ jen.Code) {
	return jen.Id(s.ReqArgName()).Dot("Request").Call().Dot("Body")
}

func (*HandlerStrategyEcho) FormParam(formName string) (name string, typ jen.Code) {
	name = normalizeVarName(formName) + "FormParam"
	typ = jen.Id(name).Op(":=").Id("f").Dot("Get").Call(jen.Lit(formName))
	return
}

func (s *HandlerStrategyEcho) FormParams() (typ jen.Code) {
	typ = jen.Custom(jen.Options{Multi: true},
		jen.List(jen.Id("f"), jen.Err()).Op(":=").Id(s.ReqArgName()).Dot("FormParams").Call(),
		jen.Do(gen.CheckErr(
			jen.Return(jen.Nil(), jen.Err()),
		)),
	)
	return
}

func (s *HandlerStrategyEcho) MultipartFormParam(formName string) (name string, typ jen.Code) {
	name = normalizeVarName(formName) + "MpFormParam"
	typ = jen.Id(name).Op(":=").Id("f").Dot("Get").Call(jen.Lit(formName))
	return
}

func (s *HandlerStrategyEcho) MultipartFormParams(multipartMaxMemory int64) (typ jen.Code) {
	typ = jen.Custom(jen.Options{Multi: true},
		jen.List(jen.Id("f"), jen.Err()).Op(":=").Id(s.ReqArgName()).Dot("FormParams").Call(),
		jen.Do(gen.CheckErr(
			jen.Return(jen.Nil(), jen.Err()),
		)),
	)
	return
}

func (*HandlerStrategyEcho) ReqType() jen.Code {
	return jen.Qual(echoPkg, "Context")
}

func (*HandlerStrategyEcho) RespType() jen.Code {
	return jen.Qual(echoPkg, "Context")
}

func (*HandlerStrategyEcho) MiddlewareType() jen.Code {
	return jen.Qual(echoPkg, "MiddlewareFunc")
}

func (*HandlerStrategyEcho) LibType() jen.Code {
	return jen.Op("*").Qual(echoPkg, "Echo")
}

func (s *HandlerStrategyEcho) HandlerFunc(method string, pattern string, result, middlewares jen.Code, bodyFunc ...jen.Code) jen.Code {

	return jen.Id(s.LibArgName()).Dot("Add").Params(
		jen.Lit(method),
		jen.Lit(pattern),
		jen.Func().Params(jen.Id(s.ReqArgName()).Qual(echoPkg, "Context")).Params(jen.Id("_").Error()).BlockFunc(func(g *jen.Group) {
			g.Add(bodyFunc...)
			if result != nil {
				g.Id("err").Op("=").Id(s.RespArgName()).Dot("JSON").Call(jen.Lit(200), result)
				g.Do(gen.CheckErr(
					s.SetHeader(jen.Lit("content-type"), jen.Lit("text/plain")),
					jen.Id(s.RespArgName()).Dot("Response").Call().Dot("WriteHeader").Call(jen.Lit(500)),
					jen.Id(s.RespArgName()).Dot("Response").Call().Dot("Write").Call(jen.Index().Byte().Call(jen.Id("err").Dot("Error").Call())),
					jen.Return(),
				))
			}
			g.Return(jen.Nil())
		}),
		jen.Add(middlewares).Op("..."),
	)
}

func (s *HandlerStrategyEcho) SetHeader(k jen.Code, v jen.Code) (typ jen.Code) {
	return jen.Id(s.RespArgName()).Dot("Response").Call().Dot("Header").Call().Dot("Add").Call(k, v)
}

func (s *HandlerStrategyEcho) WriteError(statusCode, data jen.Code) (typ jen.Code) {
	typ = jen.Custom(jen.Options{Multi: true},
		jen.Id("err").Op("=").Id(s.RespArgName()).Dot("JSON").Call(statusCode, data),
		jen.Do(gen.CheckErr(
			s.SetHeader(jen.Lit("content-type"), jen.Lit("text/plain")),
			jen.Id(s.RespArgName()).Dot("Response").Call().Dot("WriteHeader").Call(jen.Lit(500)),
			jen.Id(s.RespArgName()).Dot("Response").Call().Dot("Write").Call(jen.Index().Byte().Call(jen.Id("err").Dot("Error").Call())),
		)),
	)
	return
}

func (*HandlerStrategyEcho) RespArgName() string {
	return "ctx"
}

func (*HandlerStrategyEcho) ReqArgName() string {
	return "ctx"
}

func (*HandlerStrategyEcho) LibArgName() string {
	return "e"
}

func (*HandlerStrategyEcho) UsePathParams() bool {
	return true
}

func NewHandlerStrategyEcho() *HandlerStrategyEcho {
	return &HandlerStrategyEcho{}
}
