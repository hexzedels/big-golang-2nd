package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"scheduler/scheduler/pkg/entity"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

type NATSJobSubscriber struct {
	js  jetstream.JetStream
	log *zap.Logger
}

func NewNATSJobSubscriber(ctx context.Context, log *zap.Logger, natsURL string) (*NATSJobSubscriber, error) {
	// Connect to NATS server
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	newJS, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("failed to create jetstream: %w", err)
	}

	log.Info("Connected to NATS JetStream", zap.String("url", natsURL))

	return &NATSJobSubscriber{
		js:  newJS,
		log: log,
	}, nil
}

// Subscribe subscribes to queued job subjects and calls handler for each job
func (s *NATSJobSubscriber) Subscribe(ctx context.Context, handler func(ctx context.Context, job *entity.Job) error) error {
	// Subscribe to interval queued jobs
	intervalCons, err := s.js.CreateOrUpdateConsumer(ctx, "JOBS", jetstream.ConsumerConfig{
		FilterSubject: "JOBS.interval.queued",
		Durable:       "worker-interval",
	})
	if err != nil {
		return fmt.Errorf("failed to create interval consumer: %w", err)
	}

	// Subscribe to once queued jobs
	onceCons, err := s.js.CreateOrUpdateConsumer(ctx, "JOBS", jetstream.ConsumerConfig{
		FilterSubject: "JOBS.once.queued",
		Durable:       "worker-once",
	})
	if err != nil {
		return fmt.Errorf("failed to create once consumer: %w", err)
	}

	// Process interval jobs
	go s.processMessages(ctx, intervalCons, handler)
	// Process once jobs
	go s.processMessages(ctx, onceCons, handler)

	return nil
}

func (s *NATSJobSubscriber) processMessages(ctx context.Context, cons jetstream.Consumer, handler func(ctx context.Context, job *entity.Job) error) {
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

			var dto jobDTO
			if err := json.Unmarshal(msg.Data(), &dto); err != nil {
				s.log.Error("Failed to unmarshal job", zap.Error(err))
				msg.Nak()
				continue
			}

			job, err := dto.toEntity()
			if err != nil {
				s.log.Error("Failed to convert DTO to entity", zap.Error(err))
				msg.Nak()
				continue
			}

			if err := handler(ctx, job); err != nil {
				s.log.Error("Handler failed", zap.Error(err))
				msg.Nak()
				continue
			}

			msg.Ack()
		}
	}
}

// jobDTO is a JSON-serializable representation of a job
type jobDTO struct {
	ID             string  `json:"id"`
	Kind           int     `json:"kind"`
	Status         string  `json:"status"`
	Interval       *string `json:"interval,omitempty"`
	Once           *int64  `json:"once,omitempty"`
	LastFinishedAt int64   `json:"lastFinishedAt"`
	Payload        any     `json:"payload"`
}

func (d *jobDTO) toEntity() (*entity.Job, error) {
	job := &entity.Job{
		ID:             d.ID,
		Kind:           entity.JobKind(d.Kind),
		Status:         entity.JobStatus(d.Status),
		LastFinishedAt: d.LastFinishedAt,
		Payload:        d.Payload,
	}

	if d.Interval != nil {
		duration, err := parseDuration(*d.Interval)
		if err != nil {
			return nil, fmt.Errorf("failed to parse interval: %w", err)
		}
		job.Interval = &duration
	}

	if d.Once != nil {
		job.Once = d.Once
	}

	return job, nil
}

func parseDuration(s string) (time.Duration, error) {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	return duration, nil
}
