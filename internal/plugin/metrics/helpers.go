package metrics

import (
	"github.com/555f/gg/pkg/types"
)

func shortMethodName(f *types.Func) string {
	if named, ok := f.Sig.Recv.(*types.Named); ok {
		return "(" + named.Name + ")." + f.Name
	}
	return f.Name
}
