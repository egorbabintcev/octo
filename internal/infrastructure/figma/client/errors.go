package figma

import "errors"

var (
	ErrUnknown            = errors.New("unknown error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidRequest     = errors.New("invalid request data")
)
