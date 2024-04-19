package gg

import (
	"go/ast"
	stdtypes "go/types"

	"github.com/555f/gg/pkg/types"

	"golang.org/x/tools/go/packages"
)

func Module(pkgs []*packages.Package) (*types.Module, error) {
	for _, p := range pkgs {
		if p.Module != nil && p.Module.Main {
			result, err := types.NewDecoder(p, nil).Decode(p.Module)
			if err != nil {
				return nil, err
			}
			m, _ := result.(*types.Module)
			return m, nil
		}
	}
	return nil, nil
}

func TraverseDefs(pkgs []*packages.Package, c func(pkg *packages.Package, id *ast.Ident, obj stdtypes.Object) error) error {
	for _, pkg := range pkgs {
		for id, obj := range pkg.TypesInfo.Defs {
			if obj == nil {
				continue
			}
			if err := c(pkg, id, obj); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyNodeset(s NodeSet) NodeSet {
	out := make(NodeSet, len(s))
	for k, v := range s {
		out[k] = v
	}
	return out
}

func copyDependency(m Dependency) Dependency {
	out := make(Dependency, len(m))
	for k, v := range m {
		out[k] = copyNodeset(v)
	}
	return out
}

func removeFromDependency(dm Dependency, key, node string) {
	nodes := dm[key]
	if len(nodes) == 1 {
		delete(dm, key)
	} else {
		delete(nodes, node)
	}
}
