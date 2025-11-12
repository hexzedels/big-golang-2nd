package app

import (
	"context"
	"fmt"
	"net/http"
	"scheduler/scheduler/config"
	"scheduler/scheduler/internal/adapter/publisher"
	"scheduler/scheduler/internal/adapter/repo/postgres"
	"scheduler/scheduler/internal/adapter/subscriber"
	"scheduler/scheduler/internal/cases"
	"scheduler/scheduler/internal/input/http/gen"
	"scheduler/scheduler/internal/input/http/handler"

	"go.uber.org/zap"
)

func Start(cfg config.Config) error {
	pgPool, err := postgres.NewPostgresPool(context.Background(), cfg.PG)
	if err != nil {
		return fmt.Errorf("new postgres pool: %w", err)
	}

	jobsRepo := postgres.NewJobsRepo(pgPool)
	execRepo := postgres.NewExecutionsRepo(pgPool)

	log, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("new zap logger: %w", err)
	}

	// Create NATS JetStream publisher
	pub, err := publisher.NewNATSJobPublisher(context.Background(), log, cfg.NATSURL)
	if err != nil {
		// Log error but continue with nil publisher (graceful degradation)
		log.Warn("Failed to create NATS publisher, continuing without publisher", zap.Error(err))
	}

	scheduler := cases.NewSchedulerCase(
		jobsRepo,
		execRepo,
		pub,
		cfg.SchedulerInterval,
		log,
	)
	srv := handler.NewServer(scheduler)

	// Start scheduler tick loop
	ctx := context.Background()
	go func() {
		if err := scheduler.Start(ctx); err != nil {
			log.Error("Scheduler tick loop failed", zap.Error(err))
		}
	}()

	// Create and start NATS completion subscriber
	completionSub, err := subscriber.NewNATSCompletionSubscriber(ctx, log, cfg.NATSURL)
	if err != nil {
		log.Warn("Failed to create NATS completion subscriber, continuing without subscriber", zap.Error(err))
	} else {
		if err := completionSub.Subscribe(ctx, func(ctx context.Context, completion subscriber.JobCompletion) error {
			return scheduler.HandleJobCompletion(ctx, completion.JobID, completion.Status, completion.FinishedAt)
		}); err != nil {
			log.Warn("Failed to subscribe to completion messages", zap.Error(err))
		}
	}

	h := gen.NewStrictHandler(srv, nil)
	r := gen.HandlerWithOptions(h, gen.ChiServerOptions{})

	return http.ListenAndServe(cfg.HTTPPort, r)
}
