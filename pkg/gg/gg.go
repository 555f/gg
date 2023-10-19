package gg

import (
	"fmt"
	"go/ast"
	"go/token"
	stdtypes "go/types"
	"sort"
	"strings"

	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/types"

	"github.com/hashicorp/go-multierror"
	stdpackages "golang.org/x/tools/go/packages"
)

func Run(version, wd string, packages []*stdpackages.Package, plugins map[string]any) (allFiles []file.File, errs error) {
	module, err := Module(packages)
	if err != nil {
		errs = multierror.Append(errs, errors.Error("the golang module was not found, see for more details https://go.dev/blog/using-go-modules", token.Position{}))
		return
	}
	structs, err := findStructs(packages)
	if err != nil {
		errs = multierror.Append(errs, err)
	}
	interfaces, err := findInterfaces(packages)
	if err != nil {
		errs = multierror.Append(errs, err)
	}

	sort.Slice(interfaces, func(i, j int) bool {
		return strings.Compare(interfaces[i].Named.Name, interfaces[j].Named.Name) > 0
	})
	sort.Slice(structs, func(i, j int) bool {
		return strings.Compare(structs[i].Named.Name, structs[j].Named.Name) > 0
	})

	interfaceSet := map[string][]*Interface{}
	structSet := map[string][]*Struct{}
	pluginUsesSet := map[string][]token.Position{}

	for _, s := range structs {
		if !s.Named.IsExported {
			if len(s.Named.Tags) > 0 {
				errs = multierror.Append(errs, errors.Warn("tags are defined for the structure, but it is not exportable", s.Named.Position))
			}
			continue
		}
		if len(s.Named.Tags) == 0 {
			continue
		}
		for _, t := range s.Named.Tags.GetSlice("gg") {
			pluginUsesSet[t.Value] = append(pluginUsesSet[t.Value], s.Named.Position)
			structSet[t.Value] = append(structSet[t.Value], s)
		}
	}

	for _, iface := range interfaces {
		if !iface.Named.IsExported {
			if len(iface.Named.Tags) > 0 {
				errs = multierror.Append(errs, errors.Warn("tags are defined for the interface, but it is not exportable", iface.Named.Position))
			}
			continue
		}
		if len(iface.Named.Tags) == 0 {
			continue
		}
		for _, t := range iface.Named.Tags.GetSlice("gg") {
			pluginUsesSet[t.Value] = append(pluginUsesSet[t.Value], iface.Named.Position)
			interfaceSet[t.Value] = append(interfaceSet[t.Value], iface)
		}
	}

	var pluginGraph = newGraph()

	pkgPath := module.Path + strings.Replace(wd, module.Dir, "", -1)

	for name, f := range pluginFactories {
		if len(interfaceSet[name]) > 0 || len(structSet[name]) > 0 {
			options, _ := plugins[name].(map[string]any)
			ctx := &Context{
				Version:     version,
				pluginGraph: pluginGraph,
				Workdir:     wd,
				PkgPath:     pkgPath,
				Module:      module,
				Interfaces:  interfaceSet[name],
				Structs:     structSet[name],
				Options:     Options{m: options},
			}
			plugin := f(ctx)
			if err := pluginGraph.add(plugin); err != nil {
				errs = multierror.Append(errs, err)
			}
		}
	}
	sortedPlugins := pluginGraph.Sorted()

	for name, positions := range pluginUsesSet {
		plugin, ok := pluginGraph.plugins[name]
		if !ok {
			for _, pos := range positions {
				errs = multierror.Append(errs, errors.Warn(fmt.Sprintf("plugin not found: %s", name), pos))
			}
			continue
		}
		for _, dep := range plugin.Dependencies() {
			if _, ok := pluginUsesSet[dep]; !ok {
				for _, pos := range positions {
					errs = multierror.Append(errs, errors.Error(fmt.Sprintf("%s depends on: %s, you need to add it", plugin.Name(), dep), pos))
				}
			}
		}
	}

	for i := len(sortedPlugins) - 1; i >= 0; i-- {
		plugin := sortedPlugins[i]
		files, err := plugin.Exec()
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		allFiles = append(allFiles, files...)
	}
	return
}

func findInterfaces(packages []*stdpackages.Package) (result Interfaces, err error) {
	err = findTypes(
		packages,
		func(tp stdtypes.Type) (ok bool) {
			_, ok = tp.(*stdtypes.Interface)
			return
		},
		func(namedType *types.Named) (err error) {
			iface := namedType.Interface()
			result = append(result, &Interface{
				Named: namedType,
				Type:  iface,
			})
			return
		},
	)
	if err != nil {
		return
	}
	return
}

func findStructs(packages []*stdpackages.Package) (result []*Struct, err error) {
	err = findTypes(
		packages,
		func(tp stdtypes.Type) (ok bool) {
			_, ok = tp.(*stdtypes.Struct)
			return
		},
		func(namedType *types.Named) (err error) {
			st := namedType.Struct()
			result = append(result, &Struct{
				Named: namedType,
				Type:  st,
			})
			return
		},
	)
	if err != nil {
		return
	}
	return
}

func findTypes(packages []*stdpackages.Package, checkTypeFn func(tp stdtypes.Type) bool, callbackFn func(namedType *types.Named) error) (err error) {
	err = TraverseDefs(packages, func(pkg *stdpackages.Package, id *ast.Ident, obj stdtypes.Object) error {
		if id.Obj == nil {
			return nil
		}
		if id.Obj.Kind != ast.Typ {
			return nil
		}
		named, ok := obj.Type().(*stdtypes.Named)
		if !ok {
			return nil
		}

		if named.Obj().IsAlias() {
			return nil
		}

		if !ok {
			return nil
		}
		if !checkTypeFn(named.Underlying()) {
			return nil
		}

		entity, err := types.NewDecoder(pkg, packages).Decode(obj)
		if err != nil {
			return err
		}
		if namedType, ok := entity.(*types.Named); ok {
			if err := callbackFn(namedType); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return
	}
	return
}
