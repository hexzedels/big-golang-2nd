package publisher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"scheduler/scheduler/internal/entity"
	"scheduler/scheduler/internal/port"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

var _ port.JobPublisher = (*NATSJobPublisher)(nil)

type NATSJobPublisher struct {
	js     jetstream.JetStream
	stream jetstream.Stream
	log    *zap.Logger
}

func NewNATSJobPublisher(ctx context.Context, log *zap.Logger, natsURL string) (*NATSJobPublisher, error) {
	// Connect to NATS server
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	newJS, _ := jetstream.New(nc)

	stream, err := newJS.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "JOBS",
		Subjects: []string{"JOBS.*"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create or update stream: %w", err)
	}

	log.Info("Connected to NATS JetStream", zap.String("url", natsURL))

	return &NATSJobPublisher{
		js:     newJS,
		stream: stream,
		log:    log,
	}, nil
}

func (p *NATSJobPublisher) Publish(ctx context.Context, job *entity.Job) error {
	// Convert entity to JSON-serializable DTO
	dto := jobDTO{
		ID:             job.ID,
		Kind:           int(job.Kind),
		Status:         string(job.Status),
		LastFinishedAt: job.LastFinishedAt,
		Payload:        job.Payload,
	}

	// Convert interval if present
	if job.Interval != nil {
		durationStr := job.Interval.String()
		dto.Interval = &durationStr
	}

	// Convert once if present
	if job.Once != nil {
		dto.Once = job.Once
	}

	// Serialize job to JSON
	data, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	if job.Kind == entity.JobUndefined {
		return errors.New("undefined job kind")
	}

	// Construct subject based on job kind and status
	subject := p.subjectForJob(job)

	_, err = p.js.Publish(ctx, subject, data)
	if err != nil {
		p.log.Error("Failed to publish job",
			zap.String("job_id", job.ID),
			zap.String("subject", subject),
			zap.Error(err))
		return fmt.Errorf("failed to publish job to NATS: %w", err)
	}

	p.log.Info("Published job to NATS",
		zap.String("job_id", job.ID),
		zap.String("subject", subject))

	return nil
}

// jobDTO is a JSON-serializable representation of a job
type jobDTO struct {
	ID             string  `json:"id"`
	Kind           int     `json:"kind"` // JobKind as int
	Status         string  `json:"status"`
	Interval       *string `json:"interval,omitempty"` // duration as string
	Once           *int64  `json:"once,omitempty"`
	LastFinishedAt int64   `json:"lastFinishedAt"`
	Payload        any     `json:"payload"`
}

// subjectForJob constructs the NATS subject for a job based on its kind and status
func (p *NATSJobPublisher) subjectForJob(job *entity.Job) string {
	switch job.Kind {
	case entity.JobKindInterval:
		return fmt.Sprintf("JOBS.interval.%s", job.Status)
	case entity.JobKindOnce:
		return fmt.Sprintf("JOBS.once.%s", job.Status)
	}

	return ""
}
