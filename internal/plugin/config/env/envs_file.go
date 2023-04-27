package env

import (
	"github.com/555f/gg/internal/plugin/config/options"
	"github.com/555f/gg/pkg/file"
)

func GenEnvsFile(c options.Config) func(f *file.TxtFile) {
	return func(f *file.TxtFile) {
		walkFields(c.Fields, func(parent *options.ConfigField, field options.ConfigField) {
			var envName string
			pathFields := resolvePathFields(parent)
			for i := len(pathFields) - 1; i >= 0; i-- {
				envName += pathFields[i].Name + "_"
			}
			envName += field.Name
			f.WriteText("%s=%s\n", envName, field.DevValue)
		})
	}
}
