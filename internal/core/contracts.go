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

// Container contains log file.
type Container interface {
	// Name returns container name.
	// It can be file name or for example link to source on another service.
	Name() string
	// Logs returns channel for collecting log lines.
	//
	// This method must be called. Repeated calls returns the same channel.
	// Returned channel is non-blocking on send, closes by Close.
	//
	// Errors: Any, io.EOF, ErrNotFound.
	Logs() <-chan json.RawMessage
	// Err returns channel for listening errors.
	//
	// This method must be called. Repeated calls returns the same channel.
	// Returned channel is non-blocking on send, closes by Close.
	//
	// If you get error from this channel, all channels will close.
	// Only one msg.
	//
	// Errors: Any, io.EOF.
	Err() <-chan error
	// Close container. You shouldn't use Container after close.
	// Also, it will close all channels returned by other methods.
	Close()
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
	Search(ctx context.Context, terms []Token, limit, offset int) ([]Record, error)
}

// Source responsible for getting log journal.
type Source interface {
	// Close source. You shouldn't use Source after close.
	// Also, it will close all channels returned by other methods.
	Close()
	// New returns channel for listening event about new Container.
	//
	// This method must be called. Repeated calls returns the same channel.
	// Returned channel is non-blocking on send, closes by Close.
	New() <-chan Container
	// Err returns channel for listening errors from source space.
	//
	// This method must be called. Repeated calls returns the same channel.
	// Returned channel is non-blocking on send, closes by Close.
	//
	// If you get error from this channel, all channels from Disk will close.
	// Only one msg.
	//
	// Errors: Any.
	Err() chan error
}
