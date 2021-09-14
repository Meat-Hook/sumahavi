package tokenizer

import (
	"encoding/json"

	"github.com/Meat-Hook/sumahavi/internal/core"
)

var _ core.Tokenizer = &Service{}

// Service see doc in the core.Tokenizer.
type Service struct{}

// Tokens implements core.Tokenizer.
func (s *Service) Tokens(raw json.RawMessage) ([]core.Token, error) { panic("implement me") }
