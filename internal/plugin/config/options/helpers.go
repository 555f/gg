package options

import "github.com/555f/gg/pkg/types"

func extractFields(v any) []*types.StructFieldType {
	switch t := v.(type) {
	default:
		return nil
	case *types.Struct:
		return t.Fields
	case *types.Named:
		return extractFields(t.Type)
	}
}
