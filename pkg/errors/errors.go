package errors

import (
	"go/token"
)

type Level string

type WarningError struct {
	text string
	pos  token.Position
}

func (e *WarningError) Error() string {
	if !e.pos.IsValid() {
		return e.text
	}
	return e.pos.String() + ": " + e.text
}

type FailedError struct {
	text string
	pos  token.Position
}

func (e *FailedError) Error() string {
	if !e.pos.IsValid() {
		return e.text
	}
	return e.pos.String() + ": " + e.text
}

func Error(text string, position token.Position) error {
	return &FailedError{
		text: text,
		pos:  position,
	}
}

func Warn(text string, position token.Position) error {
	return &WarningError{
		text: text,
		pos:  position,
	}
}

func IsFailed(e error) bool {
	_, ok := e.(*FailedError)
	return ok
}

func IsWarning(e error) bool {
	_, ok := e.(*WarningError)
	return ok
}
