package types

import (
	"go/token"
)

type Vars []*Var

func (v Vars) HasError() bool {
	return v.Error() != nil
}

func (v Vars) Error() *Var {
	for _, vv := range v {
		if vv.IsError {
			return vv
		}
	}
	return nil
}

func (v Vars) Len() int {
	return len(v)
}

func (v Vars) LenFunc(f func(v *Var) bool) int {
	var offset int
	for _, val := range v {
		if f(val) {
			offset++
		}
	}
	return v.Len() - offset
}

type Var struct {
	Name       string
	Embedded   bool
	Exported   bool
	IsField    bool
	IsVariadic bool
	IsContext  bool
	IsError    bool
	IsChan     bool
	IsPointer  bool
	IsString   bool
	Type       any
	Title      string
	Zero       string
	Tags       Tags
	Position   token.Position
}
