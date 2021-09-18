package source

import (
	"github.com/Meat-Hook/sumahavi/internal/core"
)

var _ core.Source = &Service{}

// Service see doc in the core.Source.
type Service struct{}

// New implements core.Source.
func (s *Service) New() <-chan core.Container { panic("implement me") }

// Close implements core.Source.
func (s *Service) Close() { panic("implement me") }

// Err implements core.Source.
func (s *Service) Err() chan error { panic("implement me") }
