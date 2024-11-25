package metrics

import (
	"fmt"
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
	prometheusPkg = "github.com/prometheus/client_golang/prometheus"
	timePkg       = "time"
)

type Plugin struct {
	ctx *gg.Context
}

func (p *Plugin) Name() string { return "metrics" }

func (p *Plugin) Exec() (files []file.File, errs error) {
	output := filepath.Join(p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("output", "internal/metrics/metrics.go"))

	f := file.NewGoFile(p.ctx.Module, output)

	for _, iface := range p.ctx.Interfaces {
		for _, m := range iface.Type.Methods {
			constName := strcase.ToLowerCamel(iface.Named.Name + m.Name)
			f.Const().Id(constName).Op("=").Lit(m.FullName)
			f.Const().Id(constName + "Short").Op("=").Lit(m.ShortName)
		}
	}

	f.Type().Id("prometheusCollector").Interface(
		jen.Qual(prometheusPkg, "Collector"),
		jen.Id("Requests").Params().Params(jen.Op("*").Qual(prometheusPkg, "CounterVec")),
		jen.Id("ErrRequests").Params().Params(jen.Op("*").Qual(prometheusPkg, "CounterVec")),
		jen.Id("Duration").Params().Params(jen.Op("*").Qual(prometheusPkg, "HistogramVec")),
	)

	for _, iface := range p.ctx.Interfaces {
		nameConstructor := "Metrics" + strcase.ToCamel(iface.Named.Name) + "Middleware"
		nameStruct := strcase.ToCamel(iface.Named.Name) + "MetricMiddleware"
		nameConstScopeName := strcase.ToLowerCamel(iface.Named.Name) + "ScopeName"
		scopeName := strcase.ToSnake(iface.Named.Name)

		tag, ok := iface.Named.Tags.Get("metrics-middleware-type")
		if !ok {
			errs = multierror.Append(errs, fmt.Errorf("interface %s.%s not set `metrics-middleware-type` option", iface.Named.Pkg.Path, iface.Named.Name))
			continue
		}
		pkgMiddleware, nameMiddleware, err := p.ctx.Module.ParseImportPath(tag.Value)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("interface %s.%s option `metrics-middleware-type` invalid format: %s", iface.Named.Pkg.Path, iface.Named.Name, err.Error()))
			continue
		}

		f.Const().Id(nameConstScopeName).Op("=").Lit(scopeName)

		f.Type().Id(nameStruct).StructFunc(func(g *jen.Group) {
			g.Id("next").Do(f.Qual(iface.Named.Pkg.Path, iface.Named.Name))
			g.Id("c").Id("prometheusCollector")
		})
		for _, method := range iface.Type.Methods {
			opts, err := makeMethodOptions(p.ctx.Module, method)
			if err != nil {
				errs = multierror.Append(errs, err)
				continue
			}
			if opts.Skip {
				errs = multierror.Append(errs, errors.Warn("the method is marked as skipped", method.Position))
				continue
			}
			var (
				callParams []jen.Code
				results    []jen.Code
				errorVars  []*types.Var
				contextVar *types.Var
			)

			for _, param := range method.Sig.Params {
				callParam := jen.Id(param.Name)
				if param.IsVariadic {
					callParam.Op("...")
				}
				callParams = append(callParams, callParam)

				if param.IsContext {
					if contextVar != nil {
						errs = multierror.Append(errs, errors.Warn(fmt.Sprintf("found another parameter %s with the type \"context.Context\"", param.Name), param.Position))
						continue
					}
					contextVar = param
					continue
				}
			}
			for _, result := range method.Sig.Results {
				results = append(results, jen.Id(result.Name))
				if result.IsError {
					errorVars = append(errorVars, result)
				}
			}
			f.Func().
				Params(
					jen.Id("m").Op("*").Id(nameStruct),
				).
				Id(method.Name).Add(types.Convert(method.Sig, f.Import)).
				BlockFunc(func(g *jen.Group) {
					constName := strcase.ToLowerCamel(iface.Named.Name + method.Name)

					labelCodes := []jen.Code{
						jen.Lit("method").Op(":").Lit(""),
						jen.Lit("code").Op(":").Lit(""),
						jen.Lit("scopeName").Op(":").Id(nameConstScopeName),
						jen.Lit("methodNameShort").Op(":").Id(constName + "Short"),
						jen.Lit("methodNameFull").Op(":").Id(constName),
					}
					errsLabelCodes := append(labelCodes, jen.Lit("errorCode").Op(":").Lit("respFailed"))

					g.Defer().
						Func().
						Params(
							jen.Id("now").Qual(timePkg, "Time"),
						).
						BlockFunc(func(g *jen.Group) {
							g.Id("m").Dot("c").Dot("Requests").Call().Dot("With").Call(
								jen.Qual(prometheusPkg, "Labels").Values(labelCodes...),
							).Dot("Inc").Call()
							if len(errorVars) > 0 {
								for _, e := range errorVars {
									g.If(jen.Id(e.Name)).Op("!=").Nil().Block(
										jen.Id("m").Dot("c").Dot("ErrRequests").Call().Dot("With").Call(
											jen.Qual(prometheusPkg, "Labels").Values(errsLabelCodes...),
										).Dot("Inc").Call(),
									)
								}
							}
							g.Id("m").Dot("c").Dot("Duration").Call().Dot("With").Call(
								jen.Qual(prometheusPkg, "Labels").Values(labelCodes...),
							).Dot("Observe").Call(jen.Qual(timePkg, "Since").Call(jen.Id("now")).Dot("Seconds").Call())
						}).Call(jen.Qual(timePkg, "Now").Call())
					if len(results) > 0 {
						g.List(results...).Op("=").Id("m").Dot("next").Dot(method.Name).Call(callParams...)
					} else {
						g.Id("m").Dot("next").Dot(method.Name).Call(callParams...)
					}
					if len(results) > 0 {
						g.Return()
					}
				})
		}
		f.Func().
			Id(nameConstructor).
			Params(
				jen.Id("c").Id("prometheusCollector"),
			).
			Do(f.Import(pkgMiddleware, nameMiddleware)).
			Block(
				jen.Return(
					jen.Func().Params(jen.Id("next").Do(f.Qual(iface.Named.Pkg.Path, iface.Named.Name))).Do(f.Qual(iface.Named.Pkg.Path, iface.Named.Name)).Block(
						jen.Return(
							jen.Op("&").Id(nameStruct).Values(
								jen.Id("next").Op(":").Id("next"),
								jen.Id("c").Op(":").Id("c"),
							),
						),
					),
				),
			)
	}
	return []file.File{f}, errs
}

func (p *Plugin) Dependencies() []string { return []string{"middleware"} }
