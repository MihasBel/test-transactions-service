package deliver

import (
	"context"
	"encoding/json"
	"time"

	"github.com/MihasBel/test-transactions-servise/adapters/broker"
	"github.com/MihasBel/test-transactions-servise/internal/app"
	"github.com/MihasBel/test-transactions-servise/internal/rep"
	model "github.com/MihasBel/test-transactions-servise/models"
	"github.com/rs/zerolog/log"
)

const (
	tranTopic = "transaction"
)

// GRPC represent GRPC service
type GRPC struct {
	app string
	cfg app.Configuration
	s   rep.Storage
	b   *broker.Broker
}

// New Create new instance of GRPC. Should use only in main.
func New(config app.Configuration, s rep.Storage, b *broker.Broker) *GRPC {

	grpc := GRPC{
		app: "",
		cfg: config,
		s:   s,
		b:   b,
	}

	return &grpc
}

// Start an application
func (r *GRPC) Start(_ context.Context) error {
	errCh := make(chan error)
	log.Debug().Msgf("start listening %q", r.cfg.Address)
	if err := r.b.Subscribe(context.Background(), tranTopic); err != nil {
		errCh <- err
	}

	go func() {
		for val := range r.b.Ch {
			tran := model.Transaction{}
			if err := json.Unmarshal(val, &tran); err != nil {
				log.Error().Err(err)
			}
			if err := r.s.PlaceTransaction(context.Background(), tran); err != nil {
				log.Error().Err(err)
			}
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(time.Duration(r.cfg.StartTimeout) * time.Second):
		return nil
	}
}

// Stop an application
func (r *GRPC) Stop(_ context.Context) error {
	errCh := make(chan error)
	log.Debug().Msgf("stopping %q", r.cfg.Address)
	go func() {
		close(r.b.Ch)
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(time.Duration(r.cfg.StopTimeout) * time.Second):
		return nil

	}
}
