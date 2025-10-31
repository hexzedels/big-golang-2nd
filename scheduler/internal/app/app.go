package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"scheduler/scheduler/config"
	"scheduler/scheduler/internal/adapter/publisher"
	"scheduler/scheduler/internal/adapter/repo/memory"
	"scheduler/scheduler/internal/adapter/repo/postgres"
	"scheduler/scheduler/internal/cases"
	"scheduler/scheduler/internal/input/http/gen"
	"scheduler/scheduler/internal/input/http/handler"
	"scheduler/scheduler/internal/port/repo"
	"time"

	"go.uber.org/zap"
)

func Start(cfg config.Config) error {
	// TODO: Create jobs repo
	// Use in-memory repo by default to enable running without external deps
	var jobsRepo repo.Jobs
	_ = postgres.NewJobsRepo // keep import for later postgres wiring
	jobsRepo = memory.NewJobsRepo()

	rawInterval := os.Getenv("SCHEDULER_INTERVAL")
	interval, err := time.ParseDuration(rawInterval)
	if err != nil {
		return fmt.Errorf("parse interval: %w", err)
	}

	log, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("new zap logger: %w", err)
	}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		return fmt.Errorf("NATS_URL environment variable is not set")
	}

	// Create NATS JetStream publisher
	pub, err := publisher.NewNATSJobPublisher(context.Background(), log, natsURL)
	if err != nil {
		// Log error but continue with nil publisher (graceful degradation)
		log.Warn("Failed to create NATS publisher, continuing without publisher", zap.Error(err))
	}

	scheduler := cases.NewSchedulerCase(jobsRepo, pub, interval, log)
	srv := handler.NewServer(scheduler)

	h := gen.NewStrictHandler(srv, nil)
	r := gen.HandlerWithOptions(h, gen.ChiServerOptions{})

	return http.ListenAndServe(":8090", r)
}
