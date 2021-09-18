package container

import (
	"encoding/json"

	"github.com/Meat-Hook/sumahavi/internal/core"
)

var _ core.Container = &Service{}

// Service see doc in the core.Container.
type Service struct{}

// Logs implements core.Container.
func (s *Service) Logs() <-chan json.RawMessage { panic("implement me") }

// Logs implements core.Err.
func (s *Service) Err() <-chan error { panic("implement me") }

// Logs implements core.Close.
func (s *Service) Close() { panic("implement me") }
