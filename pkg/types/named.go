package types

import (
	"go/token"
)

type Named struct {
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

func (n *Named) IsSystemType() bool {
	key := n.Pkg.Path + "." + n.Name
	switch key {
	case "time.Time", "time.Duration":
		return true
	}
	return false
}

func (n *Named) Interface() *Interface {
	i, _ := n.Type.(*Interface)
	return i
}

func (n *Named) Struct() *Struct {
	i, _ := n.Type.(*Struct)
	return i
}

func (n *Named) Basic() *Basic {
	i, _ := n.Type.(*Basic)
	return i
}

var NamedTyp = map[string]*Named{
	"time.Time": {
		Name: "Time",
		Type: &Struct{},
		Pkg: &PackageType{
			Name: "Time",
			Path: "time",
		},
	},
}
