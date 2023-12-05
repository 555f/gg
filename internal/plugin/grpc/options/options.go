package options

import (
	"bytes"
	"strconv"

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

func paramJSONFromType(name string, t any) string {
	var buf bytes.Buffer

	buf.WriteString("  " + strconv.Quote(name) + ":")

	switch v := t.(type) {
	case *types.Named:
		name := v.Pkg.Path + "." + v.Name
		switch name {
		default:
			if st := v.Struct(); st != nil {
				buf.WriteString("{\n")
				for _, f := range v.Struct().Fields {
					name := f.Var.Name
					if t, err := f.SysTags.Get("json"); err == nil {
						name = t.Value()
					}
					buf.WriteString("  " + paramJSONFromType(name, f.Var.Type))
				}
				buf.WriteString("  }")
			} else {
				buf.WriteString("\"\"")
			}
		case "time.Time":
			buf.WriteString("\"\"")
		}

	case *types.Slice, *types.Array:
		buf.WriteString("[]")
	case *types.Basic:
		if v.IsNumeric() {
			buf.WriteString("0")
		} else {
			buf.WriteString("\"\"")
		}
	}
	return buf.String()
}

func (params EndpointParams) ToJSON() string {
	var buf bytes.Buffer
	buf.WriteString("{\n")
	for i, p := range params {
		if i > 0 {
			buf.WriteString(",\n")
		}
		buf.WriteString(paramJSONFromType(p.Name, p.Type))
	}
	buf.WriteString("}\n")
	return buf.String()
}

type Endpoint struct {
	Name          string
	MethodName    string
	RPCMethodName string
	Title         string
	Description   string
	OpenapiTags   []string
	Params        EndpointParams
	Results       []*EndpointResult
	Context       *types.Var
	Error         *types.Var
	Sig           *types.Sign
}

type EndpointParam struct {
	Type            any
	Title           string
	Name            string
	FldName         string
	FldNameUnExport string
	Format          string
	Omitempty       bool
	IsVariadic      bool
	Zero            string
	Required        bool
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
	Omitempty       bool
}

func Decode(iface *gg.Interface) (opts Iface, errs error) {
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
		epOpts, err := endpointDecode(opts, method)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		opts.Endpoints = append(opts.Endpoints, epOpts)
	}
	return
}

func endpointDecode(ifaceOpts Iface, method *types.Func) (opts Endpoint, errs error) {
	opts.Name = strcase.ToLowerCamel(ifaceOpts.Name) + method.Name + "Endpoint"
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
		p, err := makeEndpointParam(nil, param)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		opts.Params = append(opts.Params, p)
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
		varOpts, err := resultDecode(result)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		opts.Results = append(opts.Results, &varOpts)
	}
	if opts.Error == nil {
		errs = multierror.Append(errs, errors.Error("the return of an error by the method is required", method.Position))
	}
	return
}

func paramDecode(param *types.Var) (opts EndpointParam, err error) {
	opts.Format = "lowerCamel"
	if t, ok := param.Tags.Get("grpc-name"); ok {
		opts.Name = t.Value
		for _, option := range t.Options {
			if option == "omitempty" {
				opts.Omitempty = true
			}
			if v, ok := t.Param("format"); ok {
				opts.Format = v
			}
		}
	}
	if _, ok := param.Tags.Get("grpc-required"); ok {
		opts.Required = true
	}
	return
}

func resultDecode(result *types.Var) (opts EndpointResult, err error) {
	opts.Type = result.Type
	opts.Format = "lowerCamel"
	opts.FldName = result.Name
	opts.FldNameExport = strcase.ToCamel(result.Name)
	opts.FldNameUnExport = strcase.ToLowerCamel(result.Name)
	if t, ok := result.Tags.Get("grpc-name"); ok {
		opts.Name = t.Value
		for _, option := range t.Options {
			if option == "omitempty" {
				opts.Omitempty = true
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
	return
}
