package options

import (
	"strings"

	"github.com/555f/gg/pkg/gen"

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
	EnvsFile            string
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
		config.EnvsFile = t.Value
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

func DecodeField(parent *ConfigField, f *types.Var) (cf ConfigField, errs error) {
	cf.Parent = parent
	cf.FieldName = f.Name
	cf.Type = f.Type
	cf.Name = f.Name
	cf.Description = strings.ReplaceAll(strings.Trim(f.Title, "\n"), "\n", "<br/>")

	if t, ok := f.Tags.Get("cfg-name"); ok {
		cf.Name = t.Value
	}
	if _, ok := f.Tags.Get("cfg-required"); ok {
		cf.Required = true
	}
	if t, ok := f.Tags.Get("cfg-dev-value"); ok {
		cf.DevValue = t.Value
	}
	if _, ok := f.Tags.Get("cfg-use-zero"); ok {
		cf.UseZero = true
	}

	cf.Name = strings.ToUpper(strcase.ToScreamingSnake(cf.Name))
	cf.Zero = types.ZeroValueJS(f.Type)

	for _, structField := range gen.ExtractFields(f.Type) {
		switch t := structField.Type.(type) {
		case *types.Named:
			switch t.Pkg.Path {
			case "net/url":
				if t.Name == "URL" && !t.IsPointer {
					errs = multierror.Append(errs, errors.Error("invalid net/url.URL type, there must be a pointer", structField.Position))
				}
			}
		}

		pf, err := DecodeField(&cf, structField)
		if err != nil {
			return ConfigField{}, errors.Error(err.Error(), structField.Position)
		}
		cf.Fields = append(cf.Fields, pf)
	}
	return
}
