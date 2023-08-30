package middleware

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"

	. "github.com/dave/jennifer/jen"
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

		f.Type().Id(nameMiddleware).Func().Params(Qual(iface.Named.Pkg.Path, iface.Named.Name)).Qual(iface.Named.Pkg.Path, iface.Named.Name)

		f.Func().
			Id(nameMiddlewareChain).
			Params(
				Id("outer").Id(nameMiddleware),
				Id("others").Op("...").Id(nameMiddleware),
			).
			Id(nameMiddleware).
			BlockFunc(func(group *Group) {
				group.ReturnFunc(func(group *Group) {
					group.Func().
						Params(
							Id("next").Qual(iface.Named.Pkg.Path, iface.Named.Name),
						).
						Qual(iface.Named.Pkg.Path, iface.Named.Name).
						BlockFunc(func(group *Group) {
							group.For(
								Id("i").Op(":=").Len(Id("others")).Op("-1"),
								Id("i").Op(">=").Lit(0),
								Id("i").Op("--"),
							).BlockFunc(func(group *Group) {
								group.Id("next").Op("=").Id("others").Index(Id("i")).Call(Id("next"))
							})
							group.Return(Id("outer").Call(Id("next")))
						})
				})
			})

		f.Type().Id(nameBaseMiddleware).Struct(
			Id("next").Qual(iface.Named.Pkg.Path, iface.Named.Name),
			Id("mediator").Any(),
		)

		for _, method := range iface.Type.Methods {
			var callParams []Code
			for _, param := range method.Sig.Params {
				callParam := Id(param.Name)
				if param.IsVariadic {
					callParam.Op("...")
				}
				callParams = append(callParams, callParam)
			}

			f.Func().
				Params(
					Id("m").Op("*").Id(nameBaseMiddleware),
				).
				Id(method.Name).Add(types.Convert(method.Sig, f.Import)).
				BlockFunc(func(g *Group) {
					g.Defer().Func().Params().BlockFunc(func(g *Group) {
						g.If(List(Id("s"), Id("ok")).Op(":=").Id("m").Dot("mediator").Assert(Id(p.NameBaseMiddlewareMethodIface(iface.Named, method))), Id("ok")).Block(
							Id("s").Dot(method.Name).Call(callParams...),
						)
					}).Call()
					g.Return(
						Id("m").Dot("next").Dot(method.Name).Call(callParams...),
					)
				})
		}

		for _, method := range iface.Type.Methods {
			f.Type().Id(p.NameBaseMiddlewareMethodIface(iface.Named, method)).Interface(
				Id(method.Name).Add(types.Convert(method.Sig.Params, f.Import)),
			)
		}

		f.Func().Id(strcase.ToCamel(nameBaseMiddleware)).Params(
			Id("mediator").Any(),
		).Id(nameMiddleware).Block(
			Return(Func().Params(Id("next").Qual(iface.Named.Pkg.Path, iface.Named.Name)).Params(Qual(iface.Named.Pkg.Path, iface.Named.Name))).Block(
				Return(Op("&").Id(nameBaseMiddleware).Values(
					Id("next").Op(":").Id("next"),
					Id("mediator").Op(":").Id("mediator"),
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
