package webhooks

import "errors"

var (
	ErrInternal = errors.New("internal error")
	ErrUknown   = errors.New("unknown error")
	ErrConflict = errors.New("conflict error")
	ErrNotFound = errors.New("failed to find webhook")
)
