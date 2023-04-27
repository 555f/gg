package options

import (
	"strings"

	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/hashicorp/go-multierror"
)

type Config struct {
	Name                string
	ConstructName       string
	PkgPath             string
	Fields              []ConfigField
	MarkdownDoc         bool
	MarkdownTitle       string
	MarkdownDescription string
	MarkdownEnvFile     string
}

type ConfigField struct {
	Parent      *ConfigField
	Name        string
	Description string
	Required    bool
	UseZero     bool
	FieldName   string
	Fields      []ConfigField
	Type        any
	Zero        string
	DevValue    string
}

func Decode(st *gg.Struct) (config Config, errs error) {
	config.Name = st.Named.Name
	config.PkgPath = st.Named.Pkg.Path
	config.ConstructName = "New"
	config.MarkdownTitle = st.Named.Title
	config.MarkdownDescription = st.Named.Description

	if t, ok := st.Named.Tags.Get("cfg-constructor-name"); ok {
		config.ConstructName = t.Value
	}
	if _, ok := st.Named.Tags.Get("cfg-md-doc"); ok {
		config.MarkdownDoc = true
	}
	if t, ok := st.Named.Tags.Get("cfg-env-file"); ok {
		config.MarkdownEnvFile = t.Value
	}
	for _, field := range st.Type.Fields {
		cf, err := DecodeField(nil, field)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		config.Fields = append(config.Fields, cf)
	}
	return
}

func DecodeField(parent *ConfigField, f *types.StructFieldType) (cf ConfigField, errs error) {
	cf.Parent = parent
	cf.FieldName = f.Var.Name
	cf.Type = f.Var.Type
	cf.Name = f.Var.Name
	cf.Description = strings.ReplaceAll(strings.Trim(f.Var.Title, "\n"), "\n", "<br/>")

	if t, ok := f.Var.Tags.Get("cfg-name"); ok {
		cf.Name = t.Value
	}
	if _, ok := f.Var.Tags.Get("cfg-required"); ok {
		cf.Required = true
	}
	if t, ok := f.Var.Tags.Get("cfg-dev-value"); ok {
		cf.DevValue = t.Value
	}

	cf.Name = strings.ToUpper(strcase.ToScreamingSnake(f.Var.Name))
	cf.Zero = f.Var.Zero

	for _, structField := range extractFields(f.Var.Type) {
		pf, err := DecodeField(&cf, structField)
		if err != nil {
			return ConfigField{}, errors.Error(err.Error(), structField.Var.Position)
		}
		cf.Fields = append(cf.Fields, pf)
	}
	return
}
