package cases

import (
	"context"
	"scheduler/scheduler/internal/entity"
	"scheduler/scheduler/internal/port/repo"

	"github.com/google/uuid"
)

type SchedulerCase struct {
	jobsRepo repo.Jobs
}

func NewSchedulerCase(jobsRepo repo.Jobs) *SchedulerCase {
	return &SchedulerCase{
		jobsRepo: jobsRepo,
	}
}

func (r *SchedulerCase) Create(ctx context.Context, job *entity.Job) (string, error) {
	job.ID = uuid.NewString()

	return job.ID, r.jobsRepo.Create(ctx, &repo.JobDTO{})
}
