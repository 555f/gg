package options

import (
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/hashicorp/go-multierror"
)

var nameFormatters = map[string]func(string) string{
	"lowerCamel":     strcase.ToLowerCamel,
	"kebab":          strcase.ToKebab,
	"screamingKebab": strcase.ToScreamingKebab,
	"snake":          strcase.ToSnake,
	"screamingSnake": strcase.ToScreamingSnake,
}

func formatName(s, fmt string) string {
	if f, ok := nameFormatters[fmt]; ok {
		return f(s)
	}
	return s
}

func nameStructResp(ifaceName string, methodName string) string {
	return strcase.ToLowerCamel(ifaceName + methodName + "Response")
}

func nameStructReq(ifaceName string, methodName string) string {
	return strcase.ToLowerCamel(ifaceName + methodName + "Request")
}

func nameReqDecode(ifaceName string, methodName string) string {
	return strcase.ToLowerCamel(ifaceName + methodName + "ReqDecode")
}

func makeEndpointParam(
	parent *EndpointParam,
	param *types.Var,
	paramsIdxName map[string]int,
) (epParam *EndpointParam, errs error) {
	opts, err := paramDecode(param)
	if err != nil {
		errs = multierror.Append(errs, err)
		return
	}
	epParam = &EndpointParam{
		Title:           param.Title,
		Parent:          parent,
		FldName:         strcase.ToCamel(param.Name),
		FldNameUnExport: strcase.ToLowerCamel(param.Name),
		IsVariadic:      param.IsVariadic,
		Type:            param.Type,
		Required:        opts.Required,
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
	paramType := opts.HTTPType
	if _, ok := paramsIdxName[param.Name]; ok {
		paramType = "path"
		name = param.Name
		epParam.Required = true
	}
	epParam.HTTPType = paramType
	epParam.Name = name

	if opts.Flat {
		if named, ok := param.Type.(*types.Named); ok {
			if st, ok := named.Type.(*types.Struct); ok {
				for _, field := range st.Fields {
					p, err := makeEndpointParam(epParam, field.Var, paramsIdxName)
					if err != nil {
						errs = multierror.Append(errs, err)
						continue
					}
					epParam.Params = append(epParam.Params, p)
				}
			}
		}
	}
	return
}

func checkBasicType(t any) (name string, ok bool) {
	switch t := t.(type) {
	default:
		return "", false
	case *types.Array:
		name, ok = checkBasicType(t.Value)
		if ok {
			name = "[]" + name
		}
		return
	case *types.Slice:
		name, ok = checkBasicType(t.Value)
		if ok {
			name = "[]" + name
		}
		return
	case *types.Basic:
		return t.Name, true
	}
}
