package types

import (
	"container/list"
	"strings"
)

type Struct struct {
	Fields    []*Var
	Graph     map[string]*Var
	IsPointer bool
}

func (s *Struct) Path(v string) (vv *Var) {
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
			graph = graphByType(vv.Type)
		}
	}
	return
}

func graphByType(i any) map[string]*Var {
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
