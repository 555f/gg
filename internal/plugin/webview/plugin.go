package middleware

import (
	"path/filepath"

	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

var (
	webviewPkg = "github.com/webview/webview_go"
)

type Plugin struct {
	ctx *gg.Context
}

func (p *Plugin) Name() string { return "webview" }

func (p *Plugin) Exec() ([]file.File, error) {
	f := file.NewGoFile(p.ctx.Module, p.Output())

	for _, iface := range p.ctx.Interfaces {
		optionsName := iface.Named.Name + "Options"
		optionName := iface.Named.Name + "Option"

		f.Type().Id(optionName).Func().Params(jen.Op("*").Id(optionsName))
		f.Type().Id(optionsName).StructFunc(func(g *jen.Group) {
		})

		/*for _, m := range iface.Type.Methods {
			bindName := strcase.ToLowerCamel(iface.Named.Name) + "_" + m.Name
			resultName := bindName + "Result"

			f.Type().Id(resultName).StructFunc(func(g *jen.Group) {})
		}*/
	}

	for _, iface := range p.ctx.Interfaces {
		f.Func().Id("SetupRoutes"+iface.Named.Name).Params(
			jen.Id("svc").Do(f.Qual(iface.Named.Pkg.Path, iface.Named.Name)),
			jen.Id("w").Qual(webviewPkg, "WebView"),
			jen.Id("opts").Op("...").Id(iface.Named.Name+"Option"),
		).BlockFunc(func(g *jen.Group) {
			for _, m := range iface.Type.Methods {
				bindName := strcase.ToLowerCamel(iface.Named.Name) + "_" + m.Name

				g.Id("w").Dot("Bind").Call(
					jen.Lit(bindName), jen.Func().ParamsFunc(func(g *jen.Group) {
						for _, p := range m.Sig.Params {
							if p.IsContext {
								continue
							}
							g.Id(p.Name).Add(types.Convert(p.Type, f.Qual))
						}
					}).Params(jen.Id("_").Any(), jen.Err().Error()).BlockFunc(func(g *jen.Group) {
						g.Id("result").Op(":=").StructFunc(func(g *jen.Group) {
							for _, r := range m.Sig.Results {
								if r.IsError {
									continue
								}
								g.Id(strcase.ToCamel(r.Name)).Add(types.Convert(r.Type, f.Qual)).Tag(map[string]string{"json": strcase.ToLowerCamel(r.Name)})
							}
						}).Values()

						g.Id("ctx").Op(":=").Qual("context", "TODO").Call()

						g.ListFunc(func(g *jen.Group) {
							for _, r := range m.Sig.Results {
								if r.IsError {
									g.Err()
									continue
								}
								g.Id("result").Dot(strcase.ToCamel(r.Name))
							}
						}).Op("=").Id("svc").Dot(m.Name).CallFunc(func(g *jen.Group) {
							for _, p := range m.Sig.Params {
								if p.IsContext {
									g.Id("ctx")
									continue
								}
								g.Id(p.Name)
							}
						})
						g.Return(
							jen.Id("result"),
							jen.Err(),
						)
					}),
				)
			}
		})
	}

	return []file.File{f}, nil
}

func (p *Plugin) Output() string {
	return filepath.Join(p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("output", "internal/server/server.go"))
}

func (p *Plugin) Dependencies() []string { return nil }
