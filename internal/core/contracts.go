package core

import (
	"context"
	"encoding/json"

	"github.com/gofrs/uuid"
)

// UUID for convenient testing.
type UUID interface {
	// New returns new uuid.
	New() uuid.UUID
}

// Tokenizer responsible for making tokens.
type Tokenizer interface {
	// Tokens returns all tokens by json.
	// Errors: Any.
	Tokens(json.RawMessage) ([]Token, error)
}

// Parser responsible for parse log file.
type Parser interface {
	// Parse starts parse file by path.
	// Must read file and returns one log line like a json in channel.
	//
	// Doesn't support concurrency.
	// It isn't required to call Parse or read from channel.
	// For one channel, be sure to call no more than once.
	// If you get error, channel will close.
	//
	// Errors: Any, ErrNotFound.
	Parse(context.Context, string, chan<- json.RawMessage) error
}

// Store responsible for saving log data on disk.
type Store interface {
	// Save new log data on disk.
	// There is concurrency supporting.
	//
	// Errors: Any, ErrNotFound, ErrNotFreeSpace.
	Save(context.Context, uuid.UUID, json.RawMessage) error
	// Get log data by id.
	//
	// Errors: Any, ErrNotFound.
	Get(context.Context, uuid.UUID) (json.RawMessage, error)
}

// InvertedIndex responsible for finding data id by terms.
// Contains data on disk.
type InvertedIndex interface {
	// Add new data id with terms.
	// If we haven't this terms, II will make new terms on disk.
	// If we have this terms, II will add new id for these terms.
	//
	// Errors: Any, ErrNotFreeSpace.
	Add(ctx context.Context, terms []Token, id uuid.UUID) error
	// Search data IDs by terms.
	//
	// Errors: Any, ErrNotFound.
	Search(ctx context.Context, terms []Token, limit, offset int) ([]uuid.UUID, error)
}

// Disk responsible for checking file rotation.
type Disk interface {
	// NewFile returns channel for listening event about new file.
	//
	// This method must be called. Repeated calls returns the same channel.
	// Returned channel is non-blocking on send, closes by Close.
	NewFile() <-chan string
	// Err returns channel for listening errors from disk space.
	//
	// This method must be called. Repeated calls returns the same channel.
	// Returned channel is non-blocking on send, closes by Close.
	//
	// If you get error from this channel, all channels from Disk will close.
	// Only one msg.
	//
	// Errors: Any, ErrNotFound.
	Err() chan error
}
