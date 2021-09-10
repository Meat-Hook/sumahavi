package core

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// Start runs system.
// Will start survey and parsing log data from directory.
func (c *Core) Start(ctx context.Context) error {
	var wg sync.WaitGroup
	errc := make(chan error)

	for {
		select {
		// Get new file for parsing.
		case newPath := <-c.disk.NewFile():
			wg.Add(1)
			go func() {
				defer wg.Done()

				err := c.parse(ctx, newPath)
				if err != nil {
					errc <- err
				}
			}()

		// Listening some errors.
		// If we get error, we will close all out channels and finished this function.
		case err := <-errc:
			return err
		case err := <-c.disk.Err():
			return err
		}
	}
}

func (c *Core) parse(ctx context.Context, path string) error {
	errc := make(chan error)
	source := make(chan json.RawMessage)
	defer close(errc)
	defer close(source)

	go func() {
		err := c.parser.Parse(ctx, path, source)
		if err != nil {
			errc <- fmt.Errorf("c.parser.Parse: %w", err)
		}
	}()

	for msg := range source {
		tokens, err := c.tokenizer.Tokens(msg)
		if err != nil {
			return fmt.Errorf("c.tokenizer.Tokens: %w", err)
		}

		id := c.uuid.New()

		err = c.index.Add(ctx, tokens, id)
		if err != nil {
			return fmt.Errorf("c.index.Add: %w", err)
		}

		err = c.store.Save(ctx, id, msg)
		if err != nil {
			return fmt.Errorf("c.store.Save: %w", err)
		}
	}

	return nil
}
