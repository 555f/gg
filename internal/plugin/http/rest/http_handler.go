package rest

import (
	"github.com/555f/gg/pkg/file"
	. "github.com/dave/jennifer/jen"
)

func GenHandler() func(f *file.GoFile) {
	return func(f *file.GoFile) {
		f.Func().Id("httpHandler").Params(
			Id("ep").Add(epFunc),
			Id("reqDec").Add(reqDecFunc),
			Id("respEnc").Add(respEncFunc),
		).Qual("net/http", "HandlerFunc").Block(
			Return(
				Func().Params(
					Id("rw").Qual("net/http", "ResponseWriter"),
					Id("r").Op("*").Qual("net/http", "Request"),
				).BlockFunc(func(g *Group) {
					g.Var().Id("wb").Op("=").Make(Index().Byte(), Lit(0), Lit(10485760)) // 10MB
					g.Id("buf").Op(":=").Qual("bytes", "NewBuffer").Call(Id("wb"))
					g.List(Id("written"), Id("err")).Op(":=").Qual("io", "Copy").Call(Id("buf"), Id("r").Dot("Body"))
					g.Do(serverErrorEncoder)
					g.List(Id("params"), Err()).Op(":=").Id("reqDec").Call(
						//Id("r").Dot("Context").Call(),
						Id("r"),
						Id("wb").Index(Op(":").Id("written")),
					)
					g.Do(serverErrorEncoder)
					g.List(Id("result"), Err()).Op(":=").Id("ep").Call(
						Id("r").Dot("Context").Call(),
						Id("params"),
					)
					g.Do(serverErrorEncoder)

					g.List(Id("result"), Err()).Op("=").Id("respEnc").Call(
						//Id("r").Dot("Context").Call(),
						Id("result"),
					)
					g.Do(serverErrorEncoder)

					g.Id("statusCode").Op(":=").Lit(200)

					g.Id("rw").Dot("WriteHeader").Call(Id("statusCode"))

					//if len(ep.BodyResults) > 0 {
					g.List(Id("data"), Err()).Op(":=").Qual("encoding/json", "Marshal").Call(Id("result"))
					g.Do(serverErrorEncoder)
					g.If(
						List(Id("_"), Err()).Op(":=").Id("rw").Dot("Write").Call(Id("data")),
						Err().Op("!=").Nil(),
					).Block(
						Return(),
					)
					//}

					//func profileControllerRemoveHTTPHandler(ep func(ctx context.Context, request any) (response any, err error)) http.HandlerFunc {
					//	return func(rw http.ResponseWriter, r *http.Request) {
					//	var wb = make([]byte, 0, 10485760)
					//	buf := bytes.NewBuffer(wb)
					//	written, err := io.Copy(buf, r.Body)
					//	if err != nil {
					//	serverErrorEncoder(r.Context(), rw, err)
					//	return
					//}
					//	request, err := profileControllerRemoveReqDec(r.Context(), r, wb[:written])
					//	if err != nil {
					//	serverErrorEncoder(r.Context(), rw, err)
					//	return
					//}
					//	response, err := ep(r.Context(), params)
					//	if err != nil {
					//	serverErrorEncoder(r.Context(), rw, err)
					//	return
					//}
					//	err = svc.Remove(paramId)
					//	if err != nil {
					//	serverErrorEncoder(r.Context(), rw, err)
					//	return
					//}
					//	statusCode := 204
					//	rw.WriteHeader(statusCode)
					//}
					//}

					//g.Var().Err().Error()
					//
					//if len(ep.Params) > 0 {
					//	for _, p := range ep.Params {
					//		g.Var().Id("param" + p.FldName).Add(types.Convert(p.Type, f.Import))
					//	}
					//
					//	if len(ep.PathParams) > 0 {
					//		for _, p := range ep.PathParams {
					//			g.Id("pathParams").Op(":=").Id("pathParamsFromContext").Call(Id("r").Dot("Context").Call())
					//
					//			g.If(Id("s").Op(":=").Id("pathParams").Dot("Get").Call(Lit(p.Name)), Id("s").Op("!=").Lit("")).Block(
					//				Add(gen.ParseValue(Id("s"), Id("param"+p.FldName), "=", p.Type, f.Import)),
					//				Do(gen.CheckErr(
					//					Id("serverErrorEncoder").Call(Id("r").Dot("Context").Call(), Id("rw"), Err()),
					//					Return(),
					//				)),
					//			)
					//		}
					//	}
					//
					//	if len(ep.QueryParams) > 0 {
					//		g.Id("q").Op(":=").Id("r").Dot("URL").Dot("Query").Call()
					//		for _, param := range ep.QueryParams {
					//			paramID := Id("param" + param.FldName)
					//			if param.Parent != nil {
					//				paramID = Id("param" + param.Parent.FldName).Dot(param.FldName)
					//			}
					//			g.If(Id("s").Op(":=").Id("q").Dot("Get").Call(Lit(param.Name)), Id("s").Op("!=").Lit("")).Block(
					//				Add(gen.ParseValue(Id("s"), paramID, "=", param.Type, f.Import)),
					//				Do(gen.CheckErr(
					//					Id("serverErrorEncoder").Call(Id("r").Dot("Context").Call(), Id("rw"), Err()),
					//					Return(),
					//				)),
					//			)
					//		}
					//	}
					//
					//	if len(ep.HeaderParams) > 0 {
					//		for _, p := range ep.HeaderParams {
					//			g.If(Id("s").Op(":=").Id("r").Dot("Header").Dot("Get").Call(Lit(p.Name)), Id("s").Op("!=").Lit("")).Block(
					//				Add(gen.ParseValue(Id("s"), Id("param"+p.FldName), "=", p.Type, f.Import)),
					//				Do(gen.CheckErr(
					//					Id("serverErrorEncoder").Call(Id("r").Dot("Context").Call(), Id("rw"), Err()),
					//					Return(),
					//				)),
					//			)
					//		}
					//	}
					//
					//	if len(ep.BodyParams) > 0 {
					//		switch ep.HTTPMethod {
					//		case "POST", "PUT", "PATCH", "DELETE":
					//			g.Id("contentType").Op(":=").Id("r").Dot("Header").Dot("Get").Call(Lit("content-type"))
					//
					//			g.Id("parts").Op(":=").Qual("strings", "Split").Call(Id("contentType"), Lit(";"))
					//			g.If(Len(Id("parts")).Op("==").Lit(0)).Block(
					//				Id("serverErrorEncoder").Call(Id("r").Dot("Context").Call(), Id("rw"), Qual("errors", "New").Call(Lit("invalid content type"))),
					//				Return(),
					//			)
					//			g.Id("contentType").Op("=").Id("parts").Index(Lit(0))
					//
					//			g.Switch(Id("contentType")).BlockFunc(func(g *Group) {
					//				g.Default().Block(
					//					Id("serverErrorEncoder").Call(Id("r").Dot("Context").Call(), Id("rw"), Qual("errors", "New").Call(Lit("invalid content type"))),
					//					Return(),
					//				)
					//
					//				for _, contentType := range ep.ContentTypes {
					//					switch contentType {
					//					case "xml":
					//						g.Case(Lit("application/xml")).BlockFunc(func(g *Group) {
					//							g.Var().Id("body").Id(strcase.ToLowerCamel(s.Name+ep.MethodName) + "Req")
					//							g.Var().Id("wb").Op("=").Make(Index().Byte(), Lit(0), Lit(10485760)) // 10MB
					//							g.Id("buf").Op(":=").Qual("bytes", "NewBuffer").Call(Id("wb"))
					//							g.List(Id("written"), Id("err")).Op(":=").Qual("io", "Copy").Call(Id("buf"), Id("r").Dot("Body"))
					//							g.Do(gen.CheckErr(
					//								Id("serverErrorEncoder").Call(Id("r").Dot("Context").Call(), Id("rw"), Err()),
					//								Return(),
					//							))
					//							g.Id("data").Op(":=").Id("wd").Index(Op(":").Id("written"))
					//
					//							g.Err().Op("=").Qual("encoding/xml", "Unmarshal").Call(Id("data"), Op("&").Id("body"))
					//							g.Do(gen.CheckErr(
					//								Id("serverErrorEncoder").Call(Id("r").Dot("Context").Call(), Id("rw"), Err()),
					//								Return(),
					//							))
					//							for _, p := range ep.BodyParams {
					//								g.Id("param" + p.FldName).Op("=").Id("body").Dot(p.FldName)
					//							}
					//						})
					//					case "json":
					//						g.Case(Lit("application/json")).BlockFunc(func(g *Group) {
					//							g.Var().Id("body").Id(strcase.ToLowerCamel(s.Name+ep.MethodName) + "Req")
					//							g.Var().Id("wb").Op("=").Make(Index().Byte(), Lit(0), Lit(10485760)) // 10MB
					//							g.Id("buf").Op(":=").Qual("bytes", "NewBuffer").Call(Id("wb"))
					//							g.List(Id("written"), Id("err")).Op(":=").Qual("io", "Copy").Call(Id("buf"), Id("r").Dot("Body"))
					//							g.Do(gen.CheckErr(
					//								Id("serverErrorEncoder").Call(Id("r").Dot("Context").Call(), Id("rw"), Err()),
					//								Return(),
					//							))
					//							g.Id("data").Op(":=").Id("wb").Index(Op(":").Id("written"))
					//
					//							g.Err().Op("=").Qual("encoding/json", "Unmarshal").Call(Id("data"), Op("&").Id("body"))
					//							g.Do(gen.CheckErr(
					//								Id("serverErrorEncoder").Call(Id("r").Dot("Context").Call(), Id("rw"), Err()),
					//								Return(),
					//							))
					//							for _, p := range ep.BodyParams {
					//								g.Id("param" + p.FldName).Op("=").Id("body").Dot(p.FldName)
					//							}
					//						})
					//					case "urlencoded":
					//						g.Case(Lit("application/x-www-form-urlencoded")).BlockFunc(func(g *Group) {
					//							g.Err().Op("=").Id("r").Dot("ParseForm").Call()
					//							g.Do(gen.CheckErr(
					//								Id("serverErrorEncoder").Call(Id("r").Dot("Context").Call(), Id("rw"), Err()),
					//								Return(),
					//							))
					//							for _, p := range ep.BodyParams {
					//								g.Add(gen.ParseValue(Id("r").Dot("Form").Dot("Get").Call(Lit(p.Name)), Id("param"+p.FldName), "=", p.Type, f.Import))
					//								if b, ok := p.Type.(*types.Basic); (ok && !b.IsString()) || !ok {
					//									g.Do(gen.CheckErr(
					//										Id("serverErrorEncoder").Call(Id("r").Dot("Context").Call(), Id("rw"), Err()),
					//										Return(),
					//									))
					//								}
					//							}
					//						})
					//					case "multipart":
					//						g.Case(Lit("multipart/form-data")).BlockFunc(func(g *Group) {
					//							g.Err().Op("=").Id("r").Dot("ParseMultipartForm").Call(Lit(ep.MultipartMaxMemory))
					//							g.Do(gen.CheckErr(
					//								g.Do(gen.CheckErr(
					//									Id("serverErrorEncoder").Call(Id("r").Dot("Context").Call(), Id("rw"), Err()),
					//									Return(),
					//								)),
					//							))
					//							for _, p := range ep.BodyParams {
					//								g.Add(gen.ParseValue(Id("r").Dot("FormValue").Call(Lit(p.Name)), Id("param"+p.FldName), "=", p.Type, f.Import))
					//								if b, ok := p.Type.(*types.Basic); (ok && !b.IsString()) || !ok {
					//									g.Do(gen.CheckErr(
					//										Id("serverErrorEncoder").Call(Id("r").Dot("Context").Call(), Id("rw"), Err()),
					//										Return(),
					//									))
					//								}
					//							}
					//						})
					//					}
					//				}
					//			})
					//		}
					//	}
					//}
					//
					//g.Do(func(s *Statement) {
					//	s.ListFunc(func(g *Group) {
					//		for _, r := range ep.Results {
					//			g.Id(r.FldName)
					//		}
					//		if ep.Error != nil {
					//			g.Id(ep.Error.Name)
					//		}
					//	})
					//	if len(ep.Results) > 0 {
					//		s.Op(":=")
					//	} else if ep.Error != nil {
					//		s.Op("=")
					//	}
					//	s.Id("svc").Dot(ep.MethodName).CallFunc(func(g *Group) {
					//		if ep.Context != nil {
					//			g.Id("r").Dot("Context").Call()
					//		}
					//		for _, p := range ep.Params {
					//			g.Id("param" + p.FldName)
					//		}
					//	})
					//})
					//
					//if ep.Error != nil {
					//	g.Do(gen.CheckErr(
					//		Id("serverErrorEncoder").Call(Id("r").Dot("Context").Call(), Id("rw"), Id(ep.Error.Name)),
					//		Return(),
					//	))
					//}
					//
					//statusCode := 200
					//
					//if len(ep.BodyResults) > 0 {
					//	g.Var().Id("bufResp").Qual("bytes", "Buffer")
					//	if !ep.NoWrapResponse {
					//		g.Var().Id("bodyResp").StructFunc(gen.WrapResponse(ep.WrapResponse, ep.Results, f.Import))
					//		for _, r := range ep.Results {
					//			g.Id("bodyResp").Do(func(s *Statement) {
					//				for _, name := range ep.WrapResponse {
					//					s.Dot(strcase.ToCamel(name))
					//				}
					//			}).Dot(r.FldNameExport).Op("=").Id(r.Name)
					//		}
					//		g.If(
					//			Err().Op(":=").Qual("encoding/json", "NewEncoder").Call(Op("&").Id("bufResp")).Dot("Encode").Call(Id("bodyResp")),
					//			Err().Op("!=").Nil(),
					//		).Block(
					//			Id("serverErrorEncoder").Call(Id("r").Dot("Context").Call(), Id("rw"), Id(ep.Error.Name)),
					//			Return(),
					//		)
					//	} else if len(ep.BodyResults) == 1 {
					//		g.If(
					//			Err().Op(":=").Qual("encoding/json", "NewEncoder").Call(Op("&").Id("bufResp")).Dot("Encode").Call(Id(ep.BodyResults[0].FldName)),
					//			Err().Op("!=").Nil(),
					//		).Block(
					//			Id("serverErrorEncoder").Call(Id("r").Dot("Context").Call(), Id("rw"), Id(ep.Error.Name)),
					//			Return(),
					//		)
					//	}
					//} else {
					//	statusCode = 204
					//}
					//
					//g.Id("statusCode").Op(":=").Lit(statusCode)
					//
					//g.Id("rw").Dot("WriteHeader").Call(Id("statusCode"))
					//
					//if len(ep.BodyResults) > 0 {
					//	g.If(
					//		List(Id("_"), Err()).Op(":=").Id("rw").Dot("Write").Call(Id("bufResp").Dot("Bytes").Call()),
					//		Err().Op("!=").Nil(),
					//	).Block(
					//		Return(),
					//	)
					//}
				}),
			),
		)
	}

}
