package cases

import (
	"context"
	"fmt"
	"scheduler/worker/internal/entity"
	"scheduler/worker/internal/port"
	"time"

	"go.uber.org/zap"
)

type Worker struct {
	logger      *zap.Logger
	execChecker port.ExecutionsChecker
	publisher   port.CompletionPublisher
}

func NewWorker(
	execChecker port.ExecutionsChecker,
	pub port.CompletionPublisher,
	logger *zap.Logger,
) *Worker {
	return &Worker{
		logger:      logger,
		execChecker: execChecker,
		publisher:   pub,
	}
}

func (r *Worker) RunJob(ctx context.Context, exec *entity.JobExecution) error {
	logger := r.logger.With(
		zap.String("job_id", exec.JobID),
		zap.String("execution_id", exec.ExecutionID),
	)

	logger.Info("checking for already running execution")
	needRun, err := r.execChecker.Check(ctx, exec.ExecutionID)
	if err != nil {
		return fmt.Errorf("check need run execution: %w", err)
	}

	if !needRun {
		return ErrExecutionInProgress
	}

	r.logger.Info("starting job execution")
	r.logger.Info("simulate worker")
	time.Sleep(10 * time.Second)
	r.logger.Info("finished job execution")

	return nil
}
