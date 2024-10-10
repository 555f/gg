package webview

import (
	"path/filepath"

	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
	"github.com/hashicorp/go-multierror"
)

var (
	webviewPkg = "github.com/webview/webview_go"
)

type Plugin struct {
	ctx *gg.Context
}

func (p *Plugin) Name() string { return "webview" }

func (p *Plugin) Exec() (files []file.File, errs error) {
	f := file.NewGoFile(p.ctx.Module, p.Output())

	f.Type().Id("Context").Map(jen.String()).Any()

	f.Func().Params(jen.Id("c").Id("Context")).Id("Get").Params(jen.Id("key").String()).Any().Block(
		jen.List(jen.Id("v"), jen.Id("_")).Op(":=").Id("c").Index(jen.Id("key")),
		jen.Return(jen.Id("v")),
	)

	for _, iface := range p.ctx.Interfaces {
		optionsName := iface.Named.Name + "Options"
		optionName := iface.Named.Name + "Option"
		f.Type().Id(optionName).Func().Params(jen.Op("*").Id(optionsName))
		f.Type().Id(optionsName).StructFunc(func(g *jen.Group) {})
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
						g.Id("wvCtx").Id("Context")
						for _, p := range m.Sig.Params {
							if p.IsContext {
								continue
							}
							g.Id(p.Name).Add(types.Convert(p.Type, f.Qual))
						}
					}).Params(jen.Chan().Qual(webviewPkg, "BindCallbackResult")).BlockFunc(func(g *jen.Group) {
						g.Id("ch").Op(":=").Make(jen.Chan().Qual(webviewPkg, "BindCallbackResult"))

						g.Go().Func().Params().BlockFunc(func(g *jen.Group) {
							g.Id("result").Op(":=").StructFunc(func(g *jen.Group) {
								for _, r := range m.Sig.Results {
									if r.IsError {
										continue
									}
									g.Id(strcase.ToCamel(r.Name)).Add(types.Convert(r.Type, f.Qual)).Tag(map[string]string{"json": strcase.ToLowerCamel(r.Name)})
								}
							}).Values()

							g.Id("ctx").Op(":=").Qual("context", "TODO").Call()

							tags := m.Tags.GetSlice("webview-context")
							for _, t := range tags {
								if t.Value == "" {
									errs = multierror.Append(errs, errors.Error("the path to the context key is required", t.Position))
									return
								}
								pkgPath, name, err := p.ctx.Module.ParseImportPath(t.Value)
								if err != nil {
									errs = multierror.Append(errs, err)
									return
								}
								g.Id("ctx").Op("=").Qual("context", "WithValue").Call(
									jen.Id("ctx"),
									jen.Qual(pkgPath, name),
									jen.Id("wvCtx").Dot("Get").Call(jen.Lit(name)),
								)
							}

							g.Var().Err().Error()

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

							g.Id("ch").Op("<-").Qual(webviewPkg, "BindCallbackResult").Values(
								jen.Id("Value").Op(":").Id("result"),
								jen.Id("Error").Op(":").Id("err"),
							)
						}).Call()
						g.Return(jen.Id("ch"))
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
