package clock

import (
	"time"

	"github.com/Meat-Hook/sumahavi/internal/core"
)

var _ core.Clock = &Service{}

// Service see doc in the core.Clock.
type Service struct{}

// Now implements core.Clock.
func (s *Service) Now() time.Time { return time.Now() }
