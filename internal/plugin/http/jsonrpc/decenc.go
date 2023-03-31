package jsonrpc

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gen"

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
						Id("params").Qual("encoding/json", "RawMessage"),
					).
					Params(Id("result").Any(), Err().Error()).
					BlockFunc(func(g *Group) {
						if len(ep.Params) > 0 {
							g.Id("body").Op(":=").New(Id(ep.ReqStructName))
						}

						if len(ep.BodyParams) > 0 {
							//var stParams []Code
							//for _, p := range ep.BodyParams {
							//	st := Id(p.FldNameExport).Add(types.Convert(p.HTTPType, f.Import))
							//	if p.Name != "" {
							//		st.Tag(map[string]string{"json": p.Name})
							//	}
							//	stParams = append(stParams, st)
							//}
							//g.Var().Id("body").Struct(stParams...)
							//g.Var().Id("data").Index().Byte()
							//g.List(Id("data"), Err()).
							//	Op("=").
							//	Qual("io", "ReadAll").
							//	Call(Id("prams"))
							//g.Do(checkErr(
							//	Return(),
							//))
							g.Err().Op("=").Qual("encoding/json", "Unmarshal").Call(Id("params"), Op("&").Id("body"))
							g.Do(gen.CheckErr(
								Return(),
							))
							//for _, p := range ep.BodyParams {
							//	paramID := g.Id("param")
							//	if p.Parent != nil {
							//		paramID.Dot(p.Parent.Name)
							//	}
							//	paramID.Dot(p.FldNameExport).Op("=").Id("body").Dot(p.FldNameExport)
							//}

							////switch ep.Method {
							////case "POST", "PUT", "PATCH", "DELETE":
							//g.Id("contentType").Op(":=").Id("r").Dot("Header").Dot("Get").Call(Lit("content-type"))
							//
							//g.Id("parts").Op(":=").Qual("strings", "Split").Call(Id("contentType"), Lit(";"))
							//g.If(Len(Id("parts")).Op("==").Lit(0)).Block(
							//	Return(Nil(), Qual("errors", "New").Call(Lit("invalid content type"))),
							//)
							//g.Id("contentType").Op("=").Id("parts").Index(Lit(0))
							//
							//g.Switch(Id("contentType")).BlockFunc(func(g *Group) {
							//	g.Default().Block(
							//		Return(Nil(), Qual("errors", "New").Call(Lit("invalid content type"))),
							//	)
							//	for _, bodyType := range ep.ContentTypes {
							//		switch bodyType {
							//		case "xml":
							//			g.Case(Lit("application/xml")).BlockFunc(func(g *Group) {})
							//		case "json":
							//			g.Case(Lit("application/json")).BlockFunc(func(g *Group) {
							//
							//			})
							//		}
							//	}
							//})
							////}
						}

						//if len(ep.QueryParams) > 0 {
						//	for _, p := range ep.QueryParams {
						//		paramID := Id("body")
						//		if p.Parent != nil {
						//			paramID.Dot(p.Parent.FldName)
						//		}
						//		paramID.Dot(p.FldName)
						//		g.If(Id("s").Op(":=").Id("q").Dot("Get").Call(Lit(p.Name)), Id("s").Op("!=").Lit("")).Block(
						//			Add(gen.ParseValue(Id("s"), paramID, "=", p.Type, f.Import)),
						//			Do(gen.CheckErr(Return())),
						//		)
						//	}
						//}

						if len(ep.HeaderParams) > 0 {
							for _, p := range ep.HeaderParams {
								g.If(Id("s").Op(":=").Id("r").Dot("Header").Dot("Get").Call(Lit(p.Name)), Id("s").Op("!=").Lit("")).Block(
									Add(gen.ParseValue(Id("s"), Id("body").Dot(p.FldName), "=", p.Type, f.Import)),
									Do(gen.CheckErr(Return())),
								)
							}
						}
						if len(ep.Params) > 0 {
							g.Return(Id("body"), Nil())
						} else {
							g.Return()
						}
					})
			}
		}
	}
}
