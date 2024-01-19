package rest

import (
	"github.com/555f/gg/pkg/gen"
	"github.com/dave/jennifer/jen"
)

var _ HandlerStrategy = &HandlerStrategyChi{}

type HandlerStrategyChi struct{}

func (s *HandlerStrategyChi) ID() string {
	return "chi"
}

func (s *HandlerStrategyChi) QueryParams() (typ jen.Code) {
	typ = jen.Id("q").Op(":=").Id(s.ReqArgName()).Dot("URL").Dot("Query").Call()
	return
}

func (s *HandlerStrategyChi) QueryParam(queryName string) (name string, typ jen.Code) {
	name = normalizeVarName(queryName) + "QueryParam"
	typ = jen.Id(name).Op(":=").Id("q").Dot("Get").Call(jen.Lit(queryName))
	return
}

func (s *HandlerStrategyChi) PathParam(pathName string) (name string, typ jen.Code) {
	name = normalizeVarName(pathName) + "PathParam"
	typ = jen.Id(name).Op(":=").Qual(chiPkg, "URLParam").Call(jen.Id(s.ReqArgName()), jen.Lit(pathName))
	return
}

func (s *HandlerStrategyChi) HeaderParam(headerName string) (name string, typ jen.Code) {
	name = normalizeVarName(headerName) + "HeaderParam"
	typ = jen.Id(name).Op(":=").Id(s.ReqArgName()).Dot("Header").Dot("Get").Call(jen.Lit(headerName))
	return
}

func (s *HandlerStrategyChi) BodyPathParam() (typ jen.Code) {
	return jen.Id(s.ReqArgName()).Dot("Body")
}

func (s *HandlerStrategyChi) FormParam(formName string) (name string, typ jen.Code) {
	name = normalizeVarName(formName) + "FormParam"
	typ = jen.Id(name).Op(":=").Id(s.ReqArgName()).Dot("Form").Dot("Get").Call(jen.Lit(formName))
	return
}

func (s *HandlerStrategyChi) FormParams() (typ jen.Code, hasErr bool) {
	typ = jen.Custom(jen.Options{Multi: true},
		jen.Err().Op("=").Id(s.ReqArgName()).Dot("ParseForm").Call(),
		jen.Do(gen.CheckErr(
			jen.Return(),
		)),
	)
	return
}

func (s *HandlerStrategyChi) MultipartFormParam(formName string) (name string, typ jen.Code) {
	name = normalizeVarName(formName) + "MpFormParam"
	typ = jen.Id(name).Op(":=").Id(s.ReqArgName()).Dot("FormValue").Call(jen.Lit(formName))
	return
}

func (s *HandlerStrategyChi) MultipartFormParams(multipartMaxMemory int64) (typ jen.Code, hasErr bool) {
	// typ = jen.Custom(jen.Options{Multi: true},
	// 	jen.Err().Op("=").Id(s.ReqArgName()).Dot("ParseMultipartForm").Call(jen.Lit(multipartMaxMemory)),
	// 	jen.Do(gen.CheckErr(
	// 		jen.Return(),
	// 	)),
	// )

	typ = jen.Id(s.ReqArgName()).Dot("ParseMultipartForm").Call(jen.Lit(multipartMaxMemory))

	return
}

func (*HandlerStrategyChi) ReqType() jen.Code {
	return jen.Op("*").Qual(httpPkg, "Request")
}

func (*HandlerStrategyChi) RespType() jen.Code {
	return jen.Qual(httpPkg, "ResponseWriter")
}

func (*HandlerStrategyChi) LibType() jen.Code {
	return jen.Qual(chiPkg, "Router")
}

func (s *HandlerStrategyChi) HandlerFuncParams() (in, out []jen.Code) {
	return nil, nil
}

func (s *HandlerStrategyChi) HandlerFunc(method string, pattern string, handlerFunc func(g *jen.Group)) jen.Code {
	return nil
	// 	return jen.Id(s.LibArgName()).Dot("With").Call(jen.Add(middlewares).Op("...")).Dot("MethodFunc").Params(
	// 		jen.Lit(method),
	// 		jen.Lit(pattern),
	// 		jen.Qual(httpPkg, "HandlerFunc").Call(jen.Func().Params(jen.Id("w").Qual(httpPkg, "ResponseWriter"), jen.Id(s.ReqArgName()).Op("*").Qual(httpPkg, "Request")).BlockFunc(func(g *jen.Group) {
	// 			g.Add(bodyFunc...)
	// 			if result != nil {
	// 				g.Add(s.SetHeader(jen.Lit("content-type"), jen.Lit("application/json")))
	// 				g.List(jen.Id("data"), jen.Err()).Op(":=").Qual("encoding/json", "Marshal").Call(result)

	//			g.Do(gen.CheckErr(
	//				s.SetHeader(jen.Lit("content-type"), jen.Lit("text/plain")),
	//				jen.Id(s.RespArgName()).Dot("WriteHeader").Call(jen.Lit(500)),
	//				jen.Id(s.RespArgName()).Dot("Write").Call(jen.Index().Byte().Call(jen.Id("err").Dot("Error").Call())),
	//				jen.Return()),
	//			)
	//			g.If(
	//				jen.List(jen.Id("_"), jen.Err()).Op(":=").Id(s.RespArgName()).Dot("Write").Call(jen.Id("data")),
	//				jen.Err().Op("!=").Nil(),
	//			).Block(
	//				s.SetHeader(jen.Lit("content-type"), jen.Lit("text/plain")),
	//				jen.Id(s.RespArgName()).Dot("WriteHeader").Call(jen.Lit(500)),
	//				jen.Id(s.RespArgName()).Dot("Write").Call(jen.Index().Byte().Call(jen.Id("err").Dot("Error").Call())),
	//				jen.Return(),
	//			)
	//		}
	//	})),
	//
	// )
}

func (*HandlerStrategyChi) MiddlewareType() jen.Code {
	return jen.Func().Params(jen.Qual(httpPkg, "Handler")).Qual(httpPkg, "Handler")
}

func (s *HandlerStrategyChi) SetHeader(k jen.Code, v jen.Code) (typ jen.Code) {
	return jen.Id(s.RespArgName()).Dot("Header").Call().Dot("Set").Call(k, v)
}

func (s *HandlerStrategyChi) WriteError(statusCode, data jen.Code) (typ jen.Code) {
	typ = jen.Custom(jen.Options{Multi: true},
		s.SetHeader(jen.Lit("content-type"), jen.Lit("application/json")),
		jen.Id(s.RespArgName()).Dot("WriteHeader").Call(statusCode),
		jen.List(jen.Id("bytes"), jen.Id("err")).Op(":=").Qual(jsonPkg, "Marshal").Call(data),
		jen.Do(gen.CheckErr(
			jen.Id(s.RespArgName()).Dot("WriteHeader").Call(jen.Lit(500)),
			jen.Id(s.RespArgName()).Dot("Write").Call(jen.Index().Byte().Call(jen.Id("err").Dot("Error").Call())),
			jen.Return(),
		)),
		jen.Id(s.RespArgName()).Dot("Write").Call(jen.Id("bytes")),
	)
	return
}

func (s *HandlerStrategyChi) WriteBody(body jen.Code) {

}

func (*HandlerStrategyChi) UsePathParams() bool {
	return true
}

func (*HandlerStrategyChi) RespArgName() string {
	return "w"
}

func (*HandlerStrategyChi) ReqArgName() string {
	return "r"
}

func (*HandlerStrategyChi) LibArgName() string {
	return "r"
}

func NewHandlerStrategyChi() *HandlerStrategyChi {
	return &HandlerStrategyChi{}
}
