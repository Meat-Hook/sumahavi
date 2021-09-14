package disk

import (
	"github.com/Meat-Hook/sumahavi/internal/core"
)

var _ core.Disk = &Service{}

// Service see doc in the core.Disk.
type Service struct{}

// Close implements core.Disk.
func (s *Service) Close() {
	panic("implement me")
}

// NewFile implements core.Disk.
func (s *Service) NewFile() <-chan string {
	panic("implement me")
}

// Err implements core.Disk.
func (s *Service) Err() chan error {
	panic("implement me")
}
