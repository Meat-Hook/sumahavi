package core

import (
	"context"
	"fmt"
)

// Parse runs source parsing.
// Will start survey and parsing log data from specific source.
func (c *Core) Parse(ctx context.Context, source Source) error {
	defer source.Close()

	for msg := range source.Logs() {
		record := Record{
			ID:        c.uuid.New(),
			Name:      source.Name(),
			Body:      msg,
			CreatedAt: c.clock.Now(),
		}

		tokens, err := c.tokenizer.Tokens(msg)
		if err != nil {
			return fmt.Errorf("c.tokenizer.Tokens: %w", err)
		}

		err = c.store.Save(ctx, tokens, record)
		if err != nil {
			return fmt.Errorf("c.store.Save: %w", err)
		}
	}

	return source.Err()
}
