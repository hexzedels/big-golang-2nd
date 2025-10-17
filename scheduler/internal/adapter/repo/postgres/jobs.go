package postgres

import (
	"context"
	"scheduler/scheduler/internal/port/repo"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repo.Jobs = (*JobsRepo)(nil)

type JobsRepo struct {
	_ pgxpool.Pool
}

func NewJobsRepo( /*config*/ ) *JobsRepo {
	return &JobsRepo{}
}

func (r *JobsRepo) Create(ctx context.Context, job *repo.JobDTO) error {
	panic("not implemented")
	return nil
}

func (r *JobsRepo) Read(ctx context.Context, jobID string) (*repo.JobDTO, error) {
	panic("not implemented")
	return nil, nil
}
