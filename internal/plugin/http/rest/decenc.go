package rest

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gen"
	"github.com/555f/gg/pkg/types"

	. "github.com/dave/jennifer/jen"
)

func GenEncDec(s options.Iface) func(f *file.GoFile) {
	return func(f *file.GoFile) {
		for _, ep := range s.Endpoints {
			if len(ep.Params) > 0 {
				f.Func().Id(ep.ReqDecodeName).
					Params(
						Id("ctx").Qual("context", "Context"),
						Id("r").Op("*").Qual("net/http", "Request"),
					).
					Params(Id("result").Any(), Err().Error()).
					BlockFunc(func(g *Group) {
						if len(ep.Params) > 0 {
							g.Id("param").Op(":=").New(Id(ep.ReqStructName))
						}
						if len(ep.BodyParams) > 0 {
							switch ep.HTTPMethod {
							case "POST", "PUT", "PATCH", "DELETE":
								g.Id("contentType").Op(":=").Id("r").Dot("Header").Dot("Get").Call(Lit("content-type"))

								g.Id("parts").Op(":=").Qual("strings", "Split").Call(Id("contentType"), Lit(";"))
								g.If(Len(Id("parts")).Op("==").Lit(0)).Block(
									Return(Nil(), Qual("errors", "New").Call(Lit("invalid content type"))),
								)
								g.Id("contentType").Op("=").Id("parts").Index(Lit(0))

								g.Switch(Id("contentType")).BlockFunc(func(g *Group) {
									g.Default().Block(
										Return(Nil(), Qual("errors", "New").Call(Lit("invalid content type"))),
									)
									for _, contentType := range ep.ContentTypes {
										switch contentType {
										case "xml":
											g.Case(Lit("application/xml")).BlockFunc(func(g *Group) {
												stParams := []Code{Id("XMLName").Qual("encoding/xml", "Name").Tag(map[string]string{"xml": ep.ReqRootXMLName})}

												for _, p := range ep.BodyParams {
													st := Id(p.FldName).Add(types.Convert(p.Type, f.Import))
													if p.Name != "" {
														st.Tag(map[string]string{"xml": p.Name})
													}
													stParams = append(stParams, st)
												}
												g.Var().Id("body").Struct(stParams...)
												g.Var().Id("data").Index().Byte()
												g.List(Id("data"), Err()).
													Op("=").
													Qual("io", "ReadAll").
													Call(Id("r").Dot("Body"))
												g.Do(gen.CheckErr(
													Return(),
												))
												g.Err().Op("=").Qual("encoding/xml", "Unmarshal").Call(Id("data"), Op("&").Id("body"))
												g.Do(gen.CheckErr(
													Return(),
												))
												for _, p := range ep.BodyParams {
													g.Id("param").Dot(p.FldNameUnExport).Op("=").Id("body").Dot(p.FldName)
												}
											})
										case "json":
											g.Case(Lit("application/json")).BlockFunc(func(g *Group) {
												var stParams []Code
												for _, p := range ep.BodyParams {
													st := Id(p.FldName).Add(types.Convert(p.Type, f.Import))
													if p.Name != "" {
														st.Tag(map[string]string{"json": p.Name})
													}
													stParams = append(stParams, st)
												}
												g.Var().Id("body").Struct(stParams...)
												g.Var().Id("data").Index().Byte()
												g.List(Id("data"), Err()).
													Op("=").
													Qual("io", "ReadAll").
													Call(Id("r").Dot("Body"))
												g.Do(gen.CheckErr(
													Return(),
												))
												g.Err().Op("=").Qual("encoding/json", "Unmarshal").Call(Id("data"), Op("&").Id("body"))
												g.Do(gen.CheckErr(
													Return(),
												))
												for _, p := range ep.BodyParams {
													g.Id("param").Dot(p.FldNameUnExport).Op("=").Id("body").Dot(p.FldName)
												}
											})
										case "urlencoded":
											g.Case(Lit("application/x-www-form-urlencoded")).BlockFunc(func(g *Group) {
												g.Err().Op("=").Id("r").Dot("ParseForm").Call()
												g.Do(gen.CheckErr(Return()))
												for _, p := range ep.BodyParams {
													g.Add(gen.ParseValue(Id("r").Dot("Form").Dot("Get").Call(Lit(p.Name)), Id("param").Dot(p.FldNameUnExport), "=", p.Type, f.Import))
													if b, ok := p.Type.(*types.Basic); (ok && !b.IsString()) || !ok {
														g.Do(gen.CheckErr(Return()))
													}
												}
											})
										case "multipart":
											g.Case(Lit("multipart/form-data")).BlockFunc(func(g *Group) {
												g.Err().Op("=").Id("r").Dot("ParseMultipartForm").Call(Lit(ep.MultipartMaxMemory))
												g.Do(gen.CheckErr(Return()))
												for _, p := range ep.BodyParams {
													g.Add(gen.ParseValue(Id("r").Dot("FormValue").Call(Lit(p.Name)), Id("param").Dot(p.FldNameUnExport), "=", p.Type, f.Import))
													if b, ok := p.Type.(*types.Basic); (ok && !b.IsString()) || !ok {
														g.Do(gen.CheckErr(Return()))
													}
												}
											})
										}
									}
								})
							}
						}
						if len(ep.PathParams) > 0 || len(ep.QueryParams) > 0 {
							g.Id("q").Op(":=").Id("r").Dot("URL").Dot("Query").Call()

							for _, p := range ep.PathParams {
								g.If(Id("s").Op(":=").Id("q").Dot("Get").Call(Lit(p.Name)), Id("s").Op("!=").Lit("")).Block(
									Add(gen.ParseValue(Id("s"), Id("param").Dot(p.FldNameUnExport), "=", p.Type, f.Import)),
									Do(gen.CheckErr(Return())),
								)
							}

							if len(ep.QueryParams) > 0 {
								for _, param := range ep.QueryParams {
									g.Add(makeGetQueryParam(param.Parent, param, f.Import))
								}
							}
						}

						if len(ep.HeaderParams) > 0 {
							for _, p := range ep.HeaderParams {
								g.If(Id("s").Op(":=").Id("r").Dot("Header").Dot("Get").Call(Lit(p.Name)), Id("s").Op("!=").Lit("")).Block(
									Add(gen.ParseValue(Id("s"), Id("param").Dot(p.FldNameUnExport), "=", p.Type, f.Import)),
									Do(gen.CheckErr(Return())),
								)
							}
						}
						if len(ep.Params) > 0 {
							g.Return(Id("param"), Nil())
						} else {
							g.Return()
						}
					})
			}
		}
	}
}
