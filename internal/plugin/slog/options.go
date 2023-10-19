package logging

import (
	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/types"
)

type logContext struct {
	Name    string
	LogName string
	PkgPath string
	Type    string
}

type methodOptions struct {
	Skip        bool
	LogContexts []logContext
}

func makeMethodOptions(module *types.Module, method *types.Func) (opts methodOptions, err error) {
	if _, ok := method.Tags.Get("logging-skip"); ok {
		opts.Skip = true
	}
	if t, ok := method.Tags.Get("logging-context"); ok {
		if len(t.Options) == 0 {
			err = errors.Error("the path to the context key is required", t.Position)
			return
		}
		pkgPath, name, err := module.ParseImportPath(t.Options[0])
		if err != nil {
			return methodOptions{}, errors.Error(err.Error(), t.Position)
		}
		opts.LogContexts = append(opts.LogContexts, logContext{
			Name:    name,
			LogName: t.Value,
			PkgPath: pkgPath,
		})
	}
	return
}

type paramOptions struct {
	Name    string
	Skip    bool
	Context logContext
}

// type resultOptions struct {
// 	Name string
// 	Skip bool
// }

func makeParamOptions(pkg *types.PackageType, tags types.Tags) (opts paramOptions, err error) {
	if _, ok := tags.Get("logging-skip"); ok {
		opts.Skip = true
	}
	if t, ok := tags.Get("logging-name"); ok {
		if t.Value == "" {
			err = errors.Error(t.Key+": the value cannot be empty", t.Position)
			return
		}
		opts.Name = t.Value
	}
	return
}

// func makeResultOptions(tags types.Tags) (opts resultOptions, err error) {
// 	if _, ok := tags.Get("logging-skip"); ok {
// 		opts.Skip = true
// 	}
// 	if t, ok := tags.Get("logging-name"); ok {
// 		if t.Value == "" {
// 			err = errors.Error(t.Key+": the value cannot be empty", t.Position)
// 			return
// 		}
// 		opts.Name = t.Value
// 	}
// 	return
// }
