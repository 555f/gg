package types

import "github.com/fatih/structtag"

type Struct struct {
	Fields    []*StructFieldType
	IsPointer bool
}

type StructFieldType struct {
	Var     *Var
	SysTags *structtag.Tags
}
