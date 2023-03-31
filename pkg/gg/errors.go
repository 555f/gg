package gg

import "errors"

var (
	ErrSelfReferential      = errors.New("self-referential")
	ErrCircularDependencies = errors.New("circular dependencies")
)
