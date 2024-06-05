package grpc

import (
	"fmt"
	"go/token"
	"os/exec"
	"path/filepath"
	"strings"

	pgen "github.com/555f/gg/internal/plugin/grpc/ggen"
	"github.com/555f/gg/internal/plugin/grpc/options"
	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gen"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
	"github.com/hashicorp/go-multierror"
)

const (
	pkgGRPC     = "google.golang.org/grpc"
	pkgMetadata = "google.golang.org/grpc/metadata"
	pkgStrconv  = "strconv"
)

type Plugin struct {
	ctx            *gg.Context
	protoOutput    string
	protoDirOutput string
}

// Dependencies implements gg.Plugin.
func (*Plugin) Dependencies() []string {
	return nil
}

// Exec implements gg.Plugin.
func (p *Plugin) Exec() (files []file.File, errs error) {
	serverOutput := p.ctx.Options.GetStringWithDefault("server-output", "internal/server/server.go")
	serverAbsOutput := filepath.Join(p.ctx.Workdir, serverOutput)

	clientOutput := p.ctx.Options.GetStringWithDefault("client-output", "pkg/client/client.go")
	clientAbsOutput := filepath.Join(p.ctx.Workdir, clientOutput)

	p.protoDirOutput = filepath.Dir(serverAbsOutput)
	p.protoOutput = filepath.Join(p.protoDirOutput, "grpc.proto")

	pkgServer := p.ctx.PkgPath + "/" + filepath.Dir(serverOutput)

	pbf := pgen.NewFile("server")
	pbf.Syntax("proto3")
	pbf.Option().Id("go_package").Op("=").Lit(filepath.Join(p.ctx.PkgPath, filepath.Dir(serverOutput))).Op(";")

	protoFile := file.NewTxtFile(p.protoOutput)

	serverFile := file.NewGoFile(p.ctx.Module, serverAbsOutput)
	serverFile.SetVersion(p.ctx.Version)

	clientFile := file.NewGoFile(p.ctx.Module, clientAbsOutput)
	clientFile.SetVersion(p.ctx.Version)

	clientFile.Func().Id("toString").Params(jen.Id("v").Any()).String().Block(
		jen.Switch(jen.Id("t").Op(":=").Id("v").Assert(jen.Id("type"))).Block(
			jen.Case(jen.String()).Block(
				jen.Return(jen.Id("t")),
			),
			jen.Case(jen.Int()).Block(
				jen.Return(jen.Qual(pkgStrconv, "FormatInt").Call(jen.Int64().Call(jen.Id("t")), jen.Lit(10))),
			),
			jen.Case(jen.Int64()).Block(
				jen.Return(jen.Qual(pkgStrconv, "FormatInt").Call(jen.Id("t"), jen.Lit(10))),
			),
			jen.Case(jen.Float32()).Block(
				jen.Return(jen.Qual(pkgStrconv, "FormatFloat").Call(jen.Float64().Call(jen.Id("t")), jen.Id("'f'"), jen.Lit(10), jen.Lit(32))),
			),
			jen.Case(jen.Float64()).Block(
				jen.Return(jen.Qual(pkgStrconv, "FormatFloat").Call(jen.Id("t"), jen.Id("'f'"), jen.Lit(10), jen.Lit(64))),
			),
		),
		jen.Return(jen.Lit("")),
	)

	convertStructsCode := jen.Func().Id("convertStructs").Types(
		jen.Id("A").Any(),
		jen.Id("B").Any(),
		jen.Id("IN").Index().Id("A"),
	).Params(
		jen.Id("a").Id("IN"),
		jen.Id("c").Func().Params(jen.Id("A")).Id("B"),
	).Params(jen.Id("r").Index().Id("B")).Block(
		jen.For(jen.List(jen.Id("_"), jen.Id("v")).Op(":=").Range().Id("a")).Block(
			jen.Id("r").Op("=").Append(jen.Id("r"), jen.Id("convertStruct").Types(jen.Id("A"), jen.Id("B")).Call(
				jen.Id("v"),
				jen.Id("c"),
			)),
		),
		jen.Return(),
	)
	convertStructCode := jen.Func().Id("convertStruct").Types(
		jen.Id("A").Any(),
		jen.Id("B").Any(),
	).Params(
		jen.Id("a").Id("A"),
		jen.Id("c").Func().Params(jen.Id("A")).Id("B"),
	).Params(jen.Id("r").Id("B")).Block(
		jen.Return(jen.Id("c").Call(jen.Id("a"))),
	)

	serverFile.Add(convertStructCode)
	serverFile.Add(convertStructsCode)

	clientFile.Add(convertStructCode)
	clientFile.Add(convertStructsCode)

	clientBeforeFunc := jen.Func().Params(jen.Id("ctx").Qual("context", "Context")).Qual("context", "Context")
	clientAfterFunc := jen.Func().Params(jen.Id("ctx").Qual("context", "Context"))

	makeBeforeName := func(s options.Iface) string {
		return strcase.ToLowerCamel(s.Name) + "Before"
	}

	makeAfterName := func(s options.Iface) string {
		return strcase.ToLowerCamel(s.Name) + "After"
	}

	var (
		serverServices []options.Iface
		clientServices []options.Iface
	)

	for _, iface := range p.ctx.Interfaces {
		s, err := options.Decode(p.ctx.Module, iface)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		if s.Client.Enable {
			clientServices = append(clientServices, s)
		}
		if s.Server.Enable {
			serverServices = append(serverServices, s)
		}
	}

	serverQual := func(pkgPath, name string) func(s *jen.Statement) {
		return serverFile.Import(pkgServer, name)
	}

	clientQual := func(pkgPath, name string) func(s *jen.Statement) {
		return clientFile.Import(pkgServer, name)
	}

	var walkType func(t any, visited map[string]struct{}, fn func(named *types.Named))
	walkType = func(t any, visited map[string]struct{}, fn func(named *types.Named)) {
		switch t := t.(type) {
		case *types.Slice:
			walkType(t.Value, visited, fn)
		case *types.Named:
			if t.Pkg.Path == "time" {
				return
			}
			visitID := t.Pkg.Path + t.Name
			if _, ok := visited[visitID]; ok {
				return
			}
			visited[visitID] = struct{}{}

			s := t.Struct()
			if s == nil {
				return
			}
			fn(t)

			for _, sf := range s.Fields {
				walkType(sf.Var.Type, visited, fn)
			}
		case *types.Chan:
			walkType(t.Type, visited, fn)
			return
		}
	}

	var protobufToStuctRecursive func(path jen.Statement, t any, qualFn types.QualFunc) jen.Code
	protobufToStuct := func(path jen.Statement, t any, qualFn types.QualFunc) jen.Code {
		return protobufToStuctRecursive(path, t, qualFn)
	}
	protobufToStuctRecursive = func(path jen.Statement, t any, qualFn types.QualFunc) jen.Code {
		switch t := t.(type) {
		case *types.Slice:
			if named, ok := t.Value.(*types.Named); ok {

				return jen.Id("convertStructs").Call(jen.Add(&path), jen.Func().
					Params(jen.Id("a").Op("*").Qual(pkgServer, named.Name)).Do(func(s *jen.Statement) {
					// s.Id("a")
					// if named.IsPointer {
					// s.Op("*")
					// }
					// s.Qual(named.Pkg.Path, named.Name)
				}).Do(func(s *jen.Statement) {
					if named.IsPointer {
						s.Op("*")
					}
					s.Qual(named.Pkg.Path, named.Name)
				}).Block(
					jen.Return(
						jen.Add(protobufToStuct(*jen.Id("a"), t.Value, qualFn)),
					),
				))
			}
		case *types.Named:
			switch t.Pkg.Name {
			case "time":
				switch t.Name {
				case "Duration":
					return jen.Add(&path).Dot("AsDuration").Call()
				case "Time":
					return jen.Add(&path).Dot("AsTime").Call()
				}
			}
			var values []jen.Code
			if s := t.Struct(); s != nil {
				for _, sf := range s.Fields {
					pbName := strcase.ToCamel(strcase.ToLowerCamel(sf.Var.Name))
					values = append(
						values,
						jen.Id(sf.Var.Name).Op(":").Add(
							protobufToStuctRecursive(*jen.Add(path...).Dot(pbName), sf.Var.Type, qualFn),
						),
					)
				}
			}

			code := jen.Do(func(s *jen.Statement) {
				if t.IsPointer {
					s.Op("&")
				}
				qualFn(t.Pkg.Path, t.Name)(s)
			}).Values(values...)
			return code
		case *types.Basic:
			switch {
			default:
				return jen.Add(&path)
			case t.IsInt():
				return jen.Id("int").Call(jen.Add(&path))
			case t.IsInt8():
				return jen.Id("int8").Call(jen.Add(&path))
			case t.IsInt16():
				return jen.Id("int16").Call(jen.Add(&path))
			case t.IsUint():
				return jen.Id("uint").Call(jen.Add(&path))
			case t.IsUint8():
				return jen.Id("uint8").Call(jen.Add(&path))
			case t.IsInt16():
				return jen.Id("uint16").Call(jen.Add(&path))
			}
		}
		return nil
	}

	var structToProtobufRecursive func(path jen.Statement, t any, qualFn types.QualFunc) jen.Code
	structToProtobuf := func(path jen.Statement, t any, qualFn types.QualFunc) jen.Code {
		return structToProtobufRecursive(path, t, qualFn)
	}
	structToProtobufRecursive = func(path jen.Statement, t any, qualFn types.QualFunc) jen.Code {
		switch t := t.(type) {
		case *types.Slice:
			if named, ok := t.Value.(*types.Named); ok {
				return jen.Id("convertStructs").Call(jen.Add(&path), jen.Func().
					ParamsFunc(func(g *jen.Group) {
						g.Do(func(s *jen.Statement) {
							s.Id("a")
							if named.IsPointer {
								s.Op("*")
							}
							s.Qual(named.Pkg.Path, named.Name)
						})
					}).Op("*").Id(named.Name).Block(
					jen.Return(
						jen.Add(structToProtobuf(*jen.Id("a"), t.Value, qualFn)),
					),
				))
			}
		case *types.Named:
			switch t.Pkg.Name {
			case "time":
				switch t.Name {
				case "Duration":
					return jen.Qual("google.golang.org/protobuf/types/known/durationpb", "New").Call(jen.Add(&path))
				case "Time":
					return jen.Qual("google.golang.org/protobuf/types/known/timestamppb", "New").Call(jen.Add(&path))
				}
			}
			var values []jen.Code
			if s := t.Struct(); s != nil {
				for _, sf := range s.Fields {
					pbName := strcase.ToCamel(strcase.ToLowerCamel(sf.Var.Name))
					values = append(
						values,
						jen.Id(pbName).Op(":").Add(
							structToProtobufRecursive(*jen.Add(path...).Dot(sf.Var.Name), sf.Var.Type, qualFn),
						),
					)
				}
			}
			return jen.Op("&").Do(qualFn(t.Pkg.Path, t.Name)).Values(values...)
		case *types.Basic:
			switch {
			default:
				return jen.Add(&path)
			case t.IsInt():
				return jen.Id("int64").Call(jen.Add(&path))
			case t.IsInt8(), t.IsInt16():
				return jen.Id("int32").Call(jen.Add(&path))
			case t.IsUint():
				return jen.Id("uint").Call(jen.Add(&path))
			case t.IsUint8(), t.IsInt16():
				return jen.Id("uint32").Call(jen.Add(&path))
			}
		}
		return nil
	}

	visited := map[string]struct{}{}

	for _, s := range serverServices {
		for _, ep := range s.Endpoints {
			for _, p := range ep.Results {
				walkType(p.Type, visited, func(named *types.Named) {
					var values []pgen.Code
					for i, sf := range named.Struct().Fields {
						values = append(values, pgen.Id(goType2GRPC(sf.Var.Type)).Id(strcase.ToLowerCamel(sf.Var.Name)).Op("=").Id(fmt.Sprint(i+1)))
					}
					pbf.Message().Id(named.Name).Values(values...)
				})
			}
		}
	}

	for _, s := range serverServices {
		for _, ep := range s.Endpoints {
			var paramValues, resultValues []pgen.Code
			if ep.InStream == nil && len(ep.Params) > 0 {
				for _, p := range ep.Params {
					paramValues = append(paramValues, pgen.Qual(goType2GRPC(p.Type)).Id(p.FldNameUnExport).Op("=").Id(p.Version))
				}
				pbf.Message().Id(ep.RPCMethodName + "Request").Values(paramValues...)
			}
			if ep.OutStream == nil && len(ep.Results) > 0 {
				for _, p := range ep.Results {
					resultValues = append(resultValues, pgen.Qual(goType2GRPC(p.Type)).Id(p.FldNameUnExport).Op("=").Id(p.Version))
				}
				pbf.Message().Id(ep.RPCMethodName + "Response").Values(resultValues...)
			}
		}
	}

	for _, s := range serverServices {
		var methods []pgen.Code
		for _, ep := range s.Endpoints {
			rpcMethod := pgen.RPC().Id(ep.RPCMethodName)
			if ep.InStream != nil {
				rpcMethod.Request(pgen.Stream().Qual(goType2GRPC(ep.InStream.Chan.Type)))
			} else {
				rpcMethod.Request(pgen.Id(ep.RPCMethodName + "Request"))
			}
			if len(ep.Results) > 0 {
				if ep.OutStream != nil {
					rpcMethod.Returns(pgen.Stream().Qual(goType2GRPC(ep.OutStream.Chan.Type)))
				} else {
					rpcMethod.Returns(pgen.Id(ep.RPCMethodName + "Response"))
				}
			} else {
				rpcMethod.Returns(pgen.Qual("google.protobuf.Empty"))
			}
			methods = append(methods, rpcMethod)
		}
		pbf.Service(s.Name, methods...)
	}

	protoFile.WriteText(pbf.String())

	if len(serverServices) > 0 {

		optionsName := "options"
		optionName := "Option"
		beforeFunc := jen.Func().Params(jen.Id("ctx").Qual("context", "Context")).Qual("context", "Context")
		afterFunc := jen.Func().Params(jen.Id("ctx").Qual("context", "Context"))

		serverFile.Type().Id(optionsName).StructFunc(func(g *jen.Group) {
			for _, s := range serverServices {
				g.Id(strcase.ToLowerCamel(s.Name)).Do(serverFile.Import(s.PkgPath, s.Name))
				g.Id(makeBeforeName(s)).Index().Add(beforeFunc)
				g.Id(makeAfterName(s)).Index().Add(afterFunc)
			}

		})
		serverFile.Type().Id(optionName).Func().Params(jen.Op("*").Id(optionsName))

		for _, s := range serverServices {
			beforeName := makeBeforeName(s)
			afterName := makeAfterName(s)

			serverFile.Func().Id(s.Name).Params(
				jen.Id("s").Do(serverFile.Import(s.PkgPath, s.Name)),
			).Id("Option").Block(
				jen.Return(
					jen.Func().Params(jen.Id("o").Op("*").Id(optionsName)).Block(
						jen.Id("o").Dot(strcase.ToLowerCamel(s.Name)).Op("=").Id("s"),
					),
				),
			)
			serverFile.Func().Id(s.Name + "Before").Params(jen.Id("before").Op("...").Add(beforeFunc)).Id(optionName).Block(
				jen.Return(
					jen.Func().Params(jen.Id("o").Op("*").Id(optionsName)).Block(
						jen.Id("o").Dot(beforeName).Op("=").Append(jen.Id("o").Dot(beforeName), jen.Id("before").Op("...")),
					),
				),
			)

			serverFile.Func().Id(s.Name + "After").Params(jen.Id("after").Op("...").Add(afterFunc)).Id(optionName).Block(
				jen.Return(
					jen.Func().Params(jen.Id("o").Op("*").Id(optionsName)).Block(
						jen.Id("o").Dot(afterName).Op("=").Append(jen.Id("o").Dot(afterName), jen.Id("after").Op("...")),
					),
				),
			)
		}

		for _, s := range serverServices {
			routeName := "route" + s.Name

			serverFile.Type().Id(routeName).Struct(
				jen.Id("Unimplemented"+s.Name+"Server"),
				jen.Id("svc").Do(serverFile.Import(s.PkgPath, s.Name)),
				jen.Id("before").Index().Add(beforeFunc),
				jen.Id("after").Index().Add(afterFunc),
			)

			for _, ep := range s.Endpoints {
				requestName := ep.RPCMethodName + "Request"
				responseName := ep.RPCMethodName + "Response"
				useContext := ep.InStream == nil && ep.OutStream == nil

				serverFile.Func().Params(jen.Id("r").Op("*").Id(routeName)).Id(ep.RPCMethodName).
					ParamsFunc(func(g *jen.Group) {
						if useContext {
							g.Id("ctx").Qual("context", "Context")
						}
						if ep.InStream == nil && ep.Params.Len() > 0 {
							g.Id("req").Op("*").Id(requestName)
						}
						if ep.InStream != nil || ep.OutStream != nil {
							g.Id("stream").Id(s.Name + "_" + ep.MethodName + "Server")
						}
					}).
					ParamsFunc(func(g *jen.Group) {
						if ep.InStream == nil && ep.OutStream == nil {
							if len(ep.Results) > 0 {
								g.Op("*").Id(responseName)
							} else {
								g.Op("*").Qual("google.golang.org/protobuf/types/known/emptypb", "Empty")
							}
						}
						g.Error()
					}).BlockFunc(func(g *jen.Group) {
					if !useContext {
						g.Id("ctx").Op(":=").Id("stream").Dot("Context").Call()
					}
					if ep.InStream != nil {
						g.Id("chIn").Op(":=").Make(jen.Chan().Add(types.Convert(ep.InStream.Chan.Type, serverFile.Import)))

						g.Go().Func().Params().Block(
							jen.For().BlockFunc(func(g *jen.Group) {
								g.List(jen.Id("data"), jen.Err()).Op(":=").Id("stream").Dot("Recv").Call()
								g.Do(gen.CheckErr(jen.Return()))

								g.Id("chIn").Id("<-").Add(protobufToStuct(*jen.Id("data"), ep.InStream.Chan.Type, serverFile.Import))
							}),
						).Call()
					}

					if ep.Context != nil && len(ep.MetaContexts) > 0 {
						g.List(jen.Id("md"), jen.Id("ok")).Op(":=").Qual(pkgMetadata, "FromIncomingContext").Call(jen.Id(ep.Context.Name))
						g.If(jen.Id("ok")).BlockFunc(func(g *jen.Group) {
							for _, mc := range ep.MetaContexts {
								g.If(jen.Id("values").Op(":=").Id("md").Dot("Get").Call(jen.Lit(strings.ToLower(mc.Name))), jen.Len(jen.Id("values")).Op(">").Lit(0).Op("&&").Id("values").Index(jen.Lit(0)).Op("!=").Lit("")).Block(
									jen.Id("ctx").Op("=").Qual("context", "WithValue").Call(jen.Id("ctx"), jen.Qual(mc.PkgPath, mc.Name), jen.Id("values").Index(jen.Lit(0))),
								)
							}
						})
					}

					g.For(jen.List(jen.Id("_"), jen.Id("f")).Op(":=").Range().Id("r").Dot("before")).Block(
						jen.Id("ctx").Op("=").Id("f").Call(jen.Id("ctx")),
					)

					g.Do(func(s *jen.Statement) {
						s.ListFunc(func(g *jen.Group) {
							for _, p := range ep.Results {
								g.Id(p.FldName)
							}
							if ep.Error != nil {
								g.Id(ep.Error.Name)
							}
						})
						if len(ep.Results) > 0 || ep.Error != nil {
							s.Op(":=")
						}
					}).Id("r").Dot("svc").Dot(ep.MethodName).CallFunc(func(g *jen.Group) {
						if ep.InStream != nil {
							g.Id("chIn")
						} else {
							if ep.Context != nil {
								g.Id(ep.Context.Name)
							}
							for _, p := range ep.Params {
								g.Do(func(s *jen.Statement) {
									s.Add(protobufToStuct(*jen.Id("req").Dot(p.FldName), p.Type, serverFile.Import))
								})
							}
						}
					})

					hasResponse := hasResponseEndpoint(ep)

					g.Do(gen.CheckErr(jen.ReturnFunc(func(g *jen.Group) {
						if hasResponse {
							g.Nil()
						}
						if ep.Error != nil {
							g.Err()
						}
					})))

					g.For(jen.List(jen.Id("_"), jen.Id("f")).Op(":=").Range().Id("r").Dot("after")).Block(
						jen.Id("f").Call(jen.Id("ctx")),
					)

					if ep.OutStream != nil {
						g.For(jen.Id("data").Op(":=").Range().Id(ep.OutStream.Param.FldNameUnExport)).BlockFunc(func(g *jen.Group) {
							g.Id("stream").Dot("Send").Call(structToProtobuf(*jen.Id("data"), ep.OutStream.Chan.Type, serverQual))
						})
					} else if ep.OutStream == nil && len(ep.Results) > 0 {
						g.Id("resp").Op(":=").Op("&").Id(responseName).ValuesFunc(func(g *jen.Group) {
							for _, p := range ep.Results {
								g.Id(p.FldNameExport).Op(":").Add(structToProtobuf(*jen.Id(p.FldName), p.Type, serverQual))
							}
						})

					}

					g.ReturnFunc(func(g *jen.Group) {
						if hasResponse {
							if ep.Results.Len() > 0 {
								g.Id("resp")
							} else {
								g.Nil()
							}
						}
						if ep.Error != nil {
							g.Err()
						}
					})
				})
			}
		}
	}

	serverFile.Func().Id("Register").Params(
		jen.Id("srv").Op("*").Qual(pkgGRPC, "Server"),
		jen.Id("opts").Op("...").Id("Option"),
	).BlockFunc(func(g *jen.Group) {
		g.Id("o").Op(":=").Op("&").Id("options").Values()
		g.For(jen.List(jen.Id("_"), jen.Id("f")).Op(":=").Range().Id("opts")).Block(
			jen.Id("f").Call(jen.Id("o")),
		)
		for _, s := range serverServices {
			fldSvcName := strcase.ToLowerCamel(s.Name)
			g.If(jen.Id("o").Dot(fldSvcName).Op("!=").Nil()).Block(
				jen.Id("Register"+s.Name+"Server").Call(
					jen.Id("srv"),
					jen.Op("&").Id("route"+s.Name).Values(
						jen.Id("svc").Op(":").Id("o").Dot(fldSvcName),
						jen.Id("before").Op(":").Id("o").Dot(makeBeforeName(s)),
						jen.Id("after").Op(":").Id("o").Dot(makeAfterName(s)),
					),
				),
			)
		}
	})

	beforeFunc := jen.Func().Params(jen.Id("ctx").Qual("context", "Context")).Qual("context", "Context")
	afterFunc := jen.Func().Params(jen.Id("ctx").Qual("context", "Context"))

	for _, s := range clientServices {
		clientStructName := s.Name + "Client"
		optionName := s.Name + "Option"

		clientFile.Type().Id(optionName).Func().Params(jen.Op("*").Id(clientStructName))

		clientFile.Func().Id(s.Name + "Before").Params(jen.Id("before").Op("...").Add(beforeFunc)).Id(optionName).Block(
			jen.Return(
				jen.Func().Params(jen.Id("o").Op("*").Id(clientStructName)).Block(
					jen.Id("o").Dot("before").Op("=").Append(jen.Id("o").Dot("before"), jen.Id("before").Op("...")),
				),
			),
		)

		clientFile.Func().Id(s.Name + "After").Params(jen.Id("after").Op("...").Add(afterFunc)).Id(optionName).Block(
			jen.Return(
				jen.Func().Params(jen.Id("o").Op("*").Id(clientStructName)).Block(
					jen.Id("o").Dot("after").Op("=").Append(jen.Id("o").Dot("after"), jen.Id("after").Op("...")),
				),
			),
		)

		clientFile.Type().Id(clientStructName).Struct(
			jen.Id("cc").Qual(pkgServer, s.Name+"Client"),
			jen.Id("before").Index().Add(clientBeforeFunc),
			jen.Id("after").Index().Add(clientAfterFunc),
		)
		for _, ep := range s.Endpoints {
			clientFile.Func().
				Params(jen.Id("c").Op("*").Id(clientStructName)).
				Id(ep.MethodName).
				Add(types.Convert(ep.Sig, clientFile.Import)).
				BlockFunc(func(g *jen.Group) {
					if ep.Context == nil {
						g.Id("ctx").Op(":=").Qual("context", "TODO").Call()
					}

					respVar := jen.Id("resp")
					assignOp := ":="

					if ep.InStream == nil && ep.OutStream == nil && len(ep.Results) == 0 {
						respVar = jen.Id("_")
						assignOp = "="
					}

					g.Id("ctx").Op("=").Qual(pkgMetadata, "AppendToOutgoingContext").CallFunc(func(g *jen.Group) {
						g.Id("ctx")
						for _, mc := range ep.MetaContexts {
							g.Lit(strings.ToLower(mc.Name))
							g.Id("toString").Call(jen.Id("ctx").Dot("Value").Call(jen.Qual(mc.PkgPath, mc.Name)))
						}
					})

					g.For(jen.List(jen.Id("_"), jen.Id("f")).Op(":=").Range().Id("c").Dot("before")).Block(
						jen.Id("ctx").Op("=").Id("f").Call(jen.Id("ctx")),
					)

					g.List(respVar, jen.Err()).Op(assignOp).Id("c").Dot("cc").Dot(ep.MethodName).CallFunc(func(g *jen.Group) {
						g.Id("ctx")
						if ep.InStream == nil {
							g.Op("&").Qual(pkgServer, ep.MethodName+"Request").ValuesFunc(func(g *jen.Group) {
								for _, p := range ep.Params {
									g.Id(strcase.ToCamel(p.Name)).Op(":").Add(structToProtobuf(*jen.Id(p.Name), p.Type, clientQual))
								}
							})
						}
					})

					g.Do(gen.CheckErr(
						jen.Return(),
					))

					hasResponse := hasResponseEndpoint(ep)
					if hasResponse {
						for _, p := range ep.Results {
							g.Id(p.FldName).Op("=").Do(func(s *jen.Statement) {
							}).Add(protobufToStuct(*jen.Id("resp").Dot(p.FldNameExport), p.Type, clientFile.Import))
						}
					}

					g.For(jen.List(jen.Id("_"), jen.Id("f")).Op(":=").Range().Id("c").Dot("after")).Block(
						jen.Id("f").Call(jen.Id("ctx")),
					)

					if ep.OutStream != nil {
						g.Id(ep.OutStream.Param.Name).Op("=").Make(jen.Chan().Add(types.Convert(ep.OutStream.Chan.Type, clientFile.Import)))

						g.Go().Func().Params().Block(
							jen.Defer().Close(jen.Id(ep.OutStream.Param.Name)),
							jen.For().BlockFunc(func(g *jen.Group) {
								g.List(jen.Id("data"), jen.Err()).Op(":=").Id("resp").Dot("Recv").Call()
								g.Do(gen.CheckErr(
									jen.Return(),
								))
								g.Id(ep.OutStream.Param.Name).Op("<-").Add(protobufToStuct(*jen.Id("data"), ep.OutStream.Chan.Type, clientFile.Import))
							}),
						).Call()
					}

					if ep.InStream != nil {
						g.Go().Func().Params().Block(
							jen.For(jen.Id("data").Op(":=").Range().Id(ep.InStream.Param.Name)).Block(
								jen.List(jen.Err()).Op(":=").Id("resp").Dot("Send").Call(structToProtobuf(*jen.Id("data"), ep.InStream.Chan.Type, clientQual)),
								jen.Do(gen.CheckErr(
									jen.Break(),
								)),
							),
						).Call()
					}

					g.Return()
				})
		}
		clientFile.Func().Id("New" + s.Name + "Client").Params(
			jen.Id("cc").Qual(pkgGRPC, "ClientConnInterface"),
		).Op("*").Id(clientStructName).BlockFunc(func(g *jen.Group) {
			g.Return(
				jen.Op("&").Id(clientStructName).Values(
					jen.Id("cc").Op(":").Qual(pkgServer, "New"+s.Name+"Client").Call(jen.Id("cc")),
				),
			)
		})
	}

	files = append(files, clientFile, protoFile, serverFile)

	return
}

func (p *Plugin) OnAfterGen() error {
	cmd := exec.Command("protoc",
		"-I", p.protoDirOutput,
		"--go_out=.",
		"--go_opt=paths=source_relative",
		"--go-grpc_out=.",
		"--go-grpc_opt=paths=source_relative",
		"grpc.proto",
	)
	cmd.Dir = p.protoDirOutput
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Error(string(out), token.Position{})
	}
	return nil
}

// Name implements gg.Plugin.
func (*Plugin) Name() string {
	return "grpc"
}
