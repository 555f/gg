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

type VarParser struct{}

func (*VarParser) Parse(s string) (vars []Var) {
	val := strings.TrimSpace(s)

	startPos, finishPos := -1, -1

	var i int
	for {

		switch val[i] {
		default:
			if startPos != -1 && finishPos != -1 {
				vars = append(vars, Var{ID: val[startPos+1 : finishPos], Pos: VarPos{Start: startPos, Finish: finishPos}})
				startPos, finishPos = -1, -1
			}
		case '{':
			startPos = i
		case '}':
			finishPos = i
		}

		i++

		if i >= len(val) {
			break
		}

	}

	// for val != "" {
	// 	i := 0
	// 	for i < len(val) && val[i] != '{' {
	// 		i++
	// 	}

	// 	val = val[i:]
	// 	if val == "" {
	// 		break
	// 	}

	// 	begin := i

	// 	val = val[1:]
	// 	i = 0

	// 	for i < len(val) && val[i] != '}' {
	// 		i++
	// 	}

	// 	end := begin + i + 2

	// 	id := val[:i]
	// 	val = val[i+1:]

	// 	vars = append(vars, Var{ID: id, Pos: VarPos{Begin: begin, End: end}})
	// }
	return
}

func NewVarParser() *VarParser {
	return &VarParser{}
}
