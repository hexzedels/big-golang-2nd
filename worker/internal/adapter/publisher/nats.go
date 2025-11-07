package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

type NATSCompletionPublisher struct {
	js  jetstream.JetStream
	log *zap.Logger
}

func NewNATSCompletionPublisher(ctx context.Context, log *zap.Logger, natsURL string) (*NATSCompletionPublisher, error) {
	// Connect to NATS server
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	newJS, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("failed to create jetstream: %w", err)
	}

	log.Info("Connected to NATS JetStream for completion publisher", zap.String("url", natsURL))

	return &NATSCompletionPublisher{
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

func (p *NATSCompletionPublisher) PublishCompletion(ctx context.Context, completion JobCompletion) error {
	data, err := json.Marshal(completion)
	if err != nil {
		return fmt.Errorf("failed to marshal completion: %w", err)
	}

	subject := "JOBS.completed"
	_, err = p.js.Publish(ctx, subject, data)
	if err != nil {
		p.log.Error("Failed to publish completion",
			zap.String("job_id", completion.JobID),
			zap.String("subject", subject),
			zap.Error(err))
		return fmt.Errorf("failed to publish completion to NATS: %w", err)
	}

	p.log.Info("Published job completion to NATS",
		zap.String("job_id", completion.JobID),
		zap.String("status", completion.Status),
		zap.String("subject", subject))

	return nil
}
