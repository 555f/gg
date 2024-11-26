package clienttest

import (
	"fmt"
	"strings"

	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/gen"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/555f/gg/pkg/typetransform"
	"github.com/dave/jennifer/jen"
	"github.com/jaswdr/faker/v2"
)

const (
	promCollectorName = "prometheusCollector"
	prometheusPkg     = "github.com/prometheus/client_golang/prometheus"
	jsonPkg           = "encoding/json"
)

type Config struct {
	StatusCode int
	CheckError bool
}

type ClientTestGenerator struct {
	group        *jen.Group
	pkgPath      string
	fake         faker.Faker
	qualFn       types.QualFunc
	errorWrapper *options.ErrorWrapper
}

func (g *ClientTestGenerator) basicTypeToValue(t *types.Basic, manualValue ...string) jen.Code {
	switch {
	default:
		return jen.Lit(g.fake.Lorem().Sentence(10))
	case t.IsBool():
		return jen.Lit(true)
	case t.IsInteger():
		return jen.Lit(g.fake.RandomNumber(5))
	case t.IsFloat():
		return jen.Lit(g.fake.Float64(2, 1, 100))
	case t.IsSigned():
		return jen.Lit(10)
	}
}

func (g *ClientTestGenerator) typeToValue(t any, manualValue ...string) jen.Code {
	switch u := t.(type) {
	case *types.Basic:
		c := g.basicTypeToValue(u, manualValue...)
		if u.IsPointer {
			c = jen.Id("ptr").Call(c)
		}
		return c
	case *types.Named:
		var s jen.Statement
		if u.IsPointer {
			s.Op("&")
		}
		if b := u.Basic(); b != nil {
			var value any
			switch {
			case b.IsString():
				if len(manualValue) > 0 && manualValue[0] != "" {
					value = manualValue[0]
				} else {
					value = "abc"
				}
			}
			return s.Do(g.qualFn(u.Pkg.Path, u.Name)).Call(jen.Lit(value))
		}
		return s.Do(g.qualFn(u.Pkg.Path, u.Name)).Values()
	case *types.Map:
		if _, ok := u.Key.(*types.Basic); ok {
			return types.Convert(t, g.qualFn).Values(
				jen.Add(g.typeToValue(u.Key)).Op(":").Add(g.typeToValue(u.Value)),
			)
		}
		return types.Convert(t, g.qualFn).Values()
	case *types.Slice:
		return jen.Index().Add(types.Convert(u.Value, g.qualFn)).ValuesFunc(func(group *jen.Group) {
			switch v := u.Value.(type) {
			case *types.Named:
				if s := v.Struct(); s != nil {
					group.ValuesFunc(func(group *jen.Group) {
						for _, f := range s.Fields {
							group.Id(f.Name).Op(":").Add(g.typeToValue(f.Type))
						}
					})
				}
			}
		})
	case *types.Array:
		return jen.Index(jen.Lit(u.Len)).Add(types.Convert(u.Value, g.qualFn)).Values()
	case *types.Struct:
		return types.Convert(t, g.qualFn).Values()
	case *types.Interface:
		return jen.Nil()
	}

	panic(fmt.Sprintf("unreachable %T", t))
}

func (g *ClientTestGenerator) serverResponseGenerate(group *jen.Group, cfg Config, ep options.Endpoint) {
	if !cfg.CheckError && len(ep.BodyResults) > 0 {
		if ep.NoWrapResponse && len(ep.BodyResults) == 1 {
			serverResponse := ep.BodyResults[0]

			group.Var().Id("serverResponse").Add(types.Convert(serverResponse.Type, g.qualFn))

			if named, ok := serverResponse.Type.(*types.Named); ok {
				for _, f := range named.Struct().Fields {
					for _, v := range gen.Flatten(f) {
						var manualValue string
						if t, ok := v.Var.Tags.Get("http-test-value"); ok {
							manualValue = t.Value
						}
						group.Id("serverResponse").Op(".").Add(v.Path).Op("=").Add(g.typeToValue(v.Var.Type, manualValue))
					}
				}
			}
		}
	}
}

func (g *ClientTestGenerator) bodyParamsGenerate(group *jen.Group, ep options.Endpoint) {
	if len(ep.BodyParams) > 0 {
		group.Var().Id("serverRequest").StructFunc(func(group *jen.Group) {
			for _, p := range ep.BodyParams {
				group.Id(p.FldName.Camel()).Add(types.Convert(p.Type, g.qualFn))
			}
		})
		for _, p := range ep.BodyParams {
			if named, ok := p.Type.(*types.Named); ok {
				if st := named.Struct(); st != nil {
					for _, f := range st.Fields {
						for _, v := range gen.Flatten(f) {
							if v.IsArray {
								continue
							}
							group.Id("serverRequest").Dot(p.FldName.Camel()).Op(".").Add(v.Path).Op("=").Add(g.typeToValue(v.Var.Type))
						}
					}
				}
			} else {
				group.Id("serverRequest").Dot(p.FldName.Camel()).Op("=").Add(g.typeToValue(p.Type))
			}
		}
	}
}

func (g *ClientTestGenerator) paramsGenerate(group *jen.Group, params options.EndpointParams) {
	for _, p := range params {
		if p.HTTPType == options.BodyHTTPType {
			continue
		}
		postfix := strcase.ToCamel(string(p.HTTPType))
		group.Var().Id(p.FldName.LowerCamel() + postfix).Add(types.Convert(p.Type, g.qualFn)).Op("=").Add(g.typeToValue(p.Type))
	}
}

func (g ClientTestGenerator) mockServerGenerate(group *jen.Group, ep options.Endpoint, errorWrapperName string, cfg Config) {
	pathParts := strings.Split(ep.Path, "/")
	for i, part := range pathParts {
		if strings.HasPrefix(part, ":") {
			pathParts[i] = "{" + part[1:] + "}"
		}
	}

	group.Id("mockServer").Op(":=").Qual("net/http", "NewServeMux").Call()
	group.Id("mockServer").Dot("Handle").Call(
		jen.Lit(strings.Join(pathParts, "/")),
		jen.Qual("net/http", "HandlerFunc").Call(
			jen.Func().Params(
				jen.Id("w").Qual("net/http", "ResponseWriter"),
				jen.Id("r").Op("*").Qual("net/http", "Request"),
			).BlockFunc(func(group *jen.Group) {
				if len(ep.BodyParams) > 0 {
					group.Var().Id("body").StructFunc(func(group *jen.Group) {
						for _, p := range ep.BodyParams {
							group.Id(p.FldName.Camel()).Add(types.Convert(p.Type, g.qualFn)).Tag(map[string]string{
								"json": p.Name,
							})
						}
					})

					var bodyVar jen.Code
					if ep.NoWrapRequest {
						bodyVar = jen.Op("&").Id("body").Dot(ep.BodyParams[0].FldName.Camel())
					} else {
						bodyVar = jen.Op("&").Id("body")
					}

					group.Id("_").Op("=").Qual(jsonPkg, "NewDecoder").Call(jen.Id("r").Dot("Body")).Dot("Decode").Call(bodyVar)

					for _, p := range ep.BodyParams {
						switch t := p.Type.(type) {
						default:
							group.If(jen.Id("body").Dot(p.FldName.Camel()).Op("!=").Id("serverRequest").Dot(p.FldName.Camel()).BlockFunc(func(g *jen.Group) {
								g.Id("t").Dot("Fatal").Call(jen.Lit("failed equal method " + ep.MethodShortName + " " + p.Name))
							}))
						case *types.Named:
							if st := t.Struct(); st != nil {
								for _, f := range st.Fields {
									for _, v := range gen.Flatten(f) {
										if v.IsArray {
											continue
										}
										fieldPath := v.Paths.String()

										switch v.Var.Type.(type) {
										default:
											group.If(jen.Id("body").Dot(p.FldName.Camel()).Op(".").Add(v.Path).Op("!=").Id("serverRequest").Dot(p.FldName.Camel()).Op(".").Add(v.Path)).BlockFunc(func(g *jen.Group) {
												g.Id("t").Dot("Fatal").Call(jen.Lit("failed equal method " + ep.MethodShortName + " " + fieldPath))
											})
										case *types.Map:
											group.If(jen.Op("!").Qual("reflect", "DeepEqual").Call(jen.Id("body").Dot(p.FldName.Camel()), jen.Id("serverRequest").Dot(p.FldName.Camel())).BlockFunc(func(g *jen.Group) {
												g.Id("t").Dot("Fatal").Call(jen.Lit("failed equal method " + ep.MethodShortName + " " + p.Name))
											}))
										}
									}
								}
							}
						case *types.Slice:
							switch t := t.Value.(type) {
							case *types.Basic:
								group.If(jen.Id("body").Dot(p.FldName.Camel()).Index(jen.Lit(0)).Op("!=").Id("serverRequest").Dot(p.FldName.Camel()).Index(jen.Lit(0)).BlockFunc(func(g *jen.Group) {
									g.Id("t").Dot("Fatal").Call(jen.Lit("failed equal method " + ep.MethodShortName + " " + p.Name))
								}))
							case *types.Named:
								if st := t.Struct(); st != nil {
									for _, f := range st.Fields {
										for _, v := range gen.Flatten(f) {
											if v.IsArray {
												continue
											}
											fieldPath := v.Paths.String()

											if named, ok := v.Var.Type.(*types.Named); ok && named.System() {
												group.If(jen.Op("!").Qual("reflect", "DeepEqual").Call(jen.Id("body").Dot(p.FldName.Camel()).Index(jen.Lit(0)).Op(".").Add(v.Path), jen.Id("serverRequest").Dot(p.FldName.Camel()).Index(jen.Lit(0)).Op(".").Add(v.Path))).BlockFunc(func(g *jen.Group) {
													g.Id("t").Dot("Fatal").Call(jen.Lit("failed equal method " + ep.MethodShortName + " " + fieldPath))
												})
											} else {
												group.If(jen.Id("body").Dot(p.FldName.Camel()).Index(jen.Lit(0)).Op(".").Add(v.Path).Op("!=").Id("serverRequest").Dot(p.FldName.Camel()).Index(jen.Lit(0)).Op(".").Add(v.Path)).BlockFunc(func(g *jen.Group) {
													g.Id("t").Dot("Fatal").Call(jen.Lit("failed equal method " + ep.MethodShortName + " " + fieldPath))
												})
											}
										}
									}
								}
							}
						}

						// if named, ok := p.Type.(*types.Named); ok {
						// 	if st := named.Struct(); st != nil {
						// 		for _, f := range st.Fields {
						// 			for _, v := range gen.Flatten(f) {
						// 				if v.IsArray {
						// 					continue
						// 				}
						// 				fieldPath := v.Paths.String()

						// 				group.If(jen.Id("body").Dot(p.FldName.Camel()).Op(".").Add(v.Path).Op("!=").Id("serverRequest").Dot(p.FldName.Camel()).Op(".").Add(v.Path)).BlockFunc(func(g *jen.Group) {
						// 					g.Id("t").Dot("Fatal").Call(jen.Lit("failed equal method " + ep.MethodShortName + " " + fieldPath))
						// 				})
						// 			}
						// 		}
						// 	}
						// } else {
						// 	if a, ok := p.Type.(*types.Slice); ok {

						// 	} else {

						// 	}
						// }
					}
				}

				if len(ep.QueryParams) > 0 {
					group.Id("q").Op(":=").Id("r").Dot("URL").Dot("Query").Call()
				}

				for _, p := range ep.Params {
					if p.HTTPType == options.BodyHTTPType {
						continue
					}
					switch p.HTTPType {
					case options.PathHTTPType:
						transCode, paramID, _ := typetransform.For(p.Type).
							SetAssignID(jen.Id(p.FldName.LowerCamel())).
							SetValueID(jen.Id("r").Dot("PathValue").Call(jen.Lit(p.Name))).
							SetOp(":=").
							SetQualFunc(g.qualFn).
							SetErrStatements(
								jen.Id("t").Dot("Fatal").Call(jen.Err()),
							).Parse()
						if transCode != nil {
							group.Add(transCode)
						}
						group.If(jen.Id(p.FldName.LowerCamel() + "Path").Op("!=").Add(paramID).BlockFunc(func(g *jen.Group) {
							g.Id("t").Dot("Fatal").Call(jen.Lit("failed equal method " + ep.MethodShortName + " " + p.FldName.LowerCamel()))
						}))
					case options.QueryHTTPType:
						transCode, paramID, _ := typetransform.For(p.Type).
							SetAssignID(jen.Id(p.FldName.LowerCamel())).
							SetValueID(jen.Id("q").Dot("Get").Call(jen.Lit(p.Name))).
							SetOp(":=").
							SetQualFunc(g.qualFn).
							SetErrStatements(
								jen.Id("t").Dot("Fatal").Call(jen.Err()),
							).Parse()
						if transCode != nil {
							group.Add(transCode)
						}
						group.If(jen.Id(p.FldName.LowerCamel() + "Query").Op("!=").Add(paramID).BlockFunc(func(group *jen.Group) {
							group.Id("t").Dot("Fatal").Call(jen.Lit("failed equal method " + ep.MethodShortName + " " + p.FldName.LowerCamel()))
						}))
					case options.HeaderHTTPType:
						transCode, paramID, _ := typetransform.For(p.Type).
							SetAssignID(jen.Id(p.FldName.LowerCamel())).
							SetValueID(jen.Id("r").Dot("Header").Dot("Get").Call(jen.Lit(p.Name))).
							SetOp(":=").
							SetQualFunc(g.qualFn).
							SetErrStatements(
								jen.Id("t").Dot("Fatal").Call(jen.Err()),
							).Parse()

						if transCode != nil {
							group.Add(transCode)
						}
						group.If(jen.Id(p.FldName.LowerCamel() + "Header").Op("!=").Add(paramID).BlockFunc(func(g *jen.Group) {
							g.Id("t").Dot("Fatal").Call(jen.Lit("failed equal method " + ep.MethodShortName + " " + p.FldName.LowerCamel()))
						}))
					}
				}

				// for _, p := range ep.CookieParams {
				// 	g.If(jen.Id(p.FldName.LowerCamel() + "Cookie").Op("!=").Id("r").Dot("PathValue").Call(jen.Lit(p.Name)).BlockFunc(func(g *jen.Group) {
				// 		g.Id("t").Dot("Fatal").Call(jen.Lit("failed equal method " + ep.MethodShortName + " " + p.Name))
				// 	}))
				// }

				if cfg.StatusCode != 0 {
					group.Id("w").Dot("WriteHeader").Call(jen.Lit(cfg.StatusCode))
				}
				if !cfg.CheckError && len(ep.BodyResults) > 0 {
					group.List(jen.Id("data"), jen.Id("_")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Id("serverResponse"))
					group.Id("w").Dot("Write").Call(jen.Id("data"))
				}

				if errorWrapperName != "" {
					group.List(jen.Id("data"), jen.Id("_")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Id(errorWrapperName))
					group.Id("w").Dot("Write").Call(jen.Id("data"))
				}
			}),
		),
	)
}

func (g *ClientTestGenerator) generateCheckError(group *jen.Group, ep options.Endpoint, cfg Config, errorWrapperName string) {
	if !ep.Sig.Results.HasError() {
		return
	}
	if !cfg.CheckError {
		group.Do(gen.CheckErr(
			jen.Id("t").Dot("Fatalf").Call(jen.Lit("%s: %s"), jen.Lit("failed execute method "+ep.MethodShortName), jen.Id("err")),
		))
		return
	}
	group.Do(gen.CheckNotErr(
		jen.Id("t").Dot("Fatal").Call(jen.Lit("failed execute method " + ep.MethodShortName + " error is nil")),
	))
	if g.errorWrapper != nil {
		group.Var().Id("e").Op("*").Do(g.qualFn(g.errorWrapper.Default.Named.Pkg.Path, g.errorWrapper.Default.Named.Name))
		group.If(jen.Do(g.qualFn("errors", "As")).Call(jen.Err(), jen.Op("&").Id("e"))).BlockFunc(func(group *jen.Group) {
			for _, f := range g.errorWrapper.Fields {
				group.If(
					jen.Id(errorWrapperName).Dot(f.FldName).Op("!=").Id("e").Dot(f.FldName),
				).Block(jen.Id("t").Dot("Fatal").Call(jen.Lit("failed equal error field " + ep.MethodShortName + " " + f.FldName + " not equal")))
			}
		}).Else().Block(
			jen.Id("t").Dot("Fatal").Call(jen.Lit("failed equal error " + ep.MethodShortName + " not equal")),
		)
	}
}

func (g *ClientTestGenerator) generateCheckBodyResult(group *jen.Group, ep options.Endpoint, cfg Config) {
	if cfg.CheckError || len(ep.BodyResults) <= 0 {
		return
	}
	for _, r := range ep.BodyResults {
		if named, ok := r.Type.(*types.Named); ok {
			st := named.Struct()
			if st == nil {
				continue
			}
			for _, f := range st.Fields {
				for _, v := range gen.Flatten(f) {
					if v.IsArray {
						continue
					}

					fieldPath := v.Paths.String()

					if v.Var.IsPointer {
						group.If(jen.Id(r.Name).Op(".").Add(v.Path).Op("==").Nil()).BlockFunc(func(group *jen.Group) {
							group.Id("t").Dot("Fatal").Call(jen.Lit("failed equal method " + ep.MethodShortName + " " + fieldPath + " is nil"))
						})
						if _, ok := v.Var.Type.(*types.Basic); ok {
							group.If(jen.Op("*").Id(r.Name).Op(".").Add(v.Path).Op("!=").Op("*").Id("serverResponse").Op(".").Add(v.Path)).BlockFunc(func(group *jen.Group) {
								group.Id("t").Dot("Fatal").Call(jen.Lit("failed equal method " + ep.MethodShortName + " " + fieldPath + " not equal"))
							})
						}
					} else {
						group.If(
							jen.Id(r.Name).Op(".").Add(v.Path).Op("!=").Id("serverResponse").Op(".").Add(v.Path),
						).BlockFunc(func(group *jen.Group) {
							group.Id("t").Dot("Fatal").Call(jen.Lit("failed equal method " + ep.MethodShortName + " " + fieldPath + " not equal"))
						})
					}
				}
			}
		}

	}
}

func (g *ClientTestGenerator) generateErrorWrapper(group *jen.Group, cfg Config) (errorWrapperName string) {
	if !cfg.CheckError {
		return
	}
	if g.errorWrapper == nil {
		return
	}
	errorWrapperName = strcase.ToLowerCamel(g.errorWrapper.Struct.Named.Name)
	group.Id(errorWrapperName).Op(":=").
		Do(g.qualFn(g.errorWrapper.Struct.Named.Pkg.Path, g.errorWrapper.Struct.Named.Name)).
		ValuesFunc(func(group *jen.Group) {
			for _, field := range g.errorWrapper.Fields {
				group.Id(field.FldName).Op(":").Add(g.typeToValue(field.FldType))
			}
		})
	return errorWrapperName
}

func (g *ClientTestGenerator) Generate(iface options.Iface, ep options.Endpoint, configs []Config) {
	constructName := "create" + iface.Name + "Client"

	for _, cfg := range configs {
		testMethod := fmt.Sprintf("%s_%d", ep.MethodName, cfg.StatusCode)
		testName := "Test" + iface.Name + "_" + testMethod

		g.group.Func().Id(testName).Params(jen.Id("t").Op("*").Qual("testing", "T")).BlockFunc(func(group *jen.Group) {
			g.serverResponseGenerate(group, cfg, ep)
			g.bodyParamsGenerate(group, ep)
			g.paramsGenerate(group, ep.Params)

			errorWrapperName := g.generateErrorWrapper(group, cfg)

			g.mockServerGenerate(group, ep, errorWrapperName, cfg)

			group.Id("server").Op(":=").Qual("net/http/httptest", "NewServer").Call(
				jen.Id("mockServer"),
			)

			opts := []jen.Code{
				jen.Id("server").Dot("URL"),
				jen.Lit(ep.MethodName),
				jen.Lit(cfg.StatusCode),
			}

			group.Id("client").Op(":=").Id(constructName).Call(opts...)

			group.Do(func(s *jen.Statement) {
				if ep.Sig.Results.Len() > 0 {
					s.ListFunc(func(g *jen.Group) {
						for _, r := range ep.Sig.Results {
							if r.IsError {
								g.Id(r.Name)
								continue
							}
							if !cfg.CheckError {
								g.Id(r.Name)
							} else {
								g.Id("_")
							}

						}
					})
					if cfg.CheckError || ep.Sig.Results.HasError() {
						s.Op(":=")
					} else {
						s.Op("=")
					}
				}
			}).Id("client").Dot(ep.MethodName).CallFunc(func(g *jen.Group) {
				if ep.Context != nil {
					g.Qual("context", "TODO").Call()
				}
				for _, p := range ep.Params {
					switch p.HTTPType {
					case options.HeaderHTTPType:
						g.Id(p.FldName.LowerCamel() + "Header")
					case options.CookieHTTPType:
						g.Id(p.FldName.LowerCamel() + "Cookie")
					case options.QueryHTTPType:
						g.Id(p.FldName.LowerCamel() + "Query")
					case options.BodyHTTPType:
						g.Id("serverRequest").Dot(p.FldName.Camel())
					case options.PathHTTPType:
						g.Id(p.FldName.LowerCamel() + "Path")
					}
				}
			})

			g.generateCheckError(group, ep, cfg, errorWrapperName)

			g.generateCheckBodyResult(group, ep, cfg)
		})
	}
}

func New(group *jen.Group, pkgPath string, fake faker.Faker, qualFn types.QualFunc, errorWrapper *options.ErrorWrapper) *ClientTestGenerator {
	group.Func().Id("ptr").Types(jen.Id("T").Any()).Params(jen.Id("t").Id("T")).Op("*").Id("T").Block(
		jen.Return(jen.Op("&").Id("t")),
	)

	return &ClientTestGenerator{
		group:        group,
		pkgPath:      pkgPath,
		fake:         fake,
		qualFn:       qualFn,
		errorWrapper: errorWrapper,
	}
}
