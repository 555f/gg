package grpc

import (
	"fmt"

	"github.com/555f/gg/internal/plugin/grpc/options"
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
		switch t.Pkg.Name {
		case "time":
			switch t.Name {
			case "Duration":
				return "google.protobuf.Duration"
			case "Time":
				return "google.protobuf.Timestamp"
			}
		}
		return t.Name
	case *types.Slice:
		if b, ok := t.Value.(*types.Basic); ok && b.IsByte() {
			return "bytes"
		}
		return "repeated " + goType2GRPC(t.Value)
	case *types.Array:
		if b, ok := t.Value.(*types.Basic); ok && b.IsByte() {
			return "bytes"
		}
		return "repeated " + goType2GRPC(t.Value)
	}
	panic(fmt.Sprintf("unknown go type: %T", v))
}

func hasResponseEndpoint(ep options.Endpoint) bool {
	return (ep.InStream == nil && ep.OutStream == nil)
}
