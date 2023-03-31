package middleware

import (
	"path/filepath"

	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"

	"github.com/dave/jennifer/jen"
)

type Plugin struct {
	output string
}

func (p *Plugin) Name() string { return "middleware" }

func (p *Plugin) Exec(ctx *gg.Context) ([]file.File, error) {
	p.output = filepath.Join(ctx.Module.Dir, ctx.Options.GetStringWithDefault("output", "internal/middleware/middleware.go"))
	f := file.NewGoFile(ctx.Module, p.output)
	for _, iface := range ctx.Interfaces {
		nameMiddleware := p.NameMiddleware(iface.Named)
		nameMiddlewareChain := p.NameMiddlewareChain(iface.Named)

		f.Type().Id(nameMiddleware).Func().Params(jen.Qual(iface.Named.Pkg.Path, iface.Named.Name)).Qual(iface.Named.Pkg.Path, iface.Named.Name)

		f.Func().
			Id(nameMiddlewareChain).
			Params(
				jen.Id("outer").Id(nameMiddleware),
				jen.Id("others").Op("...").Id(nameMiddleware),
			).
			Id(nameMiddleware).
			BlockFunc(func(group *jen.Group) {
				group.ReturnFunc(func(group *jen.Group) {
					group.Func().
						Params(
							jen.Id("next").Qual(iface.Named.Pkg.Path, iface.Named.Name),
						).
						Qual(iface.Named.Pkg.Path, iface.Named.Name).
						BlockFunc(func(group *jen.Group) {
							group.For(
								jen.Id("i").Op(":=").Len(jen.Id("others")).Op("-1"),
								jen.Id("i").Op(">=").Lit(0),
								jen.Id("i").Op("--"),
							).BlockFunc(func(group *jen.Group) {
								group.Id("next").Op("=").Id("others").Index(jen.Id("i")).Call(jen.Id("next"))
							})
							group.Return(jen.Id("outer").Call(jen.Id("next")))
						})
				})
			})
	}
	return []file.File{f}, nil
}

func (p *Plugin) NameMiddlewareChain(named *types.Named) string {
	return strcase.ToCamel(named.Name) + "MiddlewareChain"
}

func (p *Plugin) NameMiddleware(namedType *types.Named) string {
	return strcase.ToCamel(namedType.Name) + "Middleware"
}

func (p *Plugin) Output() string {
	return p.output
}

func (p *Plugin) Dependencies() []string { return nil }
