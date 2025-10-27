package repo

import (
	"context"
)

type Jobs interface {
	Create(ctx context.Context, job *JobDTO) error
	Read(ctx context.Context, jobID string) (*JobDTO, error)
	Update(ctx context.Context, jobID string, job *JobDTO) error
	Delete(ctx context.Context, jobID string) error
}

type JobDTO struct {
	// Interface specific entity-like struct
	ID             string
	Once           *string
	Interval       *string
	Status         string
	CreatedAt      int64
	LastFinishedAt int64
	Payload        map[string]interface{}
}

type ExecutionDTO struct {
	ID         *string
	JobID      *string
	WorkerID   *string
	Status     *string
	StartedAt  *int64
	FinishedAt *int64
}
