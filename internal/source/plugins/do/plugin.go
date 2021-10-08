package do

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog"
	"go.uber.org/ratelimit"

	"github.com/Meat-Hook/sumahavi/internal/core"
	"github.com/Meat-Hook/sumahavi/internal/source"
)

var _ core.Source = &Plugin{}

// Default values.
const (
	DefaultStepLimit      = 100
	DefaultLimitPerSecond = 1
)

type Config struct {
	Name           string
	AppID          string
	DeploymentID   string
	StepLimit      int
	LimitPerSecond int
}

// Plugin see docs source.Plugin.
type Plugin struct {
	cfg     Config
	client  *godo.Client
	logs    chan json.RawMessage // Non-blocking on send, closes by Close.
	once    sync.Once
	limiter ratelimit.Limiter
	error   error
}

// Name implements core.Source.
func (p *Plugin) Name() string { return p.cfg.Name }

// Err implements core.Source.
func (p *Plugin) Err() error { return p.error }

// Logs implements core.Source.
func (p *Plugin) Logs() <-chan json.RawMessage { return p.logs }

// Close implements core.Source.
func (p *Plugin) Close() { p.once.Do(func() { p.close() }) }

func (p *Plugin) close() {
	close(p.logs)
	if p.error == nil {
		p.error = source.ErrHasClosed
	}
}

// New build and returns new instance Plugin.
func New(ctx context.Context, token string, cfg Config) (*Plugin, error) {
	// Validation config.
	switch {
	case cfg.AppID == "":
		return nil, fmt.Errorf("%w: didn't set AppID", source.ErrInvalidConfig)
	case cfg.DeploymentID == "":
		return nil, fmt.Errorf("%w: didn't set DeploymentID", source.ErrInvalidConfig)
	case cfg.Name == "":
		return nil, fmt.Errorf("%w: didn't set Name", source.ErrInvalidConfig)
	}

	// Set default values.
	if cfg.LimitPerSecond == 0 {
		cfg.LimitPerSecond = DefaultLimitPerSecond
	}

	if cfg.StepLimit == 0 {
		cfg.StepLimit = DefaultStepLimit
	}

	p := &Plugin{
		cfg:     cfg,
		client:  godo.NewFromToken(token),
		logs:    make(chan json.RawMessage),
		once:    sync.Once{},
		limiter: ratelimit.New(cfg.LimitPerSecond),
		error:   nil,
	}

	go func() {
		logger := zerolog.Ctx(ctx)
		logger.Debug().Msg("do plugin has started parse")
		defer logger.Debug().Msg("do plugin done")
		err := p.process(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("error from DigitalOcean")
		}
	}()

	return p, nil
}

func (p *Plugin) process(ctx context.Context) error {
	defer p.Close()
	logger := zerolog.Ctx(ctx)

	errors := make(chan error)
	lines := make(chan string)

	go func() {
		err := p.parse(ctx, lines)
		if err != nil {
			errors <- err
		}
	}()

	lastReadEvent := time.Time{}
	for {
		select {
		case <-ctx.Done():
			return nil

		case msg := <-errors:
			p.error = msg
			return p.error

		case msg := <-lines:
			logger.Debug().Str("line", msg).Msg("read line")

			array := strings.SplitN(msg, " ", 3)
			name := array[0]
			doLogTime, err := time.Parse(time.RFC3339Nano, array[1])
			if err != nil {
				logger.Error().Str("line", msg).Err(err).Msg("parse time")
				continue
			}

			if doLogTime.Before(lastReadEvent) || doLogTime.Equal(lastReadEvent) {
				logger.Debug().
					Str("line", msg).
					Time("last_read", lastReadEvent).
					Time("do_time", doLogTime).
					Msg("ignore")
				continue
			}

			lastReadEvent = doLogTime

			body := json.RawMessage(array[2])
			if !json.Valid(body) {
				// TODO: add json convert.
				logger.Warn().Str("body", array[2]).Msg("ignore")
				continue
			}

			l := line{
				AppName: name,
				DOTime:  doLogTime,
				Body:    body,
			}

			buf, err := json.Marshal(l)
			if err != nil {
				return fmt.Errorf("json.Marshal: %w", err)
			}

			p.logs <- buf
		}
	}
}

type line struct {
	AppName string          `json:"app_name"`
	DOTime  time.Time       `json:"do_time"`
	Body    json.RawMessage `json:"body"`
}

func (p *Plugin) parse(ctx context.Context, result chan<- string) error {
	for {
		p.limiter.Take()

		urls, _, err := p.client.Apps.GetLogs(ctx, p.cfg.AppID, p.cfg.DeploymentID, p.cfg.Name, godo.AppLogTypeRun, false, p.cfg.StepLimit)
		if err != nil {
			return fmt.Errorf("p.client.Apps.GetLogs: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, urls.LiveURL, nil)
		if err != nil {
			return fmt.Errorf("http.NewRequestWithContext: %w", err)
		}

		buffer := &bytes.Buffer{}
		resp, err := p.client.Do(ctx, req, buffer)
		if err != nil {
			return fmt.Errorf("p.client.Do: %w", err)
		}

		err = resp.Body.Close()
		if err != nil {
			return fmt.Errorf("resp.Body.Close: %w", err)
		}

		scanner := bufio.NewScanner(buffer)

		for scanner.Scan() {
			result <- scanner.Text()
		}

		err = scanner.Err()
		if err != nil {
			return fmt.Errorf("scanner.Err: %w", err)
		}
	}
}
