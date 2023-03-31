package types

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

func traverseDecls(packages []*packages.Package, c func(pkg *packages.Package, file *ast.File, decl ast.Decl) error) error {
	for _, pkg := range packages {
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				if err := c(pkg, file, decl); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func zeroValue(t types.Type) string {
	switch u := t.Underlying().(type) {
	case *types.Basic:
		if u.Kind() == types.UnsafePointer {
			return "nil"
		}
		info := u.Info()
		switch {
		case info&types.IsBoolean != 0:
			return "false"
		case info&(types.IsFloat|types.IsComplex) != 0:
			return "0.0"
		case info&(types.IsInteger|types.IsUnsigned|types.IsUntyped) != 0:
			return "0"
		case info&types.IsString != 0:
			return `""`
		}
	case *types.Struct:
		return "{}"
	case *types.Chan, *types.Interface, *types.Map, *types.Pointer, *types.Signature, *types.Slice, *types.Array:
		return "nil"
	}
	panic("unreachable")
}

func IsError(v any) bool {
	if named, ok := v.(*Named); ok {
		if _, ok := named.Type.(*Interface); ok && named.Name == "error" {
			return true
		}
	}
	return false
}

func IsContext(t any) bool {
	if named, ok := t.(*Named); ok {
		if _, ok := named.Type.(*Interface); ok {
			if named.Name == "Context" && named.Pkg.Path == "context" {
				return true
			}
		}
	}
	return false
}
