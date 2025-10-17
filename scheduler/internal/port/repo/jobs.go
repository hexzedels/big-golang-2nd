package repo

import (
	"context"
)

type Jobs interface {
	Create(ctx context.Context, job *JobDTO) error
	Read(ctx context.Context, jobID string) (*JobDTO, error)
	// Update
	// Delete
}

type JobDTO struct {
	// Interface specific entity-like struct
}
