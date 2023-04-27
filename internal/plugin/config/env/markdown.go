package env

import (
	"github.com/555f/gg/internal/plugin/config/options"
	"github.com/555f/gg/pkg/file"
)

func GenMarkdownDoc(c options.Config) func(f *file.TxtFile) {
	return func(f *file.TxtFile) {
		markdownTitle := "Environment Variables"
		if c.MarkdownTitle != "" {
			markdownTitle = c.MarkdownTitle
		}

		f.WriteText("# %s\n\n", markdownTitle)

		if c.MarkdownDescription != "" {
			f.WriteText("%s\n\n", c.MarkdownDescription)
		}

		f.WriteText("| Name | Type | Description | Required | Use Zero |\n|------|------|------|------|------|\n")

		walkFields(c.Fields, func(parent *options.ConfigField, field options.ConfigField) {
			var envName string
			pathFields := resolvePathFields(parent)
			for i := len(pathFields) - 1; i >= 0; i-- {
				envName += pathFields[i].Name + "_"
			}
			envName += field.Name
			required := "no"
			if field.Required {
				required = "yes"
			}
			useZero := "no"
			if field.UseZero {
				useZero = "yes"
			}
			f.WriteText("|%s|<code>%s</code>|%s|%s|%s|\n", envName, getTypeSrt(field.Type), field.Description+" ", required, useZero)
		})
	}
}
