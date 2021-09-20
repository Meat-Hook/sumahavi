package core

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

// Token contains one token for making inverted index by terms.
type Token string

// Record contains all information about log line.
type Record struct {
	ID        uuid.UUID       // Generate in core layer.
	Name      string          // Source name.
	Body      json.RawMessage // Raw body.
	CreatedAt time.Time       // Will set in core layer.
}
