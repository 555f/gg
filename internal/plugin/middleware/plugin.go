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
		nameBaseMiddleware := p.NameBaseMiddleware(iface.Named)

		f.Type().Id(nameMiddleware).Func().Params(jen.Do(f.Qual(iface.Named.Pkg.Path, iface.Named.Name))).Do(f.Qual(iface.Named.Pkg.Path, iface.Named.Name))

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
							jen.Id("next").Do(f.Qual(iface.Named.Pkg.Path, iface.Named.Name)),
						).Do(f.Qual(iface.Named.Pkg.Path, iface.Named.Name)).
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

		f.Type().Id(nameBaseMiddleware).Struct(
			jen.Id("next").Do(f.Qual(iface.Named.Pkg.Path, iface.Named.Name)),
			jen.Id("mediator").Any(),
		)

		for _, method := range iface.Type.Methods {
			var callParams []jen.Code
			for _, param := range method.Sig.Params {
				callParam := jen.Id(param.Name)
				if param.IsVariadic {
					callParam.Op("...")
				}
				callParams = append(callParams, callParam)
			}

			f.Func().
				Params(
					jen.Id("m").Op("*").Id(nameBaseMiddleware),
				).
				Id(method.Name).Add(types.Convert(method.Sig, f.Import)).
				BlockFunc(func(g *jen.Group) {
					g.Defer().Func().Params().BlockFunc(func(g *jen.Group) {
						g.If(jen.List(jen.Id("s"), jen.Id("ok")).Op(":=").Id("m").Dot("mediator").Assert(jen.Id(p.NameBaseMiddlewareMethodIface(iface.Named, method))), jen.Id("ok")).Block(
							jen.Id("s").Dot(method.Name).Call(callParams...),
						)
					}).Call()

					callMethodCode := jen.Id("m").Dot("next").Dot(method.Name).Call(callParams...)

					if method.Sig.Results.HasError() {
						g.Return(
							callMethodCode,
						)
					} else {
						g.Add(callMethodCode)
					}
				})
		}

		for _, method := range iface.Type.Methods {
			f.Type().Id(p.NameBaseMiddlewareMethodIface(iface.Named, method)).Interface(
				jen.Id(method.Name).Add(types.Convert(method.Sig.Params, f.Import)),
			)
		}

		f.Func().Id(strcase.ToCamel(nameBaseMiddleware)).Params(
			jen.Id("mediator").Any(),
		).Id(nameMiddleware).Block(
			jen.Return(jen.Func().Params(jen.Id("next").Do(f.Qual(iface.Named.Pkg.Path, iface.Named.Name))).Params(jen.Do(f.Qual(iface.Named.Pkg.Path, iface.Named.Name)))).Block(
				jen.Return(jen.Op("&").Id(nameBaseMiddleware).Values(
					jen.Id("next").Op(":").Id("next"),
					jen.Id("mediator").Op(":").Id("mediator"),
				)),
			),
		)
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

func (p *Plugin) NameBaseMiddleware(namedType *types.Named) string {
	return strcase.ToLowerCamel(namedType.Name) + "BaseMiddleware"
}

func (p *Plugin) NameBaseMiddlewareMethodIface(namedType *types.Named, method *types.Func) string {
	return strcase.ToLowerCamel(namedType.Name) + method.Name + "BaseMiddleware"
}

func (p *Plugin) Output() string {
	return filepath.Join(p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("output", "internal/middleware/middleware.go"))
}

func (p *Plugin) Dependencies() []string { return nil }
