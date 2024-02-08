package grpc

import (
	"go/token"
	"os/exec"
	"path/filepath"

	"github.com/555f/gg/internal/plugin/grpc/gen"
	"github.com/555f/gg/internal/plugin/grpc/options"
	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
	"github.com/hashicorp/go-multierror"
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

// func (p *Plugin) findGRPCInterface(iface *gg.Interface) (*gg.Interface, error) {
// 	if t, ok := iface.Named.Tags.Get("grpc-interface"); ok {
// 		pkgPath, name, err := p.ctx.Module.ParseImportPath(t.Value)
// 		if err != nil {
// 			return nil, err
// 		}
// 		for _, s := range p.ctx.OthersInterfaces {
// 			if s.Named.Pkg.Path == pkgPath && name == s.Named.Name {
// 				return s, nil
// 			}
// 		}
// 	} else {
// 		return nil, errors.Error("the \"grpc-interface\" tag is required", iface.Named.Position)
// 	}
// 	return nil, errors.Error("the interface for the gRPC server was not found", iface.Named.Position)
// }

// Exec implements gg.Plugin.
func (p *Plugin) Exec() (files []file.File, errs error) {
	// grpcInterfaces := map[string][]*types.Interface{}
	serverOutput := p.ctx.Options.GetStringWithDefault("server-output", "internal/server/server.go")
	serverAbsOutput := filepath.Join(p.ctx.Workdir, serverOutput)
	// clientOutput := filepath.Join(
	// 	p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("client-output", "internal/server/client.go"),
	// )
	// openapiOutput := filepath.Join(
	// 	p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("openapi-output", "docs/openapi.yaml"),
	// )
	// apiDocOutput := filepath.Join(
	// 	p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("apidoc-output", "docs/apidoc.html"),
	// )
	var (
		serverServices []options.Iface
	// 	clientServices  []options.Iface
	// 	openapiServices []options.Iface
	)

	p.protoDirOutput = filepath.Dir(serverAbsOutput)
	p.protoOutput = filepath.Join(p.protoDirOutput, "grpc.proto")

	f := file.NewTxtFile(p.protoOutput)

	f.WriteText("syntax = \"proto3\";\n\n")
	f.WriteText("option go_package = \"%s\";\n\n", filepath.Join(p.ctx.PkgPath, filepath.Dir(serverOutput)))
	f.WriteText("package server;\n\n")

	files = append(files, f)

	for _, iface := range p.ctx.Interfaces {
		// grpcIface, err := p.findGRPCInterface(iface)
		// if err != nil {
		// errs = multierror.Append(errs, err)
		// continue
		// }
		s, err := options.Decode(iface)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		// 	if s.Openapi.Enable {
		// 		openapiServices = append(openapiServices, s)
		// 	}
		// 	if s.Client.Enable {
		// 		clientServices = append(clientServices, s)
		// 	}
		if s.Server.Enable {
			serverServices = append(serverServices, s)
		}
	}

	var walkNamedStruct func(t any, visited map[string]struct{}, fn func(named *types.Named))
	walkNamedStruct = func(t any, visited map[string]struct{}, fn func(named *types.Named)) {
		named, ok := t.(*types.Named)
		if !ok {
			return
		}
		visitID := named.Pkg.Path + named.Name
		if _, ok := visited[visitID]; ok {
			return
		}

		s := named.Struct()
		if s == nil {
			return
		}
		fn(named)

		visited[visitID] = struct{}{}

		for _, sf := range s.Fields {
			walkNamedStruct(sf.Var.Type, visited, fn)
		}
	}

	for _, s := range serverServices {
		for _, ep := range s.Endpoints {
			for _, p := range ep.Results {

				walkNamedStruct(p.Type, map[string]struct{}{}, func(named *types.Named) {
					f.WriteText("message %s {", named.Name)
					for i, sf := range named.Struct().Fields {
						f.WriteText("%s %s = %d; ", goType2GRPC(sf.Var.Type), strcase.ToLowerCamel(sf.Var.Name), i+1)
					}
					f.WriteText("}\n")
				})
			}
		}
	}

	for _, s := range serverServices {
		for _, ep := range s.Endpoints {
			f.WriteText("message %sRequest {", ep.RPCMethodName)

			for _, p := range ep.Params {
				f.WriteText("%s %s = %s; ", goType2GRPC(p.Type), p.FldNameUnExport, p.Version)
			}

			f.WriteText("}\n")

			f.WriteText("message %sResponse {", ep.RPCMethodName)
			for _, p := range ep.Results {
				f.WriteText("%s %s = %s; ", goType2GRPC(p.Type), p.FldNameUnExport, p.Version)
			}
			f.WriteText("}\n")
		}
	}

	for _, s := range serverServices {
		f.WriteText("service %s {\n", s.Name)
		for _, ep := range s.Endpoints {
			f.WriteText("\rrpc %[1]s(%[1]sRequest) returns (%[1]sResponse) {}\n", ep.RPCMethodName)
		}
		f.WriteText("}")
	}

	serverFile := file.NewGoFile(p.ctx.Module, serverAbsOutput)
	serverFile.SetVersion(p.ctx.Version)

	// clientFile := file.NewGoFile(p.ctx.Module, clientOutput)
	// clientFile.SetVersion(p.ctx.Version)

	serverBuilder := gen.NewServerBuilder(serverFile)
	// clientBuilder := gen.NewBaseClientBuilder(clientFile)

	if len(serverServices) > 0 {
		for _, s := range serverServices {
			routeName := "route" + s.Name

			serverFile.Type().Id(routeName).Struct(
				jen.Id("svc").Do(serverFile.Import(s.PkgPath, s.Name)),
			)

			for _, ep := range s.Endpoints {
				requestName := ep.RPCMethodName + "Request"
				responseName := ep.RPCMethodName + "Response"

				serverFile.Func().Params(jen.Id("r").Op("*").Id(routeName)).Id(ep.RPCMethodName).Params(
					jen.Id("ctx").Qual("context", "Context"),
					jen.Id("req").Op("*").Id(requestName),
				).Params(
					jen.Op("*").Id(responseName),
					jen.Error(),
				).BlockFunc(func(g *jen.Group) {
					var params []jen.Code
					for _, p := range ep.Params {
						switch t := p.Type.(type) {
						case *types.Basic:
							switch {
							default:
								params = append(params, jen.Id("req").Dot(p.FldName))
							case t.IsInt():
								params = append(params, jen.Id("int").Call(jen.Id("req").Dot(p.FldName)))
							}
						}
					}

					g.Id("r").Dot("svc").Dot(ep.MethodName).Call(params...)
					g.Return(jen.Nil(), jen.Nil())
				},
				)
			}
		}

		// serverBuilder.RegisterHandlerStrategy("default", func() gen.HandlerStrategy {
		// 	return gen.NewHandlerStrategyGRPC()
		// })

		// for _, iface := range serverServices {
		// 	controllerBuilder := serverBuilder.Controller(iface)

		// 	controllerBuilder.BuildHandlers()

		// 	// 		for _, ep := range iface.Endpoints {
		// 	// 			controllerBuilder.Endpoint(ep).BuildReqStruct().
		// 	// 				BuildReqDec().
		// 	// 				BuildRespStruct().
		// 	// 				Build()
		// 	// 		}
		// }

		// // serverFile.Comment("// test")

		files = append(files, serverFile)
		serverFile.Add(serverBuilder.Build())
	}

	serverFile.Func().Id("NewServer").Params().Block()

	// if len(clientServices) > 0 {
	// 	for _, iface := range clientServices {
	// 		clientBuilder.BuildStruct(iface)
	// 		for _, ep := range iface.Endpoints {
	// 			clientBuilder.Endpoint(iface, ep).
	// 				BuildReqStruct().
	// 				BuildSetters().
	// 				BuildMethod().
	// 				BuildReqMethod().
	// 				BuildResultMethod().
	// 				BuildExecuteMethod()
	// 		}
	// 		clientBuilder.BuildConstruct(iface)
	// 	}
	// 	clientFile.Add(clientBuilder.Build())
	// 	files = append(files, clientFile)
	// }
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
