package metrics

import (
	"fmt"
	"go/token"
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

type middlewarePlugin interface {
	PkgPath(named *types.Named) string
	NameMiddleware(namedType *types.Named) string
	Output() string
}

type Plugin struct {
	ctx *gg.Context
}

func (p *Plugin) Name() string { return "metrics" }

func (p *Plugin) Exec() (files []file.File, errs error) {
	output := filepath.Join(p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("output", "internal/metrics/metrics.go"))

	f := file.NewGoFile(p.ctx.Module, output)
	middlewarePlugin, ok := p.ctx.Plugin("middleware").(middlewarePlugin)
	if !ok {
		errs = multierror.Append(errs, errors.Error("middleware plugin not found", token.Position{}))
		return
	}

	for _, iface := range p.ctx.Interfaces {
		for _, m := range iface.Type.Methods {
			constName := strcase.ToLowerCamel(iface.Named.Name + m.Name)
			f.Const().Id(constName).Op("=").Lit(m.FullName)
			f.Const().Id(constName + "Short").Op("=").Lit(shortMethodName(m))
		}
	}

	for _, iface := range p.ctx.Interfaces {
		nameStruct := p.NameStruct(iface.Named)
		nameMiddleware := middlewarePlugin.NameMiddleware(iface.Named)
		pkgMiddleware := middlewarePlugin.PkgPath(iface.Named)

		f.Type().Id(nameStruct).StructFunc(func(g *jen.Group) {
			g.Id("next").Qual(iface.Named.Pkg.Path, iface.Named.Name)
			g.Id("inRequests").Op("*").Qual(prometheusPkg, "CounterVec")
			g.Id("requests").Op("*").Qual(prometheusPkg, "CounterVec")
			g.Id("errRequests").Op("*").Qual(prometheusPkg, "CounterVec")
			g.Id("duration").Op("*").Qual(prometheusPkg, "HistogramVec")
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
						jen.Lit("method").Op(":").Id(constName),
						jen.Lit("shortMethod").Op(":").Id(constName + "Short"),
					}

					g.Id("m").Dot("inRequests").Dot("With").Call(
						jen.Qual(prometheusPkg, "Labels").Values(labelCodes...),
					).Dot("Inc").Call()

					g.Defer().
						Func().
						Params(
							jen.Id("now").Qual(timePkg, "Time"),
						).
						BlockFunc(func(g *jen.Group) {
							g.Id("m").Dot("requests").Dot("With").Call(
								jen.Qual(prometheusPkg, "Labels").Values(labelCodes...),
							).Dot("Inc").Call()
							if len(errorVars) > 0 {
								for _, e := range errorVars {
									g.If(jen.Id(e.Name)).Op("!=").Nil().Block(
										jen.Id("m").Dot("errRequests").Dot("With").Call(
											jen.Qual(prometheusPkg, "Labels").Values(labelCodes...),
										).Dot("Inc").Call(),
									)
								}
							}
							g.Id("m").Dot("duration").Dot("With").Call(
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
			Id(p.NameConstructor(iface.Named)).
			Params(
				jen.Id("namespace").String(),
				jen.Id("subsystem").String(),
			).
			Do(f.Import(pkgMiddleware, nameMiddleware)).
			Block(
				jen.Return(
					jen.Func().Params(jen.Id("next").Qual(iface.Named.Pkg.Path, iface.Named.Name)).Qual(iface.Named.Pkg.Path, iface.Named.Name).Block(
						jen.Return(jen.Op("&").Id(nameStruct).Values(
							jen.Id("next").Op(":").Id("next"),
							jen.Id("inRequests").Op(":").Qual(prometheusPkg, "NewCounterVec").CallFunc(func(g *jen.Group) {
								g.Qual(prometheusPkg, "CounterOpts").Values(
									jen.Id("Namespace").Op(":").Id("namespace"),
									jen.Id("Subsystem").Op(":").Id("subsystem"),
									jen.Id("Name").Op(":").Lit("in_requests_total"),
									jen.Id("Help").Op(":").Lit("A counter for incoming requests."),
								)
								g.Index().String().Values(jen.Lit("method"), jen.Lit("shortMethod"))
							}),
							jen.Id("requests").Op(":").Qual(prometheusPkg, "NewCounterVec").CallFunc(func(g *jen.Group) {
								g.Qual(prometheusPkg, "CounterOpts").Values(
									jen.Id("Namespace").Op(":").Id("namespace"),
									jen.Id("Subsystem").Op(":").Id("subsystem"),
									jen.Id("Name").Op(":").Lit("requests_total"),
									jen.Id("Help").Op(":").Lit("A counter for complete requests."),
								)
								g.Index().String().Values(jen.Lit("method"), jen.Lit("shortMethod"))
							}),
							jen.Id("errRequests").Op(":").Qual(prometheusPkg, "NewCounterVec").CallFunc(func(g *jen.Group) {
								g.Qual(prometheusPkg, "CounterOpts").Values(
									jen.Id("Namespace").Op(":").Id("namespace"),
									jen.Id("Subsystem").Op(":").Id("subsystem"),
									jen.Id("Name").Op(":").Lit("err_requests_total"),
									jen.Id("Help").Op(":").Lit("A counter for error requests."),
								)
								g.Index().String().Values(jen.Lit("method"), jen.Lit("shortMethod"))
							}),
							jen.Id("duration").Op(":").Qual(prometheusPkg, "NewHistogramVec").CallFunc(func(g *jen.Group) {
								g.Qual(prometheusPkg, "HistogramOpts").Values(
									jen.Id("Namespace").Op(":").Id("namespace"),
									jen.Id("Subsystem").Op(":").Id("subsystem"),
									jen.Id("Name").Op(":").Lit("request_duration_histogram_seconds"),
									jen.Id("Help").Op(":").Lit("A histogram of outgoing request latencies."),
								)
								g.Index().String().Values(jen.Lit("method"), jen.Lit("shortMethod"))
							}),
						)),
					),
				),
			)
	}
	return []file.File{f}, errs
}

func (p *Plugin) NameStruct(named *types.Named) string {
	return strcase.ToCamel(named.Name) + "MetricMiddleware"
}

func (p *Plugin) NameConstructor(named *types.Named) string {
	return "Logging" + strcase.ToCamel(named.Name) + "Middleware"
}

func (p *Plugin) Dependencies() []string { return []string{"middleware"} }
