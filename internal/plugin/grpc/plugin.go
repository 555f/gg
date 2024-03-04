package grpc

import (
	"fmt"
	"go/token"
	"os/exec"
	"path/filepath"

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
	pkgGRPC = "google.golang.org/grpc"
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

	var (
		serverServices []options.Iface
		clientServices []options.Iface
	)

	for _, iface := range p.ctx.Interfaces {
		s, err := options.Decode(iface)
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
		serverFile.Type().Id("options").StructFunc(func(g *jen.Group) {
			for _, s := range serverServices {
				g.Id(strcase.ToLowerCamel(s.Name)).Do(serverFile.Import(s.PkgPath, s.Name))
			}
		})
		serverFile.Type().Id("Option").Func().Params(jen.Op("*").Id("options"))

		for _, s := range serverServices {
			serverFile.Func().Id(s.Name).Params(
				jen.Id("s").Do(serverFile.Import(s.PkgPath, s.Name)),
			).Id("Option").Block(
				jen.Return(
					jen.Func().Params(jen.Id("o").Op("*").Id("options")).Block(
						jen.Id("o").Dot(strcase.ToLowerCamel(s.Name)).Op("=").Id("s"),
					),
				),
			)
		}

		for _, s := range serverServices {
			routeName := "route" + s.Name

			serverFile.Type().Id(routeName).Struct(
				jen.Id("Unimplemented"+s.Name+"Server"),
				jen.Id("svc").Do(serverFile.Import(s.PkgPath, s.Name)),
			)

			for _, ep := range s.Endpoints {
				requestName := ep.RPCMethodName + "Request"
				responseName := ep.RPCMethodName + "Response"

				serverFile.Func().Params(jen.Id("r").Op("*").Id(routeName)).Id(ep.RPCMethodName).
					ParamsFunc(func(g *jen.Group) {
						if ep.InStream == nil && ep.OutStream == nil {
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
							for _, p := range ep.Params {
								g.Do(func(s *jen.Statement) {
									if p.IsPointer {
										s.Op("&")
									}
									s.Add(protobufToStuct(*jen.Id("req").Dot(p.FldName), p.Type, serverFile.Import))
								})
							}
						}
					})

					if ep.OutStream != nil {
						g.For(jen.Id("data").Op(":=").Range().Id(ep.OutStream.Param.FldNameUnExport)).BlockFunc(func(g *jen.Group) {
							g.Id("stream").Dot("Send").Call(structToProtobuf(*jen.Id("data"), ep.OutStream.Chan.Type, serverQual))
						})
					}

					hasResponse := hasResponseEndpoint(ep)

					g.Do(gen.CheckErr(jen.ReturnFunc(func(g *jen.Group) {
						if hasResponse {
							g.Nil()
						}
						if ep.Error != nil {
							g.Err()
						}
					})))

					if ep.OutStream == nil && len(ep.Results) > 0 {
						g.Var().Id("resp").Op("*").Id(responseName)
						for _, p := range ep.Results {
							g.Id("resp").Dot(p.FldNameExport).Op("=").Add(structToProtobuf(*jen.Id(p.FldName), p.Type, serverQual))
						}
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
					),
				),
			)
		}
	})

	for _, s := range clientServices {
		clientStructName := s.Name + "Client"
		clientFile.Type().Id(clientStructName).Struct(
			jen.Id("cc").Qual(pkgServer, s.Name+"Client"),
		)
		for _, ep := range s.Endpoints {
			clientFile.Func().
				Params(jen.Id("c").Op("*").Id(clientStructName)).
				Id(ep.MethodName).
				Add(types.Convert(ep.Sig, clientFile.Import)).
				BlockFunc(func(g *jen.Group) {
					contextVar := jen.Id("ctx")
					if ep.Context == nil {
						contextVar = jen.Qual("context", "TODO").Call()
					}

					respVar := jen.Id("resp")
					assignOp := ":="

					if ep.InStream == nil && ep.OutStream == nil && len(ep.Results) == 0 {
						respVar = jen.Id("_")
						assignOp = "="
					}

					g.List(respVar, jen.Err()).Op(assignOp).Id("c").Dot("cc").Dot(ep.MethodName).CallFunc(func(g *jen.Group) {
						g.Add(contextVar)
						if ep.InStream == nil {
							g.Op("&").Qual(pkgServer, ep.MethodName+"Request").ValuesFunc(func(g *jen.Group) {
								for _, p := range ep.Params {
									g.Id(strcase.ToCamel(p.Name)).Op(":").Add(structToProtobuf(*jen.Id(p.Name), p.Type, clientQual))
								}
							})
						}
					})

					hasResponse := hasResponseEndpoint(ep)
					if hasResponse {
						for _, p := range ep.Results {
							g.Id(p.FldName).Op("=").Do(func(s *jen.Statement) {
								if p.IsPointer {
									s.Op("&")
								}
							}).Add(protobufToStuct(*jen.Id("resp").Dot(p.FldNameExport), p.Type, clientFile.Import))
						}
					}

					g.Do(gen.CheckErr(
						jen.Return(),
					))

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
		clientFile.Func().Id("New" + s.Name + "Cient").Params(
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
