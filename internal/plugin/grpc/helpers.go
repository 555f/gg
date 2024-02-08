package grpc

import (
	"github.com/555f/gg/pkg/types"

	stdtypes "go/types"
)

func goType2GRPC(v any) string {
	switch t := v.(type) {
	case *types.Basic:
		switch t.Kind {
		case stdtypes.Float64:
			return "double"
		case stdtypes.Float32:
			return "float"
		case stdtypes.Int8, stdtypes.Int16, stdtypes.Int32:
			return "int32"
		case stdtypes.Int, stdtypes.Int64:
			return "int64"
		case stdtypes.Uint8, stdtypes.Uint16, stdtypes.Uint32:
			return "uint32"
		case stdtypes.Uint, stdtypes.Uint64:
			return "uint64"
		case stdtypes.Bool:
			return "bool"
		case stdtypes.String:
			return "string"
		}
	case *types.Named:
		return t.Name
	case *types.Array:
		if b, ok := t.Value.(*types.Basic); ok && b.IsByte() {
			return "bytes"
		}
		return "repeated " + goType2GRPC(t.Value)
	}
	panic("unknown go type")
}
