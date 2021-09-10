package core

import "errors"

// Errors.
var (
	ErrNotFound     = errors.New("not found")
	ErrNotFreeSpace = errors.New("disk doesn't have a free space")
)
