package middleware

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"

	"github.com/dave/jennifer/jen"
)

type Plugin struct {
	ctx *gg.Context
}

func (p *Plugin) Name() string { return "middleware" }

func (p *Plugin) Exec() ([]file.File, error) {
	f := file.NewGoFile(p.ctx.Module, p.Output())
	for _, iface := range p.ctx.Interfaces {
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

func (p *Plugin) PkgPath(named *types.Named) string {
	return path.Dir(path.Join(p.ctx.Module.Path, strings.Replace(p.Output(), p.ctx.Workdir, "", -1)))
}

func (p *Plugin) NameMiddlewareChain(named *types.Named) string {
	return strcase.ToCamel(named.Name) + "MiddlewareChain"
}

func (p *Plugin) NameMiddleware(namedType *types.Named) string {
	return strcase.ToCamel(namedType.Name) + "Middleware"
}

func (p *Plugin) Output() string {
	return filepath.Join(p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("output", "internal/middleware/middleware.go"))
}

func (p *Plugin) Dependencies() []string { return nil }
