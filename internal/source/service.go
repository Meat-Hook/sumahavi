package source

import (
	"sync"

	"github.com/Meat-Hook/sumahavi/internal/core"
)

var _ core.Source = &Service{}

// Service see doc in the core.Source.
type Service struct {
	name string
	once sync.Once

	containers chan core.Container // Non-blocking on send, closes by Close.
	errors     chan error          // Non-blocking on send, closes by Close.
}

// New build and returns new Service.
func New(name string) *Service {
	s := &Service{
		name:       name,
		once:       sync.Once{},
		containers: make(chan core.Container),
		errors:     make(chan error),
	}

	return s
}

// Name implements core.Source.
func (s *Service) Name() string { return s.name }

// New implements core.Source.
func (s *Service) New() <-chan core.Container { return s.containers }

// Err implements core.Source.
func (s *Service) Err() chan error { return s.errors }

// Close implements core.Source.
func (s *Service) Close() { s.once.Do(func() { s.close() }) }

func (s *Service) close() {
	close(s.containers)
	close(s.errors)
}
