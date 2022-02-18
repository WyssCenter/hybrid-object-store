package errors

import "errors"

var (
	ErrExists       = errors.New("record already exists")
	ErrNotFound     = errors.New("record not found")
	ErrNotPermitted = errors.New("user is not permitted to take that requested action")
	ErrInvalidInput = errors.New("invalid input value")
	ErrUnauthorized = errors.New("user is not authorized")
)
