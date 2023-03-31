package jsonrpc

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gen"
	"github.com/555f/gg/pkg/strcase"

	. "github.com/dave/jennifer/jen"
)

func GenServerOption(s options.Iface) func(f *file.GoFile) {
	return func(f *file.GoFile) {
		lcServiceName := strcase.ToLowerCamel(s.Name)
		optionName := lcServiceName
		serviceOptionName := s.Name + "ServerOption"
		f.Commentf("// %s apply service", s.Name)
		f.Func().Id(s.Name).Params(
			Id("svc").Qual(s.PkgPath, s.Name),
			Id("opts").Op("...").Id(serviceOptionName),
		).Id("ServerOption").Block(
			Return().Func().Params(Id("s").Op("*").Id("Server")).BlockFunc(func(g *Group) {
				optionsName := s.Name + "Options"

				g.Id(optionName).Op(":=").Op("&").Id(optionsName).Values()

				g.For(List(Id("_"), Id("opt"))).Op(":=").Range().Id("opts").Block(
					Id("opt").Call(Id(optionName)),
				)
				for _, ep := range s.Endpoints {
					lcName := strcase.ToLowerCamel(ep.MethodName)
					g.Id("s").Dot("addRoute").Call(
						Lit(ep.Path),
						Func().Params(
							Id("ctx").Qual("context", "Context"),
							Id("w").Qual("net/http", "ResponseWriter"),
							Id("r").Op("*").Qual("net/http", "Request"),
							Id("params").Qual("encoding/json", "RawMessage"),
						).Params(Any(), Error()).BlockFunc(func(g *Group) {
							reqData := Id("reqData")
							if len(ep.Params) > 0 {
								g.List(reqData, Err()).Op(":=").Id(ep.ReqDecodeName).Call(Id("ctx"), Id("r"), Id("params"))
								g.Do(gen.CheckErr(
									//Id("serverErrorEncoder").Call(Id("ctx"), Id("w"), Err()),
									Return(Nil(), Err()),
								))
							} else {
								reqData = Nil()
							}
							g.List(Id("resp"), Err()).Op(":=").Id("middlewareChain").Call(
								Append(
									Id(optionName).Dot("middleware"),
									Id(optionName).Dot(lcName).Dot("middleware").Op("..."),
								),
							).Call(Id(ep.Name).Call(Id("svc"))).Call(Id("ctx"), reqData)
							g.Do(gen.CheckErr(
								//Id("serverErrorEncoder").Call(Id("ctx"), Id("w"), Err()),
								Return(Nil(), Err()),
							))
							//g.Id("encodeJSONResponse").Call(Id("ctx"), Id("w"), Id("resp"))

							g.Return(Id("resp"), Nil())
						}),
						Append(Id(optionName).Dot("before"), Id(optionName).Dot(lcName).Dot("before").Op("...")),
						Append(Id(optionName).Dot("after"), Id(optionName).Dot(lcName).Dot("after").Op("...")),
					)
				}
			}),
		).Line()
	}
}
