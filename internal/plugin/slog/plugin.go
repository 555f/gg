package logging

import (
	"fmt"
	"go/token"
	"path/filepath"
	"strconv"

	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"

	. "github.com/dave/jennifer/jen"
	"github.com/hashicorp/go-multierror"
)

type middlewarePlugin interface {
	PkgPath(named *types.Named) string
	NameMiddleware(namedType *types.Named) string
	Output() string
}

type Plugin struct {
	ctx *gg.Context
}

func (p *Plugin) Name() string { return "slog" }

func (p *Plugin) Exec() (files []file.File, errs error) {
	output := filepath.Join(p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("output", "internal/logging/logging.go"))
	f := file.NewGoFile(p.ctx.Module, output)
	middlewarePlugin, ok := p.ctx.Plugin("middleware").(middlewarePlugin)
	if !ok {
		errs = multierror.Append(errs, errors.Error("middleware plugin not found", token.Position{}))
		return
	}

	loggerPkg := "golang.org/x/exp/slog"
	timePkg := "time"

	f.Type().Id("errLevel").Interface(Id("Level").Params().String())
	f.Type().Id("logError").Interface(Id("LogError").Params().Error())

	f.Func().Id("levelLogger").
		Params(
			Id("e").Id("errLevel"),
		).Qual(loggerPkg, "Level").
		Block(
			Switch(Id("e").Dot("Level").Call()).
				Block(
					Default().Return(Qual(loggerPkg, "LevelError")),
					Case(Lit("debug")).Return(Qual(loggerPkg, "LevelDebug")),
					Case(Lit("info")).Return(Qual(loggerPkg, "LevelInfo")),
					Case(Lit("warn")).Return(Qual(loggerPkg, "LevelWarn")),
				),
		)

	for _, iface := range p.ctx.Interfaces {
		nameStruct := p.NameStruct(iface.Named)
		nameMiddleware := middlewarePlugin.NameMiddleware(iface.Named)
		pkgMiddleware := middlewarePlugin.PkgPath(iface.Named)

		f.Type().Id(nameStruct).Struct(
			Id("next").Qual(iface.Named.Pkg.Path, iface.Named.Name),
			Id("logger").Op("*").Qual(loggerPkg, "Logger"),
		)

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
				logParams  []Code
				callParams []Code
				results    []Code
				// logResults []Code
				errorVars  []*types.Var
				contextVar *types.Var
				paramNames = map[string]int{}
			)

			for _, param := range method.Sig.Params {
				callParam := Id(param.Name)
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
				opts, err := makeParamOptions(method.Pkg, param.Tags)
				if err != nil {
					errs = multierror.Append(errs, err)
					continue
				}
				if opts.Skip {
					continue
				}
				name := param.Name
				if name == "" {
					name = opts.Name
				}
				if name == "" {
					errs = multierror.Append(errs, errors.Error("the parameter name cannot be empty or the logging-param-name parameter must be set", param.Position))
					continue
				}
				if logParam := makeLog(name, param.Type); logParam != nil {
					logParams = append(logParams, logParam)
					paramNames[name]++
				}
			}
			for _, result := range method.Sig.Results {
				results = append(results, Id(result.Name))
				if result.IsError {
					errorVars = append(errorVars, result)
					continue
				}
				// 	opts, err := makeResultOptions(result.Tags)
				// 	if err != nil {
				// 		errs = multierror.Append(errs, err)
				// 		continue
				// 	}
				// 	if opts.Skip {
				// 		continue
				// 	}
				// 	name := result.Name
				// 	if name == "" {
				// 		name = opts.Name
				// 	}
				// 	if name == "" {
				// 		errs = multierror.Append(errs, errors.Error("the result name cannot be empty or the logging-result-name parameter must be set", result.Position))
				// 		continue
				// 	}
				// 	paramNames[name]++
				// 	logResults = append(logResults, makeLog(name, result.Type))
			}

			if len(opts.LogContexts) > 0 && contextVar == nil {
				errs = multierror.Append(errs, errors.Error("to log a value from the context, you must declare it as a method parameter", method.Position))
				continue
			}

			for _, c := range opts.LogContexts {
				paramNames[c.LogName]++
				logParams = append(logParams, Lit(c.LogName), Id(contextVar.Name).Dot("Value").Call(Do(f.Import(c.PkgPath, c.Name))))
			}

			for name, length := range paramNames {
				if length > 1 {
					errs = multierror.Append(errs, errors.Warn(fmt.Sprintf("duplicate log name %s", strconv.Quote(name)), method.Position))
					continue
				}
			}

			f.Func().
				Params(
					Id("s").Op("*").Id(nameStruct),
				).
				Id(method.Name).Add(types.Convert(method.Sig, f.Import)).
				BlockFunc(func(g *Group) {
					g.Defer().
						Func().
						Params(
							Id("now").Qual(timePkg, "Time"),
						).
						BlockFunc(func(g *Group) {
							g.Id("logger").
								Op(":=").
								Qual(loggerPkg, "With").
								Call(logParams...)
							if len(errorVars) > 0 {
								for _, e := range errorVars {
									g.Id("logLever").Op(":=").Qual(loggerPkg, "LevelDebug")
									g.If(Id(e.Name)).Op("!=").Nil().Block(
										Id("logLever").Op("=").Qual(loggerPkg, "LevelError"),
										If(List(Id("e"), Id("ok")).
											Op(":=").
											Id(e.Name).Assert(Id("errLevel")).
											Op(";").Id("ok"),
										).Block(
											Id("logLever").Op("=").Id("levelLogger").Call(Id("e")),
										).Line().
											If(List(Id("e"), Id("ok")).
												Op(":=").
												Id(e.Name).Assert(Id("logError")).
												Op(";").Id("ok"),
											).Block(
											Id("logger").Op("=").Qual(loggerPkg, "With").Call(Lit(e.Name), Id("e").Dot("LogError").Call()),
										).Else().Block(
											Id("logger").Op("=").Qual(loggerPkg, "With").Call(Lit(e.Name), Id(e.Name)),
										),
									)
								}
							}
							g.Id("logger").Op("=").Qual(loggerPkg, "With").Call(Lit("dur"), Qual("time", "Since").Call(Id("now")))
							g.Id("logger").Dot("Log").Call(
								Qual("context", "TODO").Call(),
								Id("logLever"),
								Lit(fmt.Sprintf("call method - %s", method.Name)),
							)
						}).Call(Qual(timePkg, "Now").Call())

					if len(results) > 0 {
						g.List(results...).Op("=").Id("s").Dot("next").Dot(method.Name).Call(callParams...)
					} else {
						g.Id("s").Dot("next").Dot(method.Name).Call(callParams...)
					}
					if len(results) > 0 {
						g.Return()
					}
				})
		}

		f.Func().
			Id(p.NameConstructor(iface.Named)).
			Params(
				Id("logger").Op("*").Qual(loggerPkg, "Logger"),
			).
			Do(f.Import(pkgMiddleware, nameMiddleware)).
			Block(
				Return(
					Func().Params(Id("next").Qual(iface.Named.Pkg.Path, iface.Named.Name)).Qual(iface.Named.Pkg.Path, iface.Named.Name).Block(
						Return(Op("&").Id(nameStruct).Values(Dict{
							Id("next"):   Id("next"),
							Id("logger"): Id("logger"),
						})),
					),
				),
			)
	}
	return []file.File{f}, errs
}

func (p *Plugin) NameStruct(named *types.Named) string {
	return strcase.ToCamel(named.Name) + "LoggingMiddleware"
}

func (p *Plugin) NameConstructor(named *types.Named) string {
	return "Logging" + strcase.ToCamel(named.Name) + "Middleware"
}

func (p *Plugin) Dependencies() []string { return []string{"middleware"} }
