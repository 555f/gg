package rest

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/file"

	. "github.com/dave/jennifer/jen"
)

func GenOptions(s options.Iface) func(f *file.GoFile) {
	return func(f *file.GoFile) {
		var middlewareType Code

		switch s.Type {
		case "rest":
			switch s.Lib {
			case "http":
				middlewareType = Func().Params(
					Qual("net/http", "Handler"),
				).Qual("net/http", "Handler")

			case "echo":
				middlewareType = Qual("github.com/labstack/echo/v4", "MiddlewareFunc")
			}
		}

		optionsName := s.Name + "Options"
		optionName := s.Name + "Option"

		f.Type().Id(optionName).Func().Params(Op("*").Id(optionsName))

		f.Type().Id(optionsName).StructFunc(func(g *Group) {

			g.Id("middleware").Index().Add(middlewareType)
			for _, ep := range s.Endpoints {
				g.Id("middleware" + ep.MethodName).Index().Add(middlewareType)
			}
		})

		f.Func().Id(s.Name + "Middleware").Params(Id("middleware").Op("...").Add(middlewareType)).Id(optionName).Block(
			Return(
				Func().Params(Id("o").Op("*").Id(optionsName)).Block(
					Id("o").Dot("middleware").Op("=").Append(Id("o").Dot("middleware"), Id("middleware").Op("...")),
				),
			),
		)

		for _, ep := range s.Endpoints {
			f.Func().Id(s.Name + ep.MethodName + "Middleware").Params(Id("middleware").Op("...").Add(middlewareType)).Id(optionName).Block(
				Return(
					Func().Params(Id("o").Op("*").Id(optionsName)).Block(
						Id("o").Dot("middleware"+ep.MethodName).Op("=").Append(Id("o").Dot("middleware"+ep.MethodName), Id("middleware").Op("...")),
					),
				),
			)
		}

	}
}
