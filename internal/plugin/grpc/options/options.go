package options

import (
	"fmt"

	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"

	"github.com/hashicorp/go-multierror"
)

type Iface struct {
	Type        string
	Name        string
	Title       string
	Description string
	PkgPath     string
	Server      Server
	Client      Client
	Openapi     OpenAPI
	Endpoints   []Endpoint
	HTTPReq     string
}

type Server struct {
	Enable bool
}

type Client struct {
	Enable bool
}

type OpenAPI struct {
	Enable bool
	Tags   []string
}

type APIDoc struct {
	Enable bool
}

type EndpointParams []*EndpointParam

func (params EndpointParams) Len() int {
	return len(params)
}

type EndpointResults []*EndpointResult

func (params EndpointResults) Len() int {
	return len(params)
}

type InStream struct {
	Param *EndpointParam
	Chan  *types.Chan
}

type OutStream struct {
	Param *EndpointResult
	Chan  *types.Chan
}

type MetaContext struct {
	Name    string
	PkgPath string
}

type Endpoint struct {
	MethodName    string
	RPCMethodName string
	Title         string
	Description   string
	OpenapiTags   []string
	Params        EndpointParams
	Results       EndpointResults
	InStream      *InStream
	OutStream     *OutStream
	Context       *types.Var
	Error         *types.Var
	Sig           *types.Sign
	MetaContexts  []MetaContext
}

type EndpointParam struct {
	Type            any
	Title           string
	Name            string
	FldName         string
	FldNameUnExport string
	Format          string
	Zero            string
	IsOmitempty     bool
	IsVariadic      bool
	IsRequired      bool
	IsStream        bool
	IsPointer       bool
	Version         string
	Params          EndpointParams
}

type EndpointResult struct {
	Type            any
	Title           string
	Name            string
	FldName         string
	FldNameExport   string
	FldNameUnExport string
	Format          string
	IsOmitempty     bool
	IsStream        bool
	IsPointer       bool
	Version         string
}

func Decode(module *types.Module, iface *gg.Interface) (opts Iface, errs error) {
	opts.Name = iface.Named.Name
	opts.Title = iface.Named.Title
	opts.Description = iface.Named.Description
	opts.PkgPath = iface.Named.Pkg.Path
	opts.Type = "default"

	if _, ok := iface.Named.Tags.Get("grpc-server"); ok {
		opts.Server.Enable = true
	}
	if _, ok := iface.Named.Tags.Get("grpc-client"); ok {
		opts.Client.Enable = true
	}
	if _, ok := iface.Named.Tags.Get("grpc-openapi"); ok {
		opts.Openapi.Enable = true
	}
	if t, ok := iface.Named.Tags.Get("grpc-openapi-tags"); ok {
		opts.Openapi.Tags = []string{t.Value}
		opts.Openapi.Tags = append(opts.Openapi.Tags, t.Options...)
	}
	if t, ok := iface.Named.Tags.Get("grpc-req"); ok {
		opts.HTTPReq = t.Value
	}
	for _, method := range iface.Named.Interface().Methods {
		epOpts, err := endpointDecode(module, method)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		opts.Endpoints = append(opts.Endpoints, epOpts)
	}
	return
}

func endpointDecode(module *types.Module, method *types.Func) (opts Endpoint, errs error) {
	opts.MethodName = method.Name
	opts.Title = method.Title
	opts.Description = method.Description
	opts.Sig = method.Sig
	opts.RPCMethodName = strcase.ToCamel(method.Name)
	if t, ok := method.Tags.Get("grpc-openapi-tags"); ok {
		opts.OpenapiTags = []string{t.Value}
		opts.OpenapiTags = append(opts.OpenapiTags, t.Options...)
	}
	if t, ok := method.Tags.Get("grpc-name"); ok {
		opts.RPCMethodName = t.Value
	}

	tags := method.Tags.GetSlice("grpc-meta-context")
	for _, t := range tags {
		if t.Value == "" {
			errs = multierror.Append(errs, errors.Error("the path to the context key is required", t.Position))
			return
		}
		pkgPath, name, err := module.ParseImportPath(t.Value)
		if err != nil {
			errs = multierror.Append(errs, err)
			return
		}
		opts.MetaContexts = append(opts.MetaContexts, MetaContext{
			Name:    name,
			PkgPath: pkgPath,
		})
	}
	fmt.Println(opts.MetaContexts)

	for _, param := range method.Sig.Params {
		if param.IsContext {
			if opts.Context != nil {
				errs = multierror.Append(errs, errors.Error("the method has more than one context", param.Position))
			}
			opts.Context = param
			continue
		}
		if param.Name == "" {
			errs = multierror.Append(errs, errors.Error("the parameter name cannot be empty or the grpc-name parameter must be set", param.Position))
		}
		p, err := paramDecode(param)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		if param.IsChan {
			p.IsStream = true
			opts.InStream = &InStream{
				Param: p,
				Chan:  param.Type.(*types.Chan),
			}
		}
		opts.Params = append(opts.Params, p)
	}

	if opts.InStream != nil {
		if len(opts.Params) > 1 {
			errs = multierror.Append(errs, errors.Error("when streaming, there can be only one parameter", method.Position))
		}
		// if opts.InStream.Chan.Dir != types.RecvOnly {
		// errs = multierror.Append(errs, errors.Error("the channel for the request must be read-only", method.Position))
		// }
	}

	for _, result := range method.Sig.Results {
		if result.IsError {
			if result.Name == "" {
				errs = multierror.Append(errs, errors.Error("the result parameter name cannot be empty", result.Position))
			}
			if opts.Error != nil {
				errs = multierror.Append(errs, errors.Error("the result method has more than one error", result.Position))
			}
			opts.Error = result
			continue
		}
		if result.Name == "" {
			errs = multierror.Append(errs, errors.Error("the parameter name cannot be empty or the grpc-name parameter must be set", result.Position))
		}
		r, err := resultDecode(result)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		if result.IsChan {
			r.IsStream = true
			opts.OutStream = &OutStream{
				Param: r,
				Chan:  r.Type.(*types.Chan),
			}
		}
		opts.Results = append(opts.Results, r)
	}
	if opts.Error == nil {
		errs = multierror.Append(errs, errors.Error("the return of an error by the method is required", method.Position))
	}
	return
}

func paramDecode(param *types.Var) (opts *EndpointParam, err error) {
	opts = &EndpointParam{
		Title:           param.Title,
		FldName:         strcase.ToCamel(param.Name),
		FldNameUnExport: strcase.ToLowerCamel(param.Name),
		IsVariadic:      param.IsVariadic,
		IsPointer:       param.IsPointer,
		Type:            param.Type,
		Zero:            param.Zero,
	}
	tagFmt := "lowerCamel"
	if opts.Format != "" {
		tagFmt = opts.Format
	}
	name := formatName(param.Name, tagFmt)
	if opts.Name != "" {
		name = opts.Name
	}
	opts.Name = name

	opts.Format = "lowerCamel"
	if t, ok := param.Tags.Get("grpc-name"); ok {
		opts.Name = t.Value
		for _, option := range t.Options {
			if option == "omitempty" {
				opts.IsOmitempty = true
			}
			if v, ok := t.Param("format"); ok {
				opts.Format = v
			}
		}
	}
	if t, ok := param.Tags.Get("grpc-version"); ok {
		opts.Version = t.Value
	}
	if _, ok := param.Tags.Get("grpc-required"); ok {
		opts.IsRequired = true
	}
	return
}

func resultDecode(result *types.Var) (opts *EndpointResult, err error) {
	opts = &EndpointResult{
		Type:            result.Type,
		Format:          "lowerCamel",
		FldName:         result.Name,
		FldNameExport:   strcase.ToCamel(result.Name),
		FldNameUnExport: strcase.ToLowerCamel(result.Name),
		IsPointer:       result.IsPointer,
	}

	if t, ok := result.Tags.Get("grpc-name"); ok {
		opts.Name = t.Value
		for _, option := range t.Options {
			if option == "omitempty" {
				opts.IsOmitempty = true
			}
			if v, ok := t.Param("format"); ok {
				opts.Format = v
			}
		}
	}
	if opts.Name == "" {
		tagFmt := opts.Format
		if tagFmt == "" {
			tagFmt = "lowerCamel"
		}
		opts.Name = formatName(result.Name, tagFmt)
	}
	if t, ok := result.Tags.Get("grpc-version"); ok {
		opts.Version = t.Value
	}
	return
}
