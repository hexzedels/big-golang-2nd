package cases

import (
	"context"
	"fmt"
	"scheduler/scheduler/internal/entity"
	"scheduler/scheduler/internal/port/repo"
	"time"

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
	// job.ID = uuid.NewString()

	// return job.ID, r.jobsRepo.Create(ctx, &repo.JobDTO{})
	if job.Once == nil && job.Interval == nil {
		return "", fmt.Errorf("'once' or 'interval' must be specified")
	}

	job.ID = uuid.NewString()
	job.Status = "queued"
	job.CreatedAt = time.Now().Unix()
	job.LastFinishedAt = 0

	jobDTO := &repo.JobDTO{
		ID:             job.ID,
		Once:           job.Once,
		Interval:       job.Interval,
		Status:         string(job.Status),
		CreatedAt:      job.CreatedAt,
		LastFinishedAt: job.LastFinishedAt,
		Payload:        job.Payload,
	}

	err := r.jobsRepo.Create(ctx, jobDTO)
	if err != nil {
		return "", fmt.Errorf("failed to create job: %w", err)
	}

	return job.ID, nil
}

func (r *SchedulerCase) Read(ctx context.Context, jobID string) (*entity.Job, error) {
	jobDTO, err := r.jobsRepo.Read(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to read job: %w", err)
	}

	job := &entity.Job{
		ID:             jobDTO.ID,
		Once:           jobDTO.Once,
		Interval:       jobDTO.Interval,
		Status:         entity.JobStatus(jobDTO.Status),
		CreatedAt:      jobDTO.CreatedAt,
		LastFinishedAt: jobDTO.LastFinishedAt,
		Payload:        jobDTO.Payload,
	}

	return job, nil

}

func (r *SchedulerCase) Update(ctx context.Context, jobID string, job *entity.Job) error {
	jobDTO := &repo.JobDTO{
		ID:             job.ID,
		Once:           job.Once,
		Interval:       job.Interval,
		Status:         string(job.Status),
		CreatedAt:      job.CreatedAt,
		LastFinishedAt: job.LastFinishedAt,
		Payload:        job.Payload,
	}
	err := r.jobsRepo.Update(ctx, jobID, jobDTO)
	if err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}

	return nil
}

func (r *SchedulerCase) Delete(ctx context.Context, jobID string) error {
	err := r.jobsRepo.Delete(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}
	return nil
}
