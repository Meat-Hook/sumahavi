package store

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/Meat-Hook/sumahavi/internal/core"
)

var _ core.Store = &Service{}

// Service see doc in the core.Store.
type Service struct{}

func (s *Service) Save(ctx context.Context, record core.Record) error {
	panic("implement me")
}

func (s *Service) Get(ctx context.Context, u uuid.UUID) (*core.Record, error) {
	panic("implement me")
}

// New see doc in the core.UUID.
func (s *Service) New() uuid.UUID { panic("implement me") }
