package broker

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog"
)

// Broker represents kafka broker with logger and config
type Broker struct {
	b   *kafka.Consumer
	cfg Config
	l   zerolog.Logger
	Ch  chan []byte
}

// New create new Broker
func New(cfg Config, l zerolog.Logger, ch chan []byte) *Broker {

	return &Broker{
		Ch:  ch,
		cfg: cfg,
		l:   l,
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

// Stop broker
func (brk *Broker) Stop(_ context.Context) (err error) {
	if err = brk.b.Close(); err != nil {
		return err
	}
	return nil
}
