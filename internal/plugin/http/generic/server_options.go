package generic

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/strcase"

	. "github.com/dave/jennifer/jen"
)

func GenServerOption(s options.Iface) func(f *file.GoFile) {
	return func(f *file.GoFile) {
		optionsName := s.Name + "Options"
		optionName := s.Name + "ServerOption"

		f.Type().Id(optionName).Func().Params(Op("*").Id(optionsName))
		f.Type().Id(optionsName).StructFunc(func(g *Group) {
			g.Id("serverOptions")
			for _, endpoint := range s.Endpoints {
				g.Id(strcase.ToLowerCamel(endpoint.MethodName)).Id("serverOptions")
			}
		}).Line()

		f.Func().Id(s.Name + "ApplyOptions").Params(Id("options").Op("...").Id(optionName)).Id(optionName).Block(
			Return(
				Func().Params(Id("o").Op("*").Id(optionsName)).Block(
					For(List(Id("_"), Id("opt")).Op(":=")).Range().Id("options").Block(
						Id("opt").Call(Id("o")),
					),
				),
			),
		)
		f.Func().Id(s.Name + "Before").Params(Id("before").Op("...").Id("ServerBeforeFunc")).Id(optionName).Block(
			Return(
				Func().Params(Id("o").Op("*").Id(optionsName)).Block(
					Id("o").Dot("before").Op("=").Append(Id("o").Dot("before"), Id("before").Op("...")),
				),
			),
		)
		f.Func().Id(s.Name + "After").Params(Id("after").Op("...").Id("ServerAfterFunc")).Id(optionName).Block(
			Return(
				Func().Params(Id("o").Op("*").Id(optionsName)).Block(
					Id("o").Dot("after").Op("=").Append(Id("o").Dot("after"), Id("after").Op("...")),
				),
			),
		)
		f.Func().Id(s.Name + "Middleware").Params(Id("middleware").Op("...").Id("EndpointMiddleware")).Id(optionName).Block(
			Return(
				Func().Params(Id("o").Op("*").Id(optionsName)).Block(
					Id("o").Dot("middleware").Op("=").Append(Id("o").Dot("middleware"), Id("middleware").Op("...")),
				),
			),
		)

		f.Func().Id(s.Name + "ServerMiddleware").Params(Id("serverMiddleware").Op("...").Id("ServerMiddlewareFunc")).Id(optionName).Block(
			Return(
				Func().Params(Id("o").Op("*").Id(optionsName)).Block(
					Id("o").Dot("serverMiddleware").Op("=").Append(Id("o").Dot("serverMiddleware"), Id("serverMiddleware").Op("...")),
				),
			),
		)

		for _, endpoint := range s.Endpoints {
			lcName := strcase.ToLowerCamel(endpoint.MethodName)
			f.Func().Id(s.Name + endpoint.MethodName + "Before").Params(Id("before").Op("...").Id("ServerBeforeFunc")).Id(optionName).Block(
				Return(
					Func().Params(Id("o").Op("*").Id(optionsName)).Block(
						Id("o").Dot(lcName).Dot("before").Op("=").Append(Id("o").Dot(lcName).Dot("before"), Id("before").Op("...")),
					),
				),
			)
			f.Func().Id(s.Name + endpoint.MethodName + "After").Params(Id("after").Op("...").Id("ServerAfterFunc")).Id(optionName).Block(
				Return(
					Func().Params(Id("o").Op("*").Id(optionsName)).Block(
						Id("o").Dot(lcName).Dot("after").Op("=").Append(Id("o").Dot(lcName).Dot("after"), Id("after").Op("...")),
					),
				),
			)
			f.Func().Id(s.Name + endpoint.MethodName + "Middleware").Params(Id("middleware").Op("...").Id("EndpointMiddleware")).Id(optionName).Block(
				Return(
					Func().Params(Id("o").Op("*").Id(optionsName)).Block(
						Id("o").Dot(lcName).Dot("middleware").Op("=").Append(Id("o").Dot(lcName).Dot("middleware"), Id("middleware").Op("...")),
					),
				),
			)
			f.Func().Id(s.Name + endpoint.MethodName + "ServerMiddlewareFunc").Params(Id("serverMiddleware").Op("...").Id("ServerMiddlewareFunc")).Id(optionName).Block(
				Return(
					Func().Params(Id("o").Op("*").Id(optionsName)).Block(
						Id("o").Dot(lcName).Dot("serverMiddleware").Op("=").Append(Id("o").Dot(lcName).Dot("serverMiddleware"), Id("serverMiddleware").Op("...")),
					),
				),
			)
		}
	}
}
