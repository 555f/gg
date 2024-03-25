package types

import (
	"strings"

	"github.com/fatih/structtag"
)

type Struct struct {
	Fields    []*StructFieldType
	Graph     map[string]*StructFieldType
	IsPointer bool
}

func (s *Struct) Path(v string) *StructFieldType {
	parts := strings.Split(v, ".")
	graph := s.Graph
	for _, part := range parts {
		start := strings.Index(part, "[")
		if start != -1 {
			part = part[:start]
		}
		if vv, ok := graph[part]; ok {
			graph = graphByType(vv.Var.Type)
			if graph == nil || len(parts) == 1 {
				return vv
			}
		}
	}
	return nil
}

type StructFieldType struct {
	Var     *Var
	SysTags *structtag.Tags
}

func graphByType(i any) map[string]*StructFieldType {
	switch t := i.(type) {
	default:
		return nil
	case *Struct:
		return t.Graph
	case *Named:
		return graphByType(t.Type)
	case *Slice:
		return graphByType(t.Value)
	case *Array:
		return graphByType(t.Value)
	}
}
