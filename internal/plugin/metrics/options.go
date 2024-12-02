package metrics

import (
	"github.com/555f/gg/pkg/types"
	"github.com/hashicorp/go-multierror"
)

type importPath struct {
	PkgPath string
	Name    string
	Valid   bool
}

type methodOptions struct {
	Skip         bool
	LabelErrFunc importPath
}

func makeMethodOptions(module *types.Module, method *types.Func) (opts methodOptions, errs error) {
	if _, ok := method.Tags.Get("metrics-skip"); ok {
		opts.Skip = true
	}
	if t, ok := method.Tags.Get("metrics-label-err-func"); ok {
		pkgPath, name, err := module.ParseImportPath(t.Value)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		opts.LabelErrFunc = importPath{
			PkgPath: pkgPath,
			Name:    name,
			Valid:   true,
		}
	}
	return opts, errs
}
