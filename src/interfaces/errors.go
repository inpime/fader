package interfaces

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrInternal = errors.New("internal error")
)
