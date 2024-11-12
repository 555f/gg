package rest

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/gen"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/555f/gg/pkg/typetransform"
	"github.com/dave/jennifer/jen"
)

type clientEndpointBuilder struct {
	*BaseClientBuilder
	iface     options.Iface
	ep        options.Endpoint
	qualifier Qualifier
}

func (b *clientEndpointBuilder) BuildReqStruct() ClientEndpointBuilder {
	methodRequestName := b.methodRequestName()
	clientName := b.clientName()

	b.codes = append(b.codes,
		jen.Type().Id(methodRequestName).StructFunc(func(g *jen.Group) {
			g.Id("c").Op("*").Id(clientName)
			g.Id("client").Op("*").Qual(httpPkg, "Client")
			g.Id("opts").Op("*").Id(clientOptionName)
			g.Id("params").StructFunc(func(g *jen.Group) {
				if len(b.ep.Params) > 0 {
					for _, param := range b.ep.Params {
						if len(param.Params) > 0 {
							for _, childParam := range param.Params {
								g.Add(b.makeRequestStructParam(param, childParam, b.qualifier.Qual))
							}
							continue
						}
						g.Add(b.makeRequestStructParam(nil, param, b.qualifier.Qual))
					}
				}
			})
		}))
	return b
}

func (b *clientEndpointBuilder) BuildSetters() ClientEndpointBuilder {
	for _, param := range b.ep.Params {
		if len(param.Params) > 0 {
			for _, childParam := range param.Params {
				b.buildSetter(param, childParam)
			}
		} else {
			if !param.Required {
				b.buildSetter(nil, param)
			}
		}
	}
	return b
}

func (b *clientEndpointBuilder) BuildReqMethod() ClientEndpointBuilder {
	methodRequestName := b.methodRequestName()
	recvName := b.recvName()
	clientName := b.clientName()
	methodReqName := b.methodReqName()

	b.codes = append(b.codes, jen.Func().Params(jen.Id(recvName).Op("*").Id(clientName)).Id(methodReqName).
		ParamsFunc(func(g *jen.Group) {
			for _, param := range b.ep.Params {
				if param.Required {
					g.Id(param.FldName.LowerCamel()).Add(types.Convert(param.Type, b.qualifier.Qual))
				}
			}
		}).
		Op("*").Id(methodRequestName).BlockFunc(func(g *jen.Group) {
		g.Id("m").Op(":=").Op("&").Id(methodRequestName).Values(
			jen.Id("client").Op(":").Id(recvName).Dot("opts").Dot("client"),
			jen.Id("opts").Op(":").Op("&").Id(clientOptionName).Values(
				jen.Id("ctx").Op(":").Qual("context", "TODO").Call(),
			),
			jen.Id("c").Op(":").Id(recvName),
		)
		for _, param := range b.ep.Params {
			if param.Required {
				g.Id("m").Dot("params").Dot(param.FldName.LowerCamel()).Op("=").Id(param.FldName.LowerCamel())
			}
		}
		g.Return(jen.Id("m"))
	}))
	return b
}

func (b *clientEndpointBuilder) BuildMethod() ClientEndpointBuilder {
	recvName := b.recvName()
	clientName := b.clientName()
	methodReqName := b.methodReqName()

	b.codes = append(b.codes, jen.Func().Params(jen.Id(recvName).Op("*").Id(clientName)).Id(b.ep.MethodName).
		ParamsFunc(func(g *jen.Group) {
			for _, param := range b.ep.Sig.Params {
				g.Id(param.Name).Add(types.Convert(param.Type, b.qualifier.Qual))
			}
		}).
		ParamsFunc(func(g *jen.Group) {
			for _, result := range b.ep.Results {
				g.Id(result.FldName.LowerCamel()).Add(types.Convert(result.Type, b.qualifier.Qual))
			}
			g.Err().Error()
		}).
		BlockFunc(func(g *jen.Group) {
			g.ListFunc(func(g *jen.Group) {
				for _, param := range b.ep.Results {
					g.Id(param.FldName.LowerCamel())
				}
				g.Err()
			}).Op("=").Id(recvName).Dot(methodReqName).CallFunc(func(g *jen.Group) {
				for _, param := range b.ep.Params {
					if param.Required {
						g.Id(param.FldName.LowerCamel())
					}
				}
			}).CustomFunc(jen.Options{}, func(g *jen.Group) {
				buildSetters := func(params options.EndpointParams) {
					for _, param := range params {
						if param.Required {
							continue
						}
						methodSetName := param.FldName.Camel()
						fldName := jen.Id(param.FldName.LowerCamel())
						if param.Parent != nil {
							methodSetName = param.Parent.FldName.Camel() + param.FldName.Camel()
							fldName = jen.Id(param.Parent.FldName.LowerCamel()).Dot(param.FldName.String())
						}
						g.Dot("Set" + methodSetName).Call(fldName)
					}
				}

				buildSetters(b.ep.BodyParams)
				buildSetters(b.ep.QueryParams)
				buildSetters(b.ep.HeaderParams)
				buildSetters(b.ep.CookieParams)

			}).Dot("Execute").CallFunc(func(g *jen.Group) {
				if b.ep.Context != nil {
					g.Id("WithContext").Call(jen.Id(b.ep.Context.Name))
				}
			})
			g.Return()
		}))
	return b
}

func (b *clientEndpointBuilder) BuildExecuteMethod() ClientEndpointBuilder {
	methodRequestName := b.methodRequestName()
	recvName := b.recvName()

	makeParam := func(name jen.Code, t any, required bool, f func(v jen.Code) jen.Code) jen.Code {
		paramID := jen.Id(recvName).Dot("params").Add(name)

		named, isNamed := t.(*types.Named)

		var code jen.Code
		if required {
			code = f(paramID)
		} else {
			if isNamed && named.Pkg.Path == nullPkg {
				code = jen.If(jen.Add(paramID).Dot("Valid")).Block(f(paramID))
			} else {
				if isNamed {
					code = jen.If(jen.Add(paramID).Op("!=").Nil()).Block(f(jen.Call(jen.Op("*").Add(paramID))))
				} else {
					code = jen.If(jen.Add(paramID).Op("!=").Nil()).Block(f(jen.Op("*").Add(paramID)))
				}
			}
		}
		return code
	}

	b.codes = append(b.codes, jen.Func().Params(jen.Id(recvName).Op("*").Id(methodRequestName)).Id("Execute").
		Params(
			jen.Id("opts").Op("...").Id("ClientOption"),
		).
		ParamsFunc(func(g *jen.Group) {
			for _, result := range b.ep.Results {
				g.Id(result.FldName.LowerCamel()).Add(types.Convert(result.Type, b.qualifier.Qual))
			}
			g.Err().Error()
		}).
		Block(
			jen.For(jen.List(jen.Id("_"), jen.Id("o")).Op(":=").Range().Id("opts")).Block(
				jen.Id("o").Call(jen.Id(recvName).Dot("opts")),
			),
			jen.List(jen.Id("ctx"), jen.Id("cancel")).Op(":=").Qual("context", "WithCancel").Call(jen.Id(recvName).Dot("opts").Dot("ctx")),
			jen.Do(func(s *jen.Statement) {
				if len(b.ep.PathParams) > 0 {
					var paramsCall []jen.Code
					paramsCall = append(paramsCall, jen.Lit(b.ep.SprintfPath()))
					for _, name := range b.ep.ParamsNameIdx {
						paramsCall = append(paramsCall, jen.Id(recvName).Dot("params").Dot(name))
					}
					s.Id("path").Op(":=").Qual("fmt", "Sprintf").Call(paramsCall...)
				} else {
					s.Id("path").Op(":=").Lit(b.ep.Path)
				}
			}),

			jen.Id("r").Dot("opts").Dot("ctx").Op("=").Qual(ctxPkg, "WithValue").Call(
				jen.Id("r").Dot("opts").Dot("ctx"),
				jen.Id("methodContextKey"),
				jen.Id(strcase.ToLowerCamel(b.ep.MethodName)+"FullName"),
			),
			jen.Id("r").Dot("opts").Dot("ctx").Op("=").Qual(ctxPkg, "WithValue").Call(
				jen.Id("r").Dot("opts").Dot("ctx"),
				jen.Id("shortMethodContextKey"),
				jen.Id(strcase.ToLowerCamel(b.ep.MethodName)+"ShortName"),
			),
			jen.Id("r").Dot("opts").Dot("ctx").Op("=").Qual(ctxPkg, "WithValue").Call(
				jen.Id("r").Dot("opts").Dot("ctx"),
				jen.Id("scopeNameContextKey"),
				jen.Id(strcase.ToLowerCamel(b.iface.Name)+"ScopeName"),
			),

			jen.List(jen.Id("req"), jen.Err()).Op(":=").Qual(httpPkg, "NewRequestWithContext").
				Call(
					jen.Id("r").Dot("opts").Dot("ctx"),
					jen.Lit(b.ep.HTTPMethod),
					jen.Id(recvName).Dot("c").Dot("target").Op("+").Id("path"), jen.Nil(),
				),
			jen.Do(gen.CheckErr(
				jen.Id("cancel").Call(),
				jen.Return(),
			)),
			jen.CustomFunc(jen.Options{Multi: true}, func(g *jen.Group) {
				if len(b.ep.BodyParams) > 0 {
					switch b.ep.ClientContentType {
					default:
						g.Id("req").Dot("Header").Dot("Add").Call(jen.Lit("Content-Type"), jen.Lit("application/json"))
						if len(b.ep.BodyParams) == 1 && b.ep.NoWrapRequest {
							g.Var().Id("body").Add(types.Convert(b.ep.BodyParams[0].Type, b.qualifier.Qual))
						} else {
							g.Var().Id("body").StructFunc(wrapResponse(b.ep.WrapRequest, func(g *jen.Group) {
								for _, param := range b.ep.BodyParams {
									jsonTag := param.Name
									fld := g.Id(param.FldName.Camel())
									if !param.Required {
										jsonTag += ",omitempty"
										if !isNamedType(param.Type) {
											fld.Op("*")
										}
									}
									fld.Add(types.Convert(param.Type, b.qualifier.Qual)).Tag(map[string]string{"json": jsonTag})
								}
							}, b.qualifier.Qual))
						}

						if len(b.ep.BodyParams) == 1 && b.ep.NoWrapRequest {
							g.Id("body").Op("=").Id(recvName).Dot("params").Dot(b.ep.BodyParams[0].FldName.LowerCamel())
						} else {
							for _, param := range b.ep.BodyParams {
								fldName := param.FldName.LowerCamel()
								if param.Parent != nil {
									fldName = param.Parent.FldName.LowerCamel() + param.FldName.String()
								}

								g.Do(func(s *jen.Statement) {
									s.Id("body")
									for _, name := range b.ep.WrapRequest {
										s.Dot(strcase.ToCamel(name))
									}
									s.Dot(param.FldName.Camel()).Op("=").Id(recvName).Dot("params").Dot(fldName)
								})

							}
						}
						g.Var().Id("reqData").Qual("bytes", "Buffer")
						g.Err().Op("=").Qual(jsonPkg, "NewEncoder").Call(jen.Op("&").Id("reqData")).Dot("Encode").Call(jen.Id("body"))
						g.If(jen.Err().Op("!=").Nil()).Block(
							jen.Id("cancel").Call(),
							jen.Return(),
						)
						g.Id("req").Dot("Body").Op("=").Qual("io", "NopCloser").Call(jen.Op("&").Id("reqData"))
					case "urlencoded":
						g.Id("req").Dot("Header").Dot("Add").Call(jen.Lit("Content-Type"), jen.Lit("application/x-www-form-urlencoded"))
						g.Id("body").Op(":=").Qual(urlPkg, "Values").Values()
						for _, param := range b.ep.BodyParams {

							g.Add(makeParam(jen.Dot(param.FldNameWithParent()), param.Type, param.Required, func(v jen.Code) jen.Code {
								transCode, paramID, _ := typetransform.For(param.Type).SetValueID(v).SetQualFunc(b.qualifier.Qual).Format()
								if transCode != nil {
									g.Add(transCode)
								}

								return jen.Id("body").Dot("Add").Call(jen.Lit(param.Name), paramID)
							}))

						}
						g.Id("req").Dot("Body").Op("=").Qual("io", "NopCloser").Call(jen.Qual(stringsPkg, "NewReader").Call(
							jen.Id("body").Dot("Encode").Call(),
						))
					}
				}
			}),
			jen.CustomFunc(jen.Options{Multi: true}, func(g *jen.Group) {
				if len(b.ep.QueryParams) > 0 || len(b.ep.QueryValues) > 0 {
					g.Id("q").Op(":=").Id("req").Dot("URL").Dot("Query").Call()

					for _, param := range b.ep.QueryParams {
						if named, ok := param.Type.(*types.Named); ok {
							for _, f := range named.Struct().Fields {
								if tag, err := f.SysTags.Get("json"); err == nil {
									g.Add(makeParam(jen.Dot(param.FldNameWithParent()), f.Type, false, func(v jen.Code) jen.Code {
										transCode, paramID, _ := typetransform.For(param.Type).SetValueID(v).SetQualFunc(b.qualifier.Qual).Format()
										if transCode != nil {
											g.Add(transCode)
										}

										return jen.Id("q").Dot("Add").Call(jen.Lit(tag.Name), paramID)
									}))
								}
							}
						} else {
							g.Add(makeParam(jen.Dot(param.FldNameWithParent()), param.Type, param.Required, func(v jen.Code) jen.Code {
								transCode, paramID, _ := typetransform.For(param.Type).SetValueID(v).SetQualFunc(b.qualifier.Qual).Format()
								if transCode != nil {
									g.Add(transCode)
								}

								return jen.Id("q").Dot("Add").Call(jen.Lit(param.Name), paramID)
							}))
						}
					}

					for _, param := range b.ep.QueryValues {
						g.Id("q").Dot("Add").Call(jen.Lit(param.Name), jen.Lit(param.Value))
					}

					g.Id("req").Dot("URL").Dot("RawQuery").Op("=").Id("q").Dot("Encode").Call()
				}

				for _, param := range b.ep.HeaderParams {
					g.Add(makeParam(jen.Dot(param.FldNameWithParent()), param.Type, param.Required, func(v jen.Code) jen.Code {
						transCode, paramID, _ := typetransform.For(param.Type).SetValueID(v).SetQualFunc(b.qualifier.Qual).Format()
						if transCode != nil {
							g.Add(transCode)
						}

						return jen.Id("req").Dot("Header").Dot("Add").Call(jen.Lit(param.Name), paramID)
					}))
				}
				for _, param := range b.ep.CookieParams {

					g.Add(makeParam(jen.Dot(param.FldNameWithParent()), param.Type, param.Required, func(v jen.Code) jen.Code {

						transCode, paramID, _ := typetransform.For(param.Type).SetValueID(v).SetQualFunc(b.qualifier.Qual).Format()
						if transCode != nil {
							g.Add(transCode)
						}

						return jen.Id("req").Dot("AddCookie").Call(jen.Op("&").Qual(httpPkg, "Cookie").Values(
							jen.Id("Name").Op(":").Lit(param.Name),
							jen.Id("Value").Op(":").Add(paramID),
						))
					}))
				}
			}),

			jen.Id("before").Op(":=").Append(jen.Id(recvName).Dot("c").Dot("opts").Dot("before"), jen.Id(recvName).Dot("opts").Dot("before").Op("...")),
			jen.For(jen.List(jen.Id("_"), jen.Id("before")).Op(":=").Range().Id("before")).Block(
				jen.List(jen.Id("ctx"), jen.Err()).Op("=").Id("before").Call(jen.Id("ctx"), jen.Id("req")),
				jen.Do(gen.CheckErr(
					jen.Id("cancel").Call(),
					jen.Return(),
				)),
			),
			jen.List(jen.Id("resp"), jen.Err()).Op(":=").Id(recvName).Dot("client").Dot("Do").Call(jen.Id("req")),
			jen.Do(gen.CheckErr(
				jen.Id("cancel").Call(),
				jen.Return(),
			)),

			jen.Id("after").Op(":=").Append(jen.Id(recvName).Dot("c").Dot("opts").Dot("after"), jen.Id(recvName).Dot("opts").Dot("after").Op("...")),
			jen.For(jen.List(jen.Id("_"), jen.Id("after")).Op(":=").Range().Id("after")).Block(
				jen.Id("ctx").Op("=").Id("after").Call(jen.Id("ctx"), jen.Id("resp")),
			),
			jen.Defer().Id("resp").Dot("Body").Dot("Close").Call(),
			jen.Defer().Id("cancel").Call(),

			jen.If(jen.Id("resp").Dot("StatusCode").Op(">").Lit(399)).BlockFunc(func(g *jen.Group) {
				g.If(jen.Id("resp").Dot("Body").Op("==").Qual(httpPkg, "NoBody")).Block(
					jen.Id("err").Op("=").Do(b.qualifier.Qual(fmtPkg, "Errorf")).Call(jen.Lit("http error %d"), jen.Id("resp").Dot("StatusCode")),
					jen.Return(),
				)
				if b.errorWrapper != nil {
					g.Var().Id("errorWrapper").Do(b.qualifier.Qual(b.errorWrapper.Struct.Named.Pkg.Path, b.errorWrapper.Struct.Named.Name))
					g.Var().Id("bytes").Index().Byte()
					g.List(jen.Id("bytes"), jen.Id("err")).Op("=").Qual(ioPkg, "ReadAll").Call(jen.Id("resp").Dot("Body"))
					g.Do(gen.CheckErr(
						jen.Id("err").Op("=").Do(b.qualifier.Qual(fmtPkg, "Errorf")).Call(jen.Lit("http error %d: %w"), jen.Id("resp").Dot("StatusCode"), jen.Id("err")),
						jen.Return(),
					))
					g.Id("err").Op("=").Qual(jsonPkg, "Unmarshal").Call(jen.Id("bytes"), jen.Op("&").Id("errorWrapper"))
					g.Do(gen.CheckErr(
						jen.Id("err").Op("=").Qual(fmtPkg, "Errorf").Call(jen.Lit("http error %d unmarshal data %s: %w"), jen.Id("resp").Dot("StatusCode"), jen.Id("bytes"), jen.Id("err")),
						jen.Return(),
					))
					g.Id("err").Op("=").Op("&").Do(b.qualifier.Qual(b.errorWrapper.Default.Named.Pkg.Path, b.errorWrapper.Default.Named.Name)).ValuesFunc(func(g *jen.Group) {
						for _, field := range b.errorWrapper.Fields {
							g.Id(strcase.ToCamel(field.FldName)).Op(":").Id("errorWrapper").Dot(field.FldName)
						}
						if b.errorWrapper.HasStatusCode {
							g.Id("StatusCode").Op(":").Id("resp").Dot("StatusCode")
						}
					})
					g.Return()
				} else {
					g.Id("err").Op("=").Do(b.qualifier.Qual(fmtPkg, "Errorf")).Call(jen.Lit("http error %d"), jen.Id("resp").Dot("StatusCode"))
					g.Return()
				}
			}),
			jen.CustomFunc(jen.Options{Multi: true}, func(g *jen.Group) {
				switch {
				default:
					g.Return()
				case len(b.ep.BodyResults) > 0:

					if !b.ep.NoWrapResponse {
						g.Var().Id("respBody").StructFunc(wrapResponse(b.ep.WrapResponse, func(g *jen.Group) {
							for _, result := range b.ep.BodyResults {
								g.Id(result.FldName.Camel()).Add(types.Convert(result.Type, b.qualifier.Qual)).Tag(map[string]string{"json": result.Name})
							}
						}, b.qualifier.Qual))
					} else if len(b.ep.BodyResults) == 1 {
						g.Var().Id("respBody").Add(types.Convert(b.ep.BodyResults[0].Type, b.qualifier.Qual))
					}
					g.Var().Id("reader").Qual("io", "ReadCloser")
					g.Switch(jen.Id("resp").Dot("Header").Dot("Get").Call(jen.Lit("Content-Encoding"))).Block(
						jen.Default().Block(jen.Id("reader").Op("=").Id("resp").Dot("Body")),
						jen.Case(jen.Lit("gzip")).Block(
							jen.List(jen.Id("reader"), jen.Err()).Op("=").Qual("compress/gzip", "NewReader").Call(jen.Id("resp").Dot("Body")),
							jen.Do(gen.CheckErr(
								jen.Return(),
							)),
							jen.Defer().Id("reader").Dot("Close").Call(),
						),
					)
					g.Id("err").Op("=").Qual(jsonPkg, "NewDecoder").
						Call(jen.Id("reader")).Dot("Decode").
						Call(jen.Op("&").Id("respBody"))
					g.Do(gen.CheckErr(
						jen.Return(),
					))
					g.ReturnFunc(func(g *jen.Group) {
						if len(b.ep.BodyResults) > 0 {
							if !b.ep.NoWrapResponse {
								var ids []jen.Code
								for _, name := range b.ep.WrapResponse {
									ids = append(ids, jen.Dot(strcase.ToCamel(name)))
								}
								for _, result := range b.ep.BodyResults {
									g.Id("respBody").Add(ids...).Dot(strcase.ToCamel(result.Name))
								}
							} else {
								g.Id("respBody")
							}
						}
						g.Nil()
					})
				case len(b.ep.HeaderResults) > 0:
					for _, r := range b.ep.HeaderResults {
						g.Id(r.FldName.LowerCamel()).Op("=").Id("resp").Dot("Header").Dot("Get").Call(jen.Lit(r.Name))
					}
					g.Return()
				}
			}),
		))
	return b
}

func (b *clientEndpointBuilder) buildSetter(parentParam, param *options.EndpointParam) {
	methodRequestName := b.methodRequestName()
	recvName := b.recvName()

	fldName := param.FldName.LowerCamel()
	fnName := param.FldName.Camel()
	if parentParam != nil {
		fldName = parentParam.FldName.LowerCamel() + param.FldName.String()
		fnName = parentParam.FldName.Camel() + param.FldName.Camel()
	}
	b.codes = append(b.codes,
		jen.Func().Params(
			jen.Id(recvName).Op("*").Id(methodRequestName),
		).Id("Set"+fnName).Params(
			jen.Id(fldName).Add(types.Convert(param.Type, b.qualifier.Qual)),
		).Op("*").Id(methodRequestName).BlockFunc(func(g *jen.Group) {
			g.Add(jen.CustomFunc(jen.Options{}, func(g *jen.Group) {
				g.Id(recvName).Dot("params").Dot(fldName).Op("=")
				if !param.Required && !types.IsPointer(param.Type) {
					g.Op("&")
				}
				g.Id(fldName)
			}))
			g.Return(jen.Id(recvName))
		}))
}

func (b *clientEndpointBuilder) clientName() string {
	return clientStructName(b.iface)
}

func (b *clientEndpointBuilder) methodRequestName() string {
	return b.iface.Name + b.ep.MethodName + "Request"
}

func (b *clientEndpointBuilder) recvName() string {
	return "r"
}

func (b *clientEndpointBuilder) methodReqName() string {
	return b.ep.MethodName + "Request"
}

func (b *clientEndpointBuilder) makeRequestStructParam(parentParam, param *options.EndpointParam, importFn types.QualFunc) jen.Code {
	fldName := param.FldName.LowerCamel()
	if parentParam != nil {
		fldName = parentParam.FldName.LowerCamel() + param.FldName.String()
	}
	paramID := jen.Id(fldName)
	if !param.Required && !types.IsPointer(param.Type) {
		paramID.Op("*")
	}
	paramID.Add(types.Convert(param.Type, importFn))
	return paramID
}
