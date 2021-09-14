package uuid

import (
	"github.com/gofrs/uuid"

	"github.com/Meat-Hook/sumahavi/internal/core"
)

var _ core.UUID = &Service{}

// Service see doc in the core.UUID.
type Service struct{}

// New implements core.UUID.
func (s *Service) New() uuid.UUID { panic("implement me") }
