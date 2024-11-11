package env

import (
	"github.com/555f/gg/internal/plugin/config/options"
	"github.com/555f/gg/pkg/types"
)

func resolvePathFields(parent *options.ConfigField) (fields []options.ConfigField) {
	for parent != nil {
		fields = append(fields, *parent)
		parent = parent.Parent
	}
	return
}

func walkFields(fields []options.ConfigField, cb func(parent *options.ConfigField, field options.ConfigField)) {
	walkRecursiveFields(nil, fields, cb)
}

func walkRecursiveFields(parent *options.ConfigField, fields []options.ConfigField, cb func(parent *options.ConfigField, field options.ConfigField)) {
	for _, field := range fields {
		if len(field.Fields) > 0 {
			walkRecursiveFields(&field, field.Fields, cb)
			continue
		}
		cb(parent, field)
	}
}

func getTypeSrt(t any) string {
	switch t := t.(type) {
	default:
		return ""
	case *types.Map:
		return "map[string]" + getTypeSrt(t.Value)
	case *types.Array:
		return getTypeSrt(t.Value) + "[]"
	case *types.Slice:
		return getTypeSrt(t.Value) + "[]"
	case *types.Basic:
		return t.Name()
	}
}
