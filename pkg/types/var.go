package types

import (
	"go/token"
)

type Vars []*Var

func (v Vars) LenFunc(f func(v *Var) bool) int {
	l := len(v)
	var offset int
	for _, val := range v {
		if f(val) {
			offset++
		}
	}
	return l - offset
}

type Var struct {
	Name       string
	Embedded   bool
	Exported   bool
	IsField    bool
	IsVariadic bool
	IsContext  bool
	IsError    bool
	Type       any
	Title      string
	Zero       string
	Tags       Tags
	Position   token.Position
}
