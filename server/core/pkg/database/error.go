package database

import (
	"errors"
	"strings"
)

var (
	ErrExists       = errors.New("record already exists")
	ErrNotFound     = errors.New("record not found")
	ErrNotPermitted = errors.New("user is not permitted to take that requested action")
	ErrInvalidInput = errors.New("invalid input value")
	ErrUnauthorized = errors.New("user is not authorized")
)

func ConvertError(err error) error {
	// pg errors
	errMsg := err.Error()
	if strings.Contains(errMsg, "duplicate key value") {
		return ErrExists
	} else if strings.Contains(errMsg, "no rows in result set") {
		return ErrNotFound
	} else if strings.Contains(errMsg, "invalid input value") {
		return ErrInvalidInput
	} else if strings.Contains(errMsg, "already exists") {
		return ErrExists
	}

	return err
}
