package store

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/Meat-Hook/sumahavi/internal/core"
)

var _ core.Store = &Service{}

// Service see doc in the core.Store.
type Service struct{}

// Search implements core.Store.
func (s *Service) Search(ctx context.Context, terms []string, limit, offset int) ([]core.Record, error) {
	panic("implement me")
}

// Save implements core.Store.
func (s *Service) Save(ctx context.Context, tokens []core.Token, record core.Record) error {
	panic("implement me")
}

// Get implements core.Store.
func (s *Service) Get(ctx context.Context, u uuid.UUID) (*core.Record, error) {
	panic("implement me")
}
