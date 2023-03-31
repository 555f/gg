package types

import (
	"go/token"
	stdtypes "go/types"
)

type Named struct {
	origin      *stdtypes.Named
	Title       string
	Description string
	Name        string
	Type        any
	Pkg         *PackageType
	IsPointer   bool
	IsExported  bool
	IsAlias     bool
	Methods     []*Func
	Position    token.Position
	Tags        Tags
}

func (n *Named) Interface() *Interface {
	i, _ := n.Type.(*Interface)
	return i
}

func (n *Named) Struct() *Struct {
	i, _ := n.Type.(*Struct)
	return i
}
