package jsonpointer

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrUnexpectedType = errors.New("unexpected type")
	ErrOutOfBounds    = errors.New("out of bounds")
)
