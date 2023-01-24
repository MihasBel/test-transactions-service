package broker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/MihasBel/test-transactions-service/internal/rep"
	model "github.com/MihasBel/test-transactions-service/models"
	"github.com/rs/zerolog/log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog"
)

// Broker represents kafka broker with logger and config
type Broker struct {
	b   *kafka.Consumer
	cfg Config
	l   zerolog.Logger
	s   rep.Storage
	Ch  chan []byte
}

// New create new Broker
func New(cfg Config, l zerolog.Logger, s rep.Storage, ch chan []byte) *Broker {

	return &Broker{
		Ch:  ch,
		cfg: cfg,
		l:   l,
		s:   s,
	}
}

// Subscribe to get msg from kafka topic
func (brk *Broker) Subscribe(ctx context.Context, topic string) (err error) {
	brk.b, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        brk.cfg.KafkaURL,
		"group.id":                 "transaction",
		"auto.offset.reset":        "smallest",
		"allow.auto.create.topics": true,
	})

	if err != nil {
		return err
	}

	if err = brk.b.Subscribe(topic, nil); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			msg, err := brk.b.ReadMessage(-1)
			if err == nil {
				brk.l.Info().Msgf("[%v]: %s", msg.String(), msg.Value)
				brk.Ch <- msg.Value

			} else {
				brk.l.
					Error().
					Err(err).
					Interface("msg", msg).
					Msg("consumer error")

			}

		}
	}()

	return nil
}

// Start a broker
func (brk *Broker) Start(_ context.Context) error {
	errCh := make(chan error)
	if err := brk.Subscribe(context.Background(), brk.cfg.Topic); err != nil {
		errCh <- err
	}

	go func() {
		for val := range brk.Ch {
			tran := model.Transaction{}
			if err := json.Unmarshal(val, &tran); err != nil {
				log.Error().Err(err)
			}
			if err := brk.s.PlaceTransaction(context.Background(), tran); err != nil {
				log.Error().Err(err)
			}
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(time.Duration(brk.cfg.StartTimeout) * time.Second):
		return nil
	}
}

// Stop a broker
func (brk *Broker) Stop(_ context.Context) error {
	errCh := make(chan error)
	go func() {
		if err := brk.b.Close(); err != nil {
			errCh <- err
		}
		close(brk.Ch)
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(time.Duration(brk.cfg.StopTimeout) * time.Second):
		return nil

	}
}
