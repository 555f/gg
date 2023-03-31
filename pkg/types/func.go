package types

import "go/token"

type Func struct {
	Pkg         *PackageType
	FullName    string
	Name        string
	Exported    bool
	Sig         *Sign
	Title       string
	Description string
	Tags        Tags
	Position    token.Position
	Returns     []any
	Text        string
}
