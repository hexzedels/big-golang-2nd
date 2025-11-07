package app

import (
	"context"
	"fmt"
	"scheduler/scheduler/pkg/entity"
	"scheduler/worker/config"
	"scheduler/worker/internal/adapter/publisher"
	"scheduler/worker/internal/adapter/subscriber"
	"time"

	"go.uber.org/zap"
)

func Start(cfg config.Config) error {
	log, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("new zap logger: %w", err)
	}

	// Create NATS subscriber
	sub, err := subscriber.NewNATSJobSubscriber(context.Background(), log, cfg.NATSURL)
	if err != nil {
		return fmt.Errorf("new NATS subscriber: %w", err)
	}

	// Create NATS completion publisher
	pub, err := publisher.NewNATSCompletionPublisher(context.Background(), log, cfg.NATSURL)
	if err != nil {
		return fmt.Errorf("new NATS completion publisher: %w", err)
	}

	// Subscribe to jobs and handle them
	ctx := context.Background()
	if err := sub.Subscribe(ctx, func(ctx context.Context, job *entity.Job) error {
		// Log the received job
		log.Info("Received job",
			zap.String("job_id", job.ID),
			zap.String("kind", fmt.Sprintf("%d", job.Kind)),
			zap.String("status", string(job.Status)),
			zap.Any("payload", job.Payload))

		// Simulate job processing - just mark as completed
		// In a real implementation, this would do actual work
		completion := publisher.JobCompletion{
			JobID:      job.ID,
			Status:     "completed",
			FinishedAt: time.Now().UnixMilli(),
		}

		// Publish completion message back to NATS
		if err := pub.PublishCompletion(ctx, completion); err != nil {
			log.Error("Failed to publish completion", zap.Error(err))
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("subscribe to jobs: %w", err)
	}

	// Keep the worker running
	log.Info("Worker started, waiting for jobs...")
	select {}
}
