package metrics

import (
	"github.com/555f/gg/pkg/types"
)

type methodOptions struct {
	Skip bool
}

func makeMethodOptions(module *types.Module, method *types.Func) (opts methodOptions, err error) {
	if _, ok := method.Tags.Get("metrics-skip"); ok {
		opts.Skip = true
	}
	return
}
