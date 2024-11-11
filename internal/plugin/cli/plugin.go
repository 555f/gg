package cli

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gen"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/555f/gg/pkg/typetransform"
	"github.com/dave/jennifer/jen"
)

var (
	contextPkg = "context"
	osPkg      = "os"
	flagPkg    = "flag"
	fmtPkg     = "fmt"
	timePkg    = "time"
)

type Plugin struct {
	ctx *gg.Context
}

func (p *Plugin) Name() string { return "cli" }

func (p *Plugin) Exec() (files []file.File, errs error) {
	f := file.NewGoFile(p.ctx.Module, p.Output())

	f.Func().Id("Run").ParamsFunc(func(g *jen.Group) {
		g.Id("ctx").Qual(contextPkg, "Context")
		for _, iface := range p.ctx.Interfaces {
			g.Id(strcase.ToLowerCamel(iface.Named.Name)).Do(f.Qual(iface.Named.Pkg.Path, iface.Named.Name))
		}
	}).BlockFunc(func(g *jen.Group) {
		g.Switch(jen.Qual(osPkg, "Args").Index(jen.Lit(1))).BlockFunc(func(g *jen.Group) {
			for _, iface := range p.ctx.Interfaces {
				for _, m := range iface.Type.Methods {
					commandName := strcase.ToSnake(m.Name)
					if tag, ok := m.Tags.Get("cli-name"); ok {
						commandName = tag.Value
					}
					g.Case(jen.Lit(commandName)).BlockFunc(func(g *jen.Group) {
						flagCommandName := commandName + "Cmd"

						params := make(types.Vars, 0, len(m.Sig.Params))
						for _, p := range m.Sig.Params {
							if !p.IsContext {
								params = append(params, p)
							}
						}
						if len(params) > 0 {

							for _, p := range params {
								g.Var().Id(p.Name).Add(types.Convert(p.Type, f.Import))
							}

							g.Id(flagCommandName).Op(":=").Qual(flagPkg, "NewFlagSet").Call(jen.Lit(commandName), jen.Qual(flagPkg, "ExitOnError"))

							var paramArgs types.Vars

							for _, p := range params {
								if _, ok := p.Tags.Get("cli-arg"); ok {
									paramArgs = append(paramArgs, p)
									continue
								}
								nameFlag := strcase.ToLowerCamel(p.Name)
								if tag, ok := p.Tags.Get("cli-name"); ok {
									nameFlag = tag.Value
								}

								var (
									flagFuncName string
									zeroValue    jen.Code
								)

								if basic, ok := p.Type.(*types.Basic); ok {
									switch {
									case basic.IsString():
										flagFuncName = "StringVar"
										zeroValue = jen.Lit("")
									case basic.IsUint():
										flagFuncName = "UintVar"
										zeroValue = jen.Lit(0)
									case basic.IsUint64():
										flagFuncName = "Uint64Var"
										zeroValue = jen.Lit(0)
									case basic.IsInt():
										flagFuncName = "IntVar"
										zeroValue = jen.Lit(0)
									case basic.IsInt64():
										flagFuncName = "Int64"
										zeroValue = jen.Lit(0)
									case basic.IsFloat64():
										flagFuncName = "Float64Var"
										zeroValue = jen.Lit(0)
									case basic.IsBool():
										flagFuncName = "BoolVar"
										zeroValue = jen.Lit(false)
									}
								} else if named, ok := p.Type.(*types.Named); ok && named.Name == "Duration" {
									flagFuncName = "DurationVar"
									zeroValue = jen.Lit(0)
								} else {
									continue
								}

								g.Id(flagCommandName).Dot(flagFuncName).Call(jen.Op("&").Id(p.Name), jen.Lit(nameFlag), zeroValue, jen.Lit(strings.TrimSpace(p.Title)))
							}
							g.Err().Op(":=").Id(flagCommandName).Dot("Parse").Call(jen.Qual(osPkg, "Args").Index(jen.Lit(2).Op(":")))
							g.Do(gen.CheckErr(
								jen.Qual(fmtPkg, "Println").Call(jen.Err()),
								jen.Return(),
							))
							for _, p := range paramArgs {
								t, _ := p.Tags.Get("cli-arg")
								argIndex, _ := strconv.ParseInt(t.Value, 10, 64)

								transCode, _, _ := typetransform.For(p.Type).
									SetAssignID(jen.Id(p.Name)).
									SetValueID(jen.Id(flagCommandName).Dot("Arg").Call(jen.Lit(int(argIndex)))).
									SetQualFunc(f.Import).
									SetOp("=").
									SetErrStatements(
										jen.Qual(fmtPkg, "Println").Call(jen.Err()),
										jen.Return(),
									).Parse()

								g.Add(transCode)

								// g.Add(gen.ParseValue(jen.Id(flagCommandName).Dot("Arg").Call(jen.Lit(int(argIndex))), jen.Id(p.Name), "=", p.Type, f.Import, func() jen.Code {
								// 	return jen.Do(gen.CheckErr(
								// 		jen.Qual(fmtPkg, "Println").Call(jen.Err()),
								// 		jen.Return(),
								// 	))
								// }))
							}
						}

						hasError := m.Sig.Results.HasError()

						g.Do(func(s *jen.Statement) {
							ifaceArgName := strcase.ToLowerCamel(iface.Named.Name)

							if m.Sig.Results.Len() > 0 {
								s.ListFunc(func(g *jen.Group) {
									for _, r := range m.Sig.Results {
										if !r.IsString && !r.IsError {
											g.Id("_")
											continue
										}
										g.Id(r.Name)
									}
								})
								if m.Sig.Results.Len() == 1 && hasError {
									s.Op("=")
								} else {
									s.Op(":=")
								}
							}

							s.Id(ifaceArgName).Dot(m.Name).CallFunc(func(g *jen.Group) {
								g.Id("ctx")
								for _, p := range params {
									g.Id(p.Name)
								}
							})
						})
						if hasError {
							g.Do(gen.CheckErr(
								jen.Qual(fmtPkg, "Println").Call(jen.Err()),
								jen.Return(),
							))
						}

						for _, r := range m.Sig.Results {
							if !r.IsString {
								continue
							}
							g.Qual(fmtPkg, "Println").Call(jen.Id(r.Name))
						}
					})
				}
			}
		})
	})

	return []file.File{f}, nil
}

func (p *Plugin) Output() string {
	return filepath.Join(p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("output", "internal/cli/cli.go"))
}

func (p *Plugin) Dependencies() []string { return nil }
