package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"scheduler/worker/config"
	"scheduler/worker/internal/adapter/publisher"
	"scheduler/worker/internal/adapter/subscriber"
	"scheduler/worker/internal/cases"
	"syscall"

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

	workerCases := cases.NewWorker(nil, pub, log.Named("worker"))
	// Subscribe to jobs and handle them
	ctx := context.Background()
	if err := sub.Subscribe(ctx, workerCases.RunJob); err != nil {
		return fmt.Errorf("subscribe to jobs: %w", err)
	}

	// Keep the worker running
	log.Info("Worker started, waiting for jobs...")

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	return nil
}
