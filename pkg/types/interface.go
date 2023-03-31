package types

import stdtypes "go/types"

type Interface struct {
	origin          *stdtypes.Interface
	Methods         []*Func
	Embedded        []any
	ExplicitMethods []*Func
}

func (t *Interface) Implement(n *Named) bool {
	return stdtypes.Implements(stdtypes.NewPointer(n.origin), t.origin)
}
