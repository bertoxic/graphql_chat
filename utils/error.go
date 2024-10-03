package utils

import "errors"

var (
	ErrValidation     = errors.New("validation error")
	ErrNotFound       = errors.New("not found ")
	ErrInternalServer = errors.New("Internal server error")
	ErrInvalidInput   = errors.New("invalid input")
	ErrUserExists     = errors.New("user already exists")
)
