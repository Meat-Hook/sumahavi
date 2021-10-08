package core

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

// Clock for convenient testing.
type Clock interface {
	// Now returns actual time.
	Now() time.Time
}

// UUID for convenient testing.
type UUID interface {
	// New returns new uuid.
	New() uuid.UUID
}

// Tokenizer responsible for making tokens.
type Tokenizer interface {
	// Tokens returns all tokens by json.
	//
	// Errors: Any.
	Tokens(json.RawMessage) ([]Token, error)
}

// Store responsible for saving log data on source.
// Contains data on source.
type Store interface {
	// Save new log data on source.
	// There is concurrency supporting.
	//
	// Errors: Any, ErrNotFound, ErrNotFreeSpace.
	Save(context.Context, []Token, Record) error
	// Get log data by id.
	//
	// Errors: Any, ErrNotFound.
	Get(context.Context, uuid.UUID) (*Record, error)
	// Search log data by some terms.
	//
	// Errors: Any, ErrNotFound.
	Search(ctx context.Context, terms []string, limit, offset int) ([]Record, error)
}

// Source responsible for getting log journal.
type Source interface {
	// Name returns source name.
	Name() string
	// Close source. You shouldn't use Source after close.
	// Also, it will close all channels returned by other methods.
	//
	// You can repeat it after first call this method.
	Close()
	// Logs returns channel for collecting log lines.
	//
	// If you don't read this channel parser won't parse data.
	//
	// This method must be called. Repeated calls returns the same channel.
	// Returned channel is non-blocking on send, closes by Close.
	//
	// Errors: Any, io.EOF, ErrNotFound.
	Logs() <-chan json.RawMessage
	// Err returns last getting errors.
	//
	// Can be nil.
	//
	// Errors: Any.
	Err() error
}
