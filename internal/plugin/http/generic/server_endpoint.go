package generic

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/types"

	. "github.com/dave/jennifer/jen"
)

func GenServerEndpoints(s options.Iface) func(f *file.GoFile) {
	return func(f *file.GoFile) {
		for _, ep := range s.Endpoints {
			var (
				requestParams    []Code
				responseParams   []Code
				returnItems      []Code
				methodCallParams []Code
				assignResults    []Code
			)

			for _, param := range ep.Params {
				if len(param.Params) > 0 {
					for _, childParam := range param.Params {
						requestParams = append(requestParams, makeRequestParam(param, childParam, f.Import))
					}
				} else {
					requestParams = append(requestParams, makeRequestParam(nil, param, f.Import))
				}
			}

			if ep.Context != nil {
				methodCallParams = append(methodCallParams, Id("ctx"))
			}

			for _, param := range ep.Params {
				if len(param.Params) > 0 {
					if named, ok := param.Type.(*types.Named); ok {
						methodCallParams = append(methodCallParams, Qual(named.Pkg.Path, named.Name).ValuesFunc(func(g *Group) {
							for _, childParam := range param.Params {
								g.Id(childParam.FldName).Op(":").Id("r").Dot(param.FldNameUnExport + childParam.FldName)
							}
						}))
					}
				} else {
					st := Id("r").Dot(param.FldNameUnExport)
					if param.IsVariadic {
						st.Op("...")
					}
					methodCallParams = append(methodCallParams, st)
				}
			}

			if len(ep.BodyResults) > 0 {
				if len(ep.BodyResults) == 1 && ep.DisabledWrapResponse {
					returnItems = append(returnItems, Id(ep.BodyResults[0].FldNameUnExport))
					assignResults = append(assignResults, Id(ep.BodyResults[0].FldNameUnExport))
				} else {
					fields := Dict{}
					for _, result := range ep.BodyResults {
						responseParams = append(
							responseParams,
							Id(result.FldNameExport).Add(types.Convert(result.Type, f.Import)).Tag(map[string]string{"json": result.Name}),
						)
						assignResults = append(assignResults, Id(result.FldNameUnExport))
						fields[Id(result.FldNameExport)] = Id(result.FldNameUnExport)
					}
					st := Id(ep.RespStructName).Values(fields)
					returnItems = append(returnItems, st)
				}
			} else {
				returnItems = append(returnItems, Nil())
			}

			if len(requestParams) > 0 {
				f.Type().Id(ep.ReqStructName).Struct(requestParams...).Line()
			}

			if len(responseParams) > 0 {
				f.Type().Id(ep.RespStructName).Struct(responseParams...).Line()
			}

			epFuncDo := func(s *Statement) {
				s.Func().
					Params(Id("ctx").Qual("context", "Context"), Id("request").Any()).
					Params(Any(), Error())
			}

			if ep.Error != nil {
				returnItems = append(returnItems, Id(ep.Error.Name))
				assignResults = append(assignResults, Id(ep.Error.Name))
			} else {
				returnItems = append(returnItems, Nil())
			}

			f.Func().Id(ep.Name).
				Params(Id("s").Qual(s.PkgPath, s.Name)).
				Do(epFuncDo).
				Block(
					Return().Do(epFuncDo).BlockFunc(
						func(g *Group) {
							if len(requestParams) > 0 {
								g.Id("r").Op(":=").Id("request").Assert(Op("*").Id(ep.ReqStructName))
							}
							g.Do(func(s *Statement) {
								if len(assignResults) > 0 {
									s.List(assignResults...).Op(":=")
								}
								s.Id("s").Dot(ep.MethodName).Call(methodCallParams...)
							})
							g.Return(returnItems...)
						},
					),
				).
				Line()
		}
	}
}
