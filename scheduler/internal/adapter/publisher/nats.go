package publisher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"scheduler/scheduler/internal/entity"
	"scheduler/scheduler/internal/port"
	pkgentity "scheduler/scheduler/pkg/entity"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

var _ port.JobPublisher = (*NATSJobPublisher)(nil)

type NATSJobPublisher struct {
	js jetstream.JetStream
	// stream jetstream.Stream
	log *zap.Logger
}

func NewNATSJobPublisher(ctx context.Context, log *zap.Logger, natsURL string) (*NATSJobPublisher, error) {
	// Connect to NATS server
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	newJS, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("failed to create jetstream context: %w", err)
	}

	streams := newJS.ListStreams(ctx)
	for stream := range streams.Info() {
		log.Info("Stream", zap.Any("name", stream.Config))
	}

	log.Info("Connected to NATS JetStream", zap.String("url", natsURL))

	return &NATSJobPublisher{
		js:  newJS,
		log: log,
	}, nil
}

func (p *NATSJobPublisher) Publish(
	ctx context.Context,
	kind entity.JobKind,
	exec *entity.Execution,
	paylaod any,
) error {
	if kind == entity.JobUndefined {
		return errors.New("undefined job kind")
	}

	// Convert entity to JSON-serializable DTO
	dto := pkgentity.ExecutionDTO{
		ID:       exec.ID,
		JobID:    exec.JobID,
		Status:   pkgentity.ExecutionStatus(exec.Status),
		QueuedAt: exec.QueuedAt,
		Payload:  paylaod,
	}

	// Serialize job to JSON
	data, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	// Construct subject based on job kind and status
	subject := p.subjectForJob(kind, dto.Status)

	_, err = p.js.Publish(ctx, subject, data)
	if err != nil {
		p.log.Error("Failed to publish job",
			zap.String("execution_id", exec.ID),
			zap.String("job_id", exec.JobID),
			zap.String("subject", subject),
			zap.Error(err))
		return fmt.Errorf("failed to publish job to NATS: %w", err)
	}

	p.log.Info("Published job to NATS",
		zap.String("execution_id", exec.ID),
		zap.String("job_id", exec.JobID),
		zap.String("subject", subject))

	return nil
}

// subjectForJob constructs the NATS subject for a job based on its kind and status
func (p *NATSJobPublisher) subjectForJob(
	kind entity.JobKind,
	status pkgentity.ExecutionStatus,
) string {
	switch kind {
	case entity.JobKindInterval:
		return fmt.Sprintf("JOBS.interval.%s", status)
	case entity.JobKindOnce:
		return fmt.Sprintf("JOBS.once.%s", status)
	}

	return ""
}
