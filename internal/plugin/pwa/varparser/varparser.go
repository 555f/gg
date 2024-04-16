package varparser

import (
	"strings"
)

type VarPos struct {
	Start, Finish int
}

type Var struct {
	ID  string
	Pos VarPos
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
			vars = append(vars, Var{ID: val[startPos+1 : finishPos], Pos: VarPos{Start: startPos, Finish: finishPos + 1}})
			startPos, finishPos = -1, -1
		}
		i++
		if i >= len(val) {
			break
		}
	}
	return
}
