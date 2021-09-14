package index

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/Meat-Hook/sumahavi/internal/core"
)

var _ core.InvertedIndex = &Service{}

// Service see doc in the core.InvertedIndex.
type Service struct{}

// Add implements core.InvertedIndex.
func (s *Service) Add(ctx context.Context, terms []core.Token, id uuid.UUID) error {
	panic("implement me")
}

// Search implements core.InvertedIndex.
func (s *Service) Search(ctx context.Context, terms []core.Token, limit, offset int) ([]uuid.UUID, error) {
	panic("implement me")
}
