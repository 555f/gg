package logging

import (
	"fmt"
	"go/token"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"

	. "github.com/dave/jennifer/jen"
	"github.com/hashicorp/go-multierror"
)

type middlewarePlugin interface {
	NameMiddleware(named *types.Named) string
	Output() string
}

type Plugin struct {
	ctx *gg.Context
}

func (p *Plugin) Name() string { return "logging" }

func (p *Plugin) Exec() (files []file.File, errs error) {
	output := filepath.Join(p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("output", "internal/logging/logging.go"))
	f := file.NewGoFile(p.ctx.Module, output)
	middlewarePlugin, ok := p.ctx.Plugin("middleware").(middlewarePlugin)
	if !ok {
		errs = multierror.Append(errs, errors.Error("middleware plugin not found", token.Position{}))
		return
	}

	loggerPkg := "github.com/go-kit/log"
	levelPkg := "github.com/go-kit/log/level"
	timePkg := "time"

	f.Type().Id("errLevel").Interface(Id("Level").Params().String())
	f.Type().Id("logError").Interface(Id("LogError").Params().Error())

	f.Func().Id("levelLogger").
		Params(
			Id("e").Id("errLevel"),
			Id("logger").Qual(loggerPkg, "Logger"),
		).Qual(loggerPkg, "Logger").
		Block(
			Switch(Id("e").Dot("Level").Call()).
				Block(
					Default().Return(Qual(levelPkg, "Error").Call(Id("logger"))),
					Case(Lit("debug")).Return(Qual(levelPkg, "Debug").Call(Id("logger"))),
					Case(Lit("info")).Return(Qual(levelPkg, "Info").Call(Id("logger"))),
					Case(Lit("warn")).Return(Qual(levelPkg, "Warn").Call(Id("logger"))),
				),
		)

	for _, iface := range p.ctx.Interfaces {
		nameStruct := p.NameStruct(iface.Named)
		nameMiddleware := middlewarePlugin.NameMiddleware(iface.Named)
		pkgMiddleware := path.Dir(path.Join(p.ctx.Module.Path, strings.Replace(middlewarePlugin.Output(), p.ctx.Module.Dir, "", -1)))

		f.Type().Id(nameStruct).Struct(
			Id("next").Qual(iface.Named.Pkg.Path, iface.Named.Name),
			Id("logger").Qual(loggerPkg, "Logger"),
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
				logParams = []Code{
					Id("s").Dot("logger"),
					Lit("message"), Lit(fmt.Sprintf("call method - %s", method.Name)),
				}
				callParams []Code
				results    []Code
				logResults = []Code{Id("logger")}
				errorVar   *types.Var
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
				logParams = append(logParams, Lit(name), makeParamLog(param))
				paramNames[name]++
			}

			for _, result := range method.Sig.Results {
				results = append(results, Id(result.Name))
				if errorVar == nil && result.IsError {
					errorVar = result
					continue
				}
				opts, err := makeResultOptions(result.Tags)
				if err != nil {
					errs = multierror.Append(errs, err)
					continue
				}
				if opts.Skip {
					continue
				}
				name := result.Name
				if name == "" {
					name = opts.Name
				}
				if name == "" {
					errs = multierror.Append(errs, errors.Error("the result name cannot be empty or the logging-result-name parameter must be set", result.Position))
					continue
				}
				paramNames[name]++
				logResults = append(logResults, Lit(name), makeParamLog(result))
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
								Qual(loggerPkg, "WithPrefix").
								Call(logParams...)
							if errorVar != nil {
								g.If(Id(errorVar.Name)).Op("!=").Nil().Block(
									If(List(Id("e"), Id("ok")).
										Op(":=").
										Id(errorVar.Name).Assert(Id("errLevel")).
										Op(";").Id("ok"),
									).Block(
										Id("logger").Op("=").Id("levelLogger").Call(Id("e"), Id("logger")),
									).Else().Block(
										Id("logger").Op("=").Qual(levelPkg, "Error").Call(Id("logger")),
									).Line().
										If(List(Id("e"), Id("ok")).
											Op(":=").
											Id(errorVar.Name).Assert(Id("logError")).
											Op(";").Id("ok"),
										).Block(
										Id("logger").Op("=").Qual(loggerPkg, "WithPrefix").Call(Id("logger"), Lit(errorVar.Name), Id("e").Dot("LogError").Call()),
									).Else().Block(
										Id("logger").Op("=").Qual(loggerPkg, "WithPrefix").Call(Id("logger"), Lit(errorVar.Name), Id(errorVar.Name)),
									),
								).Else().Block(
									Id("logger").Op("=").Qual(levelPkg, "Debug").Call(Id("logger")),
									Id("logger").Op("=").Qual(loggerPkg, "WithPrefix").Call(logResults...),
								)
							} else {
								g.Id("logger").Op("=").Qual(levelPkg, "Debug").Call(Id("logger"))
							}
							g.Id("_").Op("=").Id("logger").Dot("Log").Call(Lit("dur"), Qual("time", "Since").Call(Id("now")))
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
				Id("logger").Qual(loggerPkg, "Logger"),
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
