package postgres

import (
	"context"
	"scheduler/scheduler/internal/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// var _ repo.Jobs = (*JobsRepo)(nil)

type JobsRepo struct {
	_ pgxpool.Pool
}

func NewJobsRepo( /*config*/ ) *JobsRepo {
	return &JobsRepo{}
}

func (r *JobsRepo) Create(ctx context.Context, job *entity.Job) error {
	panic("not implemented")
	return nil
}

func (r *JobsRepo) Read(ctx context.Context, jobID string) (*entity.Job, error) {
	panic("not implemented")
	return nil, nil
}
