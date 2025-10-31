package repo

import (
	"context"
	"errors"
	"scheduler/scheduler/internal/entity"
)

var (
	ErrJobNotFound = errors.New("job not found")
)

type Jobs interface {
	Upsert(ctx context.Context, jobs []*entity.Job) error
	Create(ctx context.Context, job *entity.Job) error
	Read(ctx context.Context, jobID string) (*entity.Job, error)
	List(ctx context.Context) ([]*entity.Job, error)
	// Delete
}
