package app

import (
	"context"
	"os"

	"github.com/MihasBel/test-transactions-service/adapters/broker"
	"github.com/MihasBel/test-transactions-service/adapters/pg"
	"github.com/MihasBel/test-transactions-service/delivery/grpc/server"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Lifecycle to start and stop modules
type Lifecycle interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
type cmp struct {
	Service Lifecycle
	Name    string
}

// App represents application
type App struct {
	log  *zerolog.Logger
	cmps []cmp
	cfg  Configuration
}

// New application constructor
func New(cfg Configuration) *App {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("cmp", "app").Logger()
	return &App{
		log:  &l,
		cfg:  cfg,
		cmps: []cmp{},
	}
}

// Start an application
func (a *App) Start(ctx context.Context) error {
	a.log.Info().Msg("starting app")

	db := pg.New(a.cfg.PG, *a.log)
	chanRes := make(chan []byte)
	b := broker.New(a.cfg.Broker, *a.log, db, chanRes)
	srv := server.New(a.cfg.GRPC, db, *a.log)

	a.cmps = append(
		a.cmps,
		cmp{db, "storage"},
		cmp{b, "kafka"},
		cmp{srv, "grpc"},
	)

	okCh, errCh := make(chan struct{}), make(chan error)
	go func() {
		for _, c := range a.cmps {
			a.log.Info().Msgf("%v is starting", c.Name)
			if err := c.Service.Start(ctx); err != nil {
				a.log.Error().Err(err).Msgf("Cannot start %v", c.Name)
				errCh <- errors.Wrapf(err, "Cannot start %v", c.Name)
			}
		}
		okCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return errors.New("Start timeout")
	case err := <-errCh:
		return err
	case <-okCh:
		return nil
	}
}

// Stop an application
func (a *App) Stop(ctx context.Context) error {
	a.log.Info().Msg("shutting down service...")

	okCh, errCh := make(chan struct{}), make(chan error)
	go func() {
		for i := len(a.cmps) - 1; i >= 0; i-- {
			c := a.cmps[i]
			a.log.Info().Msgf("%v is stopping", c.Name)
			if err := c.Service.Start(ctx); err != nil {
				a.log.Error().Err(err).Msgf("Cannot stop %v", c.Name)
				errCh <- errors.Wrapf(err, "Cannot stop %v", c.Name)
			}
		}
		okCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return errors.New("Stop timeout")
	case err := <-errCh:
		return err
	case <-okCh:
		return nil
	}
}
