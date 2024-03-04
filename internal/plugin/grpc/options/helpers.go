package options

import (
	"github.com/555f/gg/pkg/strcase"
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
