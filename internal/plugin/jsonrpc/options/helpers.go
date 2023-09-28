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

func makeEndpointParam(
	parent *EndpointParam,
	param *types.Var,
) (epParam *EndpointParam, errs error) {
	opts, err := paramDecode(param)
	if err != nil {
		errs = multierror.Append(errs, err)
		return
	}
	epParam = &EndpointParam{
		Title:           param.Title,
		FldName:         strcase.ToCamel(param.Name),
		FldNameUnExport: strcase.ToLowerCamel(param.Name),
		IsVariadic:      param.IsVariadic,
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
	epParam.Name = name
	return
}
