package rest

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gen"
	"github.com/555f/gg/pkg/strcase"

	. "github.com/dave/jennifer/jen"
)

func GenServerOption(s options.Iface) func(f *file.GoFile) {
	return func(f *file.GoFile) {
		serviceOptionName := s.Name + "ServerOption"
		f.Commentf("// %s apply service", s.Name)
		f.Func().Id(s.Name).Params(
			Id("svc").Do(f.Import(s.PkgPath, s.Name)),
			Id("opts").Op("...").Id(serviceOptionName),
		).Id("ServerOption").Block(
			Return().Func().Params(Id("s").Op("*").Id("Server")).BlockFunc(func(g *Group) {
				optionName := strcase.ToLowerCamel(s.Name)
				optionsName := s.Name + "Options"

				g.Id(optionName).Op(":=").Op("&").Id(optionsName).Values()

				g.For(List(Id("_"), Id("opt"))).Op(":=").Range().Id("opts").Block(
					Id("opt").Call(Id(optionName)),
				)
				for _, ep := range s.Endpoints {
					lcName := strcase.ToLowerCamel(ep.MethodName)

					g.Id(lcName+"Opts").Op(":=").Index().Id("RouteOption").Values(
						Id("Before").Call(Append(Id(optionName).Dot("before"), Id(optionName).Dot(lcName).Dot("before").Op("...")).Op("...")),
						Id("After").Call(Append(Id(optionName).Dot("after"), Id(optionName).Dot(lcName).Dot("after").Op("...")).Op("...")),
						Id("Middleware").Call(Append(Id(optionName).Dot("serverMiddleware"), Id(optionName).Dot(lcName).Dot("serverMiddleware").Op("...")).Op("...")),
					)

					g.Id("s").Dot("AddRoute").Call(
						Lit(ep.HTTPMethod),

						Do(func(s *Statement) {
							if len(ep.PathParams) == 0 {
								s.Lit(ep.Pattern)
							} else {
								s.Id("Regex").Values(
									Id("Pattern").Op(":").Lit(ep.Pattern),
									Id("Params").Op(":").Index().String().ValuesFunc(func(g *Group) {
										for _, name := range ep.ParamsNameIdx {
											g.Lit(name)
										}
									}),
								)
							}
						}),
						Qual("net/http", "HandlerFunc").Call(
							Func().Params(
								Id("rw").Qual("net/http", "ResponseWriter"),
								Id("r").Op("*").Qual("net/http", "Request"),
							).BlockFunc(func(g *Group) {
								reqData := Id("reqData")
								if len(ep.Params) > 0 {
									g.List(reqData, Err()).Op(":=").Id(ep.ReqDecodeName).Call(Id("r").Dot("Context").Call(), Id("r"))
									g.Do(gen.CheckErr(
										Id("serverErrorEncoder").Call(Id("r").Dot("Context").Call(), Id("rw"), Err()),
										Return(),
									))
								} else {
									reqData = Nil()
								}
								g.List(Id("resp"), Err()).Op(":=").Id("middlewareChain").Call(
									Append(
										Id(optionName).Dot("middleware"),
										Id(optionName).Dot(lcName).Dot("middleware").Op("..."),
									),
								).Call(Id(ep.Name).Call(Id("svc"))).Call(Id("r").Dot("Context").Call(), reqData)
								g.Do(gen.CheckErr(
									Id("serverErrorEncoder").Call(Id("r").Dot("Context").Call(), Id("rw"), Err()),
									Return(),
								))
								g.Id("encodeJSONResponse").Call(Id("r").Dot("Context").Call(), Id("rw"), Id("resp"))
							}),
						),
						Id(lcName+"Opts").Op("..."),
					)
				}
			}),
		).Line()
	}
}
