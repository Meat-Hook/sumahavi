package source

import (
	"encoding/json"

	"github.com/Meat-Hook/sumahavi/internal/core"
)

var _ core.Container = &container{}

// See doc in the core.Container.
type container struct{}

// Logs implements core.Container.
func (s *container) Logs() <-chan json.RawMessage { panic("implement me") }

// Err implements core.Container.
func (s *container) Err() <-chan error { panic("implement me") }

// Close implements core.Container.
func (s *container) Close() { panic("implement me") }
