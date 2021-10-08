package source

import "errors"

// Errors.
var (
	ErrHasClosed     = errors.New("source has closed")
	ErrInvalidConfig = errors.New("invalid config")
)
