package gg

import (
	"go/ast"
	stdtypes "go/types"

	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"

	"golang.org/x/tools/go/packages"
)

func ZeroValue(t any, qualFunc types.QualFunc) jen.Code {
	switch t := t.(type) {
	default:
		return jen.Nil()
	case *types.Chan, *types.Interface, *types.Map, *types.Slice, *types.Array:
		return jen.Nil()
	case *types.Basic:
		switch {
		default:
			return jen.Lit(0)
		case t.IsBool():
			return jen.Lit(false)
		case t.IsString():
			return jen.Lit("")
		case t.IsInteger():
			return jen.Lit(0)
		case t.IsFloat():
			return jen.Lit(0.0)
		}
	case *types.Named:
		if t.IsPointer {
			return jen.Nil()
		}
		if t.Name == "error" {
			return jen.Nil()
		}
		return jen.Do(qualFunc(t.Pkg.Path, t.Name)).Values()
	}
}

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
