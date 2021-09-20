package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
)

// Start runs container.
// Will start survey and parsing log data from directory.
func (c *Core) Start(ctx context.Context, source Source) error {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		source.Close()
	}()

	errc := make(chan error)
	for {
		select {
		// Get new container for staring parse.
		case container := <-source.New():
			wg.Add(1)
			go func() {
				defer wg.Done()

				err := c.parse(ctx, source.Name(), container)
				if err != nil {
					errc <- fmt.Errorf("c.parse: %w", err)
				}
			}()

		// Listening some errors.
		// If we get error, we will close all out channels and finished this function.
		case err := <-errc:
			return err
		case err := <-source.Err():
			return err
		}
	}
}

func (c *Core) parse(ctx context.Context, sourceName string, container Container) error {
	defer container.Close()

	for {
		select {
		// New log line.
		case msg := <-container.Logs():
			record := Record{
				ID:        c.uuid.New(),
				Name:      sourceName,
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

		// Error from parsing process.
		case err := <-container.Err():
			switch {
			case errors.Is(err, io.EOF):
				return nil
			default:
				return err
			}
		}
	}
}
