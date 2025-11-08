package core

import "errors"

var (
	ErrNotValidRequest = errors.New("not a valid Request")
	ErrNotFound        = errors.New("not found")
)
