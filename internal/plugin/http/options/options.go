package options

import (
	"fmt"
	"go/token"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"

	"github.com/hashicorp/go-multierror"
)

var (
	fnRegex = regexp.MustCompile(`^([A-Za-z0-9_]+)\(\).+$`)
)

type ErrorWrapperField struct {
	FldName    string
	FldType    any
	Name       string
	Interface  string
	MethodName string
}

type ErrorWrapper struct {
	Struct  *gg.Struct
	Default *gg.Struct
	Fields  []ErrorWrapperField
}

type Iface struct {
	Name         string
	Title        string
	Description  string
	PkgPath      string
	Server       Server
	Client       Client
	Openapi      OpenAPI
	APIDoc       APIDoc
	Endpoints    []Endpoint
	Type         string
	Lib          string
	ErrorWrapper ErrorWrapper
	HTTPReq      string
}

type Server struct {
	Enable bool
	Errors []string
}

type Client struct {
	Enable bool
}

type OpenAPI struct {
	Enable  bool
	Tags    []string
	Headers []OpenapiHeader
}

type APIDoc struct {
	Enable bool
}

type QueryValue struct {
	Name  string
	Value string
}

type Endpoint struct {
	Name               string
	MethodName         string
	Title              string
	Description        string
	HTTPMethod         string
	Path               string
	ParamsIdxName      map[string]int
	ParamsNameIdx      []string
	MultipartMaxMemory int64
	ReqRootXMLName     string
	RespRootXMLName    string
	ContentTypes       []string
	AcceptTypes        []string
	OpenapiTags        []string
	OpenapiHeaders     []OpenapiHeader
	QueryValues        []QueryValue
	WrapResponse       []string
	Params             []*EndpointParam
	BodyParams         []*EndpointParam
	QueryParams        []*EndpointParam
	HeaderParams       []*EndpointParam
	CookieParams       []*EndpointParam
	PathParams         []*EndpointParam
	Results            []*EndpointResult
	BodyResults        []*EndpointResult
	HeaderResults      []*EndpointResult
	CookieResults      []*EndpointResult
	Errors             []string
	NoWrapResponse     bool
	TimeFormat         string
	Context            *types.Var
	Error              *types.Var
	Sig                *types.Sign
}

type EndpointParam struct {
	Parent          *EndpointParam
	Type            any
	HTTPType        string
	Title           string
	Name            string
	FldName         string
	FldNameUnExport string
	Format          string
	Omitempty       bool
	IsVariadic      bool
	Required        bool
	Zero            string
	Flat            bool
	Params          []*EndpointParam
}

type EndpointResult struct {
	Type            any
	HTTPType        string
	Title           string
	Name            string
	FldName         string
	FldNameExport   string
	FldNameUnExport string
	Format          string
	Omitempty       bool
}

type OpenapiHeader struct {
	Title    string
	Name     string
	Required bool
}

func DecodeErrorWrapper(errorWrapperPath, defaultErrorPath string, structs []*gg.Struct) (errorWrapper *ErrorWrapper, errs error) {
	var (
		errorWrapperStruct *gg.Struct
		defaultErrorStruct *gg.Struct
	)
	for _, s := range structs {
		if errorWrapperStruct != nil && defaultErrorStruct != nil {
			break
		}
		if errorWrapperPath == path.Join(s.Named.Pkg.Path, s.Named.Name) {
			errorWrapperStruct = s
		} else if defaultErrorPath == path.Join(s.Named.Pkg.Path, s.Named.Name) {
			defaultErrorStruct = s
		}
	}
	if errorWrapperStruct == nil {
		errs = multierror.Append(errs, errors.Error("not found error wrapper struct "+errorWrapperPath, token.Position{}))
		return
	}
	if defaultErrorStruct == nil {
		errs = multierror.Append(errs, errors.Error("not found default error struct "+defaultErrorPath, token.Position{}))
		return
	}
	errorWrapper = &ErrorWrapper{
		Struct:  errorWrapperStruct,
		Default: defaultErrorStruct,
	}
	if errorWrapperStruct != nil {
		for _, field := range errorWrapperStruct.Type.Fields {
			if t, ok := field.Var.Tags.Get("http-error-interface"); ok {
				name := strcase.ToLowerCamel(field.Var.Name)
				if jsonTag, err := field.SysTags.Get("json"); err == nil {
					name = jsonTag.Name
				}

				var methodName string
				matches := fnRegex.FindAllStringSubmatch(t.Value, -1)
				if len(matches) > 0 && len(matches[0]) == 2 {
					methodName = matches[0][1]
				} else {
					errs = multierror.Append(errs, errors.Error("invalid interface method", t.Position))
					continue
				}
				errorWrapper.Fields = append(errorWrapper.Fields, ErrorWrapperField{
					FldName:    field.Var.Name,
					FldType:    field.Var.Type,
					Name:       name,
					Interface:  t.Value,
					MethodName: methodName,
				})
			}
		}
	}
	return
}

func Decode(iface *gg.Interface) (opts Iface, errs error) {
	opts.Name = iface.Named.Name
	opts.Title = iface.Named.Title
	opts.Description = iface.Named.Description
	opts.PkgPath = iface.Named.Pkg.Path
	opts.Type = "rest"
	opts.Lib = "echo"
	if t, ok := iface.Named.Tags.Get("http-type"); ok {
		switch t.Value {
		default:
			errs = multierror.Append(errs, errors.Error("invalid http type, valid values rest, jsonrpc", t.Position))
		case "rest", "jsonrpc":
			opts.Type = t.Value
		}
		if len(t.Options) > 0 {
			opts.Lib = t.Options[0]
		}
	}
	if _, ok := iface.Named.Tags.Get("http-server"); ok {
		opts.Server.Enable = true
	}
	if _, ok := iface.Named.Tags.Get("http-client"); ok {
		opts.Client.Enable = true
	}
	if _, ok := iface.Named.Tags.Get("http-openapi"); ok {
		opts.Openapi.Enable = true
	}
	if t, ok := iface.Named.Tags.Get("http-openapi-tags"); ok {
		opts.Openapi.Tags = []string{t.Value}
		opts.Openapi.Tags = append(opts.Openapi.Tags, t.Options...)
	}
	openapiHeaderTags := iface.Named.Tags.GetSlice("http-openapi-header")
	for _, t := range openapiHeaderTags {
		title, _ := t.Param("title")
		oh := OpenapiHeader{
			Title: title,
			Name:  t.Value,
		}
		for _, v := range t.Options {
			switch v {
			case "required":
				oh.Required = true
			}
		}
		opts.Openapi.Headers = append(opts.Openapi.Headers, oh)
	}
	if _, ok := iface.Named.Tags.Get("http-api-doc"); ok {
		opts.APIDoc.Enable = true
	}
	errorTags := iface.Named.Tags.GetSlice("http-error")
	for _, tag := range errorTags {
		opts.Server.Errors = append(opts.Server.Errors, tag.Value)
	}
	if t, ok := iface.Named.Tags.Get("http-req"); ok {
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
	opts.TimeFormat = time.RFC3339
	opts.Sig = method.Sig

	if t, ok := method.Tags.Get("http-time-format"); ok {
		opts.TimeFormat = t.Value
	}
	if ifaceOpts.Type == "rest" {
		if t, ok := method.Tags.Get("http-method"); ok {
			switch t.Value {
			default:
				errs = multierror.Append(errs, errors.Error("invalid http method, valid values GET, HEAD, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH", t.Position))
			case "GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH":
				opts.HTTPMethod = t.Value
			}
		} else {
			errs = multierror.Append(errs, errors.Error("the http-method parameter is required", method.Position))
		}
	} else {
		opts.HTTPMethod = "POST"
	}
	if t, ok := method.Tags.Get("http-path"); ok {
		if _, err := url.Parse(t.Value); err != nil {
			errs = multierror.Append(errs, errors.Error("invalid http-path format", t.Position))
		}
		opts.Path = t.Value

	}

	if ifaceOpts.Type == "rest" {
		parts := strings.Split(opts.Path, "/")
		opts.ParamsIdxName = make(map[string]int, len(parts))
		for _, part := range parts {
			if strings.HasPrefix(part, ":") {
				paramName := part[1:]
				opts.ParamsNameIdx = append(opts.ParamsNameIdx, paramName)
				opts.ParamsIdxName[paramName] = len(opts.ParamsNameIdx) - 1
			}
		}
	}
	if opts.Path == "" && ifaceOpts.Type == "jsonrpc" {
		opts.Path = strcase.ToLowerCamel(ifaceOpts.Name) + "." + strcase.ToLowerCamel(method.Name)
	}
	if t, ok := method.Tags.Get("http-openapi-tags"); ok {
		opts.OpenapiTags = []string{t.Value}
		opts.OpenapiTags = append(opts.OpenapiTags, t.Options...)
	}
	openapiHeaderTags := method.Tags.GetSlice("http-openapi-header")
	for _, t := range openapiHeaderTags {
		title, _ := t.Param("title")
		oh := OpenapiHeader{
			Title: title,
			Name:  t.Value,
		}
		for _, v := range t.Options {
			switch v {
			case "required":
				oh.Required = true
			}
		}
		opts.OpenapiHeaders = append(opts.OpenapiHeaders, oh)
	}
	if t, ok := method.Tags.Get("http-content-types"); ok {
		opts.ContentTypes = []string{t.Value}
		opts.ContentTypes = append(opts.ContentTypes, t.Options...)

		for _, contentType := range opts.ContentTypes {
			switch contentType {
			default:
				errs = multierror.Append(errs, errors.Error("invalid http-content-types use 'json', 'xml', 'urlencoded' or 'multipart'", t.Position))
			case "json", "xml", "urlencoded", "multipart":
				switch contentType {
				case "xml":
					opts.ReqRootXMLName = t.Params["root-xml"]
					if opts.ReqRootXMLName == "" {
						errs = multierror.Append(errs, errors.Error("the root-xml parameter of the http-content-types tag is required when using the XML content type", t.Position))
					}
				case "multipart":
					opts.MultipartMaxMemory = 67108864
					if t.Params["multipart-max-memory"] != "" {
						multipartMaxMemory, err := strconv.ParseInt(t.Params["multipart-max-memory"], 10, 64)
						if err == nil {
							opts.MultipartMaxMemory = multipartMaxMemory
						} else {
							errs = multierror.Append(errs, errors.Error("invalid multipart-max-memory must be integer", t.Position))
						}
					} else {
						errs = multierror.Append(errs, errors.Warn(fmt.Sprintf("multipartMaxMemory uses the default value of %d bytes", opts.MultipartMaxMemory), t.Position))
					}
				}
			}
		}
	}
	if t, ok := method.Tags.Get("http-accept-types"); ok {
		switch t.Value {
		default:
			errs = multierror.Append(errs, errors.Error("invalid http-accept-types, use 'json', 'xml', 'urlencoded' or 'multipart'", t.Position))
		case "json", "xml", "urlencoded", "multipart":
			opts.RespRootXMLName = t.Params["root-xml"]
			opts.AcceptTypes = append([]string{t.Value}, t.Options...)
		}
	}
	queryValues := method.Tags.GetSlice("http-query-value")
	for _, q := range queryValues {
		var val string
		if len(q.Options) > 0 {
			val = q.Options[0]
		}
		opts.QueryValues = append(opts.QueryValues, QueryValue{
			Name:  q.Value,
			Value: val,
		})
	}
	if t, ok := method.Tags.Get("http-wrap-response"); ok {
		opts.WrapResponse = strings.Split(t.Value, ".")
	}
	if _, ok := method.Tags.Get("http-nowrap-response"); ok {
		opts.NoWrapResponse = true
	}
	if len(opts.WrapResponse) > 0 && opts.NoWrapResponse {
		errs = multierror.Append(errs, errors.Warn("the http-wrap-response tag conflicts with http-nowrap-response", method.Position))
	}
	errorTags := method.Tags.GetSlice("http-error")
	for _, tag := range errorTags {
		opts.Errors = append(opts.Errors, tag.Value)
	}
	if len(opts.ContentTypes) == 0 {
		opts.ContentTypes = []string{"json"}
	}
	if len(opts.AcceptTypes) == 0 {
		opts.AcceptTypes = []string{"json"}
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
			errs = multierror.Append(errs, errors.Error("the parameter name cannot be empty or the http-name parameter must be set", param.Position))
		}

		p, err := makeEndpointParam(nil, param)
		if err != nil {
			errs = multierror.Append(errs, err)
		}

		if _, ok := opts.ParamsIdxName[param.Name]; ok {
			p.HTTPType = "path"
			p.Name = param.Name
			p.Required = true
		}

		opts.Params = append(opts.Params, p)
	}
	for _, result := range method.Sig.Results {
		if result.IsError {
			if result.Name == "" {
				errs = multierror.Append(errs, errors.Error("the parameter name cannot be empty or the http-name parameter must be set", result.Position))
			}
			if opts.Error != nil {
				errs = multierror.Append(errs, errors.Error("the method has more than one error", result.Position))
			}
			opts.Error = result
			continue
		}
		if result.Name == "" {
			errs = multierror.Append(errs, errors.Error("the parameter name cannot be empty or the http-name parameter must be set", result.Position))
		}
		if name, ok := checkBasicType(result.Type); ok {
			if opts.NoWrapResponse {
				errs = multierror.Append(errs, errors.Error("the \"@http-nowrap-response\" tag cannot be used for basic type "+name, result.Position))
			}
			if name == "[]byte" {
				errs = multierror.Append(errs, errors.Error("the []byte type is not supported, use marshaling and unmarshalling of non-standard formats for the response", result.Position))
			}
		}

		varOpts, err := resultDecode(result)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		opts.Results = append(opts.Results, &varOpts)
	}

	var fillTypeParams func(params []*EndpointParam)
	fillTypeParams = func(params []*EndpointParam) {
		for _, param := range params {
			if len(param.Params) > 0 {
				fillTypeParams(param.Params)
				continue
			}
			switch param.HTTPType {
			case "path":
				opts.PathParams = append(opts.PathParams, param)
			case "cookie":
				opts.CookieParams = append(opts.CookieParams, param)
			case "query":
				opts.QueryParams = append(opts.QueryParams, param)
			case "header":
				opts.HeaderParams = append(opts.HeaderParams, param)
			case "body":
				opts.BodyParams = append(opts.BodyParams, param)
			}
		}
	}

	fillTypeParams(opts.Params)

	for _, result := range opts.Results {
		switch result.HTTPType {
		case "cookie":
			opts.CookieResults = append(opts.CookieResults, result)
		case "header":
			opts.HeaderResults = append(opts.HeaderResults, result)
		case "body":
			opts.BodyResults = append(opts.BodyResults, result)
		}
	}
	if len(opts.BodyParams) > 0 && (opts.HTTPMethod != "POST" && opts.HTTPMethod != "PUT" && opts.HTTPMethod != "DELETE" && opts.HTTPMethod != "PATCH") {
		errs = multierror.Append(errs, errors.Error("only HTTP POST, PUT, PATCH and DELETE methods can have a request body", method.Position))
	}
	if len(opts.PathParams) != len(opts.ParamsNameIdx) {
		errs = multierror.Append(errs, errors.Error("the method has no parameters found for the http-path tag, the required parameters: "+strings.Join(opts.ParamsNameIdx, ", "), method.Position))
	}
	if opts.NoWrapResponse && len(opts.Results) != 1 {
		errs = multierror.Append(errs, errors.Error("the \"@http-nowrap-response\" tag can be used for only one return parameter", method.Position))
	}
	return
}

func paramDecode(param *types.Var) (opts EndpointParam, err error) {
	opts.HTTPType = "body"
	opts.Format = "lowerCamel"
	if t, ok := param.Tags.Get("http-name"); ok {
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
	if t, ok := param.Tags.Get("http-type"); ok {
		opts.HTTPType = t.Value
	}
	if _, ok := param.Tags.Get("http-required"); ok {
		opts.Required = true
	}
	if _, ok := param.Tags.Get("http-flat"); ok {
		opts.Flat = true
	}
	return
}

func resultDecode(result *types.Var) (opts EndpointResult, err error) {
	opts.HTTPType = "body"
	opts.Type = result.Type
	opts.Format = "lowerCamel"
	opts.FldName = result.Name
	opts.FldNameExport = strcase.ToCamel(result.Name)
	opts.FldNameUnExport = strcase.ToLowerCamel(result.Name)
	if t, ok := result.Tags.Get("http-name"); ok {
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
	if t, ok := result.Tags.Get("http-type"); ok {
		opts.HTTPType = t.Value
	}
	return
}
