package types

import (
	"container/list"
	"strings"

	"github.com/fatih/structtag"
)

type Struct struct {
	Fields    []*StructFieldType
	Graph     map[string]*StructFieldType
	IsPointer bool
}

func (s *Struct) Path(v string) (vv *StructFieldType) {
	graph := s.Graph
	l := split(v, ".")
	for e := l.Front(); e != nil; e = e.Next() {
		if graph == nil {
			break
		}
		part := e.Value.(string)
		start := strings.Index(part, "[")
		if start != -1 {
			part = part[:start]
		}
		var ok bool
		vv, ok = graph[part]
		if ok {
			graph = graphByType(vv.Var.Type)
		}
	}
	return
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

func split(s, sep string) (l *list.List) {
	l = list.New()
	i := 0
	for {
		m := strings.Index(s, sep)
		if m < 0 {
			if len(s) != 0 {
				l.PushBack(s)
			}
			break
		}
		l.PushBack(s[:m])
		s = s[m+len(sep):]
		i++
	}
	return
}
