package varparser

import (
	"strings"
)

type VarPos struct {
	Start, Finish int
}

type Var struct {
	ID   string
	Path string
	Pos  VarPos
}

func Parse(s string) (vars []Var) {
	val := strings.TrimSpace(s)
	startPos, finishPos := -1, -1
	var i int
	for {
		switch val[i] {
		case '{':
			startPos = i
		case '}':
			finishPos = i
		}
		if startPos != -1 && finishPos != -1 {
			var path string
			id := val[startPos+1 : finishPos]

			if opIdx := strings.Index(id, "("); opIdx != -1 {
				if clIdx := strings.Index(id, ")"); clIdx != -1 {
					path = id[opIdx+1 : clIdx]
				}
			}

			vars = append(vars, Var{ID: id, Path: path, Pos: VarPos{Start: startPos, Finish: finishPos + 1}})
			startPos, finishPos = -1, -1
		}
		i++
		if i >= len(val) {
			break
		}
	}
	return
}
