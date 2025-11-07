package subscriber

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

type NATSCompletionSubscriber struct {
	js  jetstream.JetStream
	log *zap.Logger
}

func NewNATSCompletionSubscriber(ctx context.Context, log *zap.Logger, natsURL string) (*NATSCompletionSubscriber, error) {
	// Connect to NATS server
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	newJS, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("failed to create jetstream: %w", err)
	}

	log.Info("Connected to NATS JetStream for completion subscriber", zap.String("url", natsURL))

	return &NATSCompletionSubscriber{
		js:  newJS,
		log: log,
	}, nil
}

type JobCompletion struct {
	JobID        string  `json:"jobId"`
	Status       string  `json:"status"` // "completed" or "failed"
	FinishedAt   int64   `json:"finishedAt"`
	ErrorMessage *string `json:"errorMessage,omitempty"`
}

// Subscribe subscribes to job completion messages and calls handler for each completion
func (s *NATSCompletionSubscriber) Subscribe(ctx context.Context, handler func(ctx context.Context, completion JobCompletion) error) error {
	// Create consumer for completion messages
	cons, err := s.js.CreateOrUpdateConsumer(ctx, "JOBS", jetstream.ConsumerConfig{
		FilterSubject: "JOBS.completed",
		Durable:       "scheduler-completion",
	})
	if err != nil {
		return fmt.Errorf("failed to create completion consumer: %w", err)
	}

	// Process completion messages
	go s.processMessages(ctx, cons, handler)

	return nil
}

func (s *NATSCompletionSubscriber) processMessages(ctx context.Context, cons jetstream.Consumer, handler func(ctx context.Context, completion JobCompletion) error) {
	iter, err := cons.Messages()
	if err != nil {
		s.log.Error("Failed to create message iterator", zap.Error(err))
		return
	}
	defer iter.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := iter.Next()
			if err != nil {
				s.log.Error("Failed to get next message", zap.Error(err))
				continue
			}

			var completion JobCompletion
			if err := json.Unmarshal(msg.Data(), &completion); err != nil {
				s.log.Error("Failed to unmarshal completion", zap.Error(err))
				msg.Nak()
				continue
			}

			if err := handler(ctx, completion); err != nil {
				s.log.Error("Handler failed", zap.Error(err))
				msg.Nak()
				continue
			}

			msg.Ack()
		}
	}
}
