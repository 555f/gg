package types

type Interface struct {
	Methods         []*Func
	Embedded        []any
	ExplicitMethods []*Func
}

// func (t *Interface) Implement(n *Named) bool {
// 	return stdtypes.Implements(stdtypes.NewPointer(n.origin), t.origin)
// }
