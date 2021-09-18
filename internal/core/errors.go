package core

import "errors"

// Errors.
var (
	ErrNotFound     = errors.New("not found")
	ErrNotFreeSpace = errors.New("source doesn't have a free space")
)
