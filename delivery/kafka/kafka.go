package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/MihasBel/test-transactions-servise/adapters/broker"
	"github.com/MihasBel/test-transactions-servise/internal/rep"
	model "github.com/MihasBel/test-transactions-servise/models"
	"github.com/rs/zerolog/log"
)

const (
	tranTopic = "transaction"
)

// Server represent Server service
type Server struct {
	app string
	cfg Config
	s   rep.Storage
	b   *broker.Broker
}

// New Create new instance of Server. Should use only in main.
func New(config Config, s rep.Storage, b *broker.Broker) *Server {

	grpc := Server{
		app: "",
		cfg: config,
		s:   s,
		b:   b,
	}

	return &grpc
}

// Start an application
func (r *Server) Start(_ context.Context) error {
	errCh := make(chan error)
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
func (r *Server) Stop(_ context.Context) error {
	errCh := make(chan error)
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
