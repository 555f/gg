package types

import (
	"go/ast"
	stdtypes "go/types"

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

func ZeroValueJS(t any) string {
	switch u := t.(type) {
	case *Basic:
		if u.Kind == stdtypes.UnsafePointer {
			return "nil"
		}
		switch {
		case u.IsBool():
			return "false"
		case u.IsFloat():
			return "0.0"
		case u.IsInteger():
			return "0"
		case u.IsString():
			return `""`
		}
	case *Named:
		return "{}"
	case *Struct:
		return "{}"
	case *Chan, *Interface, *Map, *Sign, *Slice, *Array:
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

func IsString(v any) bool {
	if basic, ok := v.(*Basic); ok && basic.IsString() {
		return true
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

func IsChan(t any) (ok bool) {
	_, ok = t.(*Chan)
	return
}

func IsPointer(t any) (ok bool) {
	switch t := t.(type) {
	case *Named:
		return t.IsPointer
	case *Basic:
		return t.IsPointer
	case *Var:
		return t.IsPointer
	case *Struct:
		return t.IsPointer
	case *Map:
		return t.IsPointer
	case *Array:
		return t.IsPointer
	case *Slice:
		return t.IsPointer
	}
	return false
}
