package parser

import (
	"context"
	"encoding/json"

	"github.com/Meat-Hook/sumahavi/internal/core"
)

var _ core.Parser = &Service{}

// Service see doc in the core.Parser.
type Service struct{}

// Parse implements core.Parser.
func (s *Service) Parse(ctx context.Context, path string, msg chan<- json.RawMessage) error {
	panic("implement me")
}
