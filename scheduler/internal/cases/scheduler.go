package cases

import (
	"context"
	"fmt"
	"scheduler/scheduler/internal/entity"
	"scheduler/scheduler/internal/port"
	"scheduler/scheduler/internal/port/repo"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SchedulerCase struct {
	jobsRepo  repo.Jobs
	running   map[string]*entity.RunningJob
	publisher port.JobPublisher
	interval  time.Duration
	logger    *zap.Logger
}

func NewSchedulerCase(
	jobsRepo repo.Jobs,
	publisher port.JobPublisher,
	interval time.Duration,
	logger *zap.Logger,
) *SchedulerCase {
	return &SchedulerCase{
		jobsRepo:  jobsRepo,
		running:   make(map[string]*entity.RunningJob),
		publisher: publisher,
		interval:  interval,
		logger:    logger,
	}
}

func (r *SchedulerCase) Create(ctx context.Context, job *entity.Job) (string, error) {
	job.ID = uuid.NewString()

	return job.ID, r.jobsRepo.Create(ctx, job)
}

func (r *SchedulerCase) Start(ctx context.Context) error {
	for {
		select {
		case <-time.NewTicker(r.interval).C:
			if err := r.tick(ctx); err != nil {

			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (r *SchedulerCase) tick(ctx context.Context) error {
	// We fetch all job

	jobs, err := r.jobsRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("list jobs: %w", err)
	}

	// We compare running and jobs from repo
	// We delete from running jobs deleted from repo

	repoJobs := make(map[string]*entity.Job, len(jobs))
	for _, j := range jobs {
		repoJobs[j.ID] = &entity.Job{
			ID: j.ID,
		}
	}

	for jobID, j := range r.running {
		if _, ok := repoJobs[jobID]; !ok {
			r.logger.Debug("stop deleted job")
			j.Cancel()
			delete(r.running, jobID)
		}
	}

	// We start new jobs that have to be started

	now := time.Now().UnixMilli()

	var updates []*entity.Job

	for jobID, j := range repoJobs {
		if _, ok := r.running[jobID]; ok {
			r.logger.Debug("skip already running job")
			continue
		}

		if j.Kind == entity.JobKindInterval {
			if now > j.Interval.Milliseconds()+j.LastFinishedAt {
				go r.runJob(ctx, j)
			}
		} else {
			// Need to run once.
			if j.LastFinishedAt == 0 && now > *j.Once {
				go r.runJob(ctx, j)
			}
		}

		j.Status = entity.JobStatusQueued
		updates = append(updates, j)
	}

	// We put new jobs to running

	if err := r.jobsRepo.Upsert(ctx, updates); err != nil {
		return fmt.Errorf("upsert started jobs: %w", err)
	}

	return nil
}

func (r *SchedulerCase) runJob(ctx context.Context, j *entity.Job) {
	ctx, cancel := context.WithCancel(ctx)
	r.running[j.ID] = &entity.RunningJob{
		Job:    j,
		Cancel: cancel,
	}

	if err := r.publisher.Publish(ctx, j); err != nil {
		r.logger.Error("publish job", zap.Error(err))
	}
}
