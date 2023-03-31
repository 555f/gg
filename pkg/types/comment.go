package types

import "go/token"

type Comment struct {
	Value    string
	IsTitle  bool
	IsTag    bool
	Position token.Position
}

type Comments []*Comment

func (c Comments) String() (result string) {
	for _, comment := range c {
		result += comment.Value + "\n"
	}
	return
}
