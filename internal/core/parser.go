package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
)

// Start runs parser.
// Will start survey and parsing log data from directory.
func (c *Core) Start(ctx context.Context) error {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		c.disk.Close()
	}()

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
					errc <- fmt.Errorf("c.parse: %w", err)
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

	for {
		select {
		// New log line.
		case msg := <-source:
			record := Record{
				ID:        c.uuid.New(),
				Body:      msg,
				CreatedAt: c.clock.Now(),
			}

			tokens, err := c.tokenizer.Tokens(msg)
			if err != nil {
				return fmt.Errorf("c.tokenizer.Tokens: %w", err)
			}

			err = c.store.Save(ctx, record)
			if err != nil {
				return fmt.Errorf("c.store.Save: %w", err)
			}

			err = c.index.Add(ctx, tokens, record.ID)
			if err != nil {
				return fmt.Errorf("c.index.Add: %w", err)
			}

		// Error from parsing process.
		case err := <-errc:
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				return nil
			case errors.Is(err, context.Canceled):
				return nil
			case errors.Is(err, io.EOF):
				return nil
			default:
				return err
			}
		}
	}
}
