package source

import (
	"encoding/json"
	"sync"

	"github.com/rs/zerolog"

	"github.com/Meat-Hook/sumahavi/internal/core"
)

var _ core.Container = &container{}

// See doc in the core.Container.
type container struct {
	logger zerolog.Logger
	once   sync.Once

	logs   chan json.RawMessage // Non-blocking on send, closes by Close.
	errors chan error           // Non-blocking on send, closes by Close.
}

// Logs implements core.Container.
func (c *container) Logs() <-chan json.RawMessage { return c.logs }

// Err implements core.Container.
func (c *container) Err() <-chan error { return c.errors }

// Close implements core.Container.
func (c *container) Close() { c.once.Do(func() { c.close() }) }

func (c *container) close() {
	close(c.logs)
	close(c.errors)
}

func (c *container) process() {
	defer c.Close()

}
