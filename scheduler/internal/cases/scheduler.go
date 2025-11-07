package cases

import (
	"context"
	"fmt"
	"scheduler/scheduler/internal/entity"
	"scheduler/scheduler/internal/port"
	"scheduler/scheduler/internal/port/repo"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SchedulerCase struct {
	jobsRepo  repo.Jobs
	running   map[string]*entity.RunningJob
	publisher port.JobPublisher
	interval  time.Duration
	mx        sync.Mutex
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
				r.logger.Error("process tick", zap.Error(err))
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (r *SchedulerCase) tick(ctx context.Context) error {
	// We fetch all job
	r.mx.Lock()
	defer r.mx.Unlock()

	jobs, err := r.jobsRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("list jobs: %w", err)
	}

	// We compare running and jobs from repo
	// We delete from running jobs deleted from repo

	repoJobs := make(map[string]*entity.Job, len(jobs))
	for _, j := range jobs {
		repoJobs[j.ID] = j
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
			if j.Once != nil {
				// Need to run once.
				if j.LastFinishedAt == 0 && now > *j.Once {
					go r.runJob(ctx, j)
				}
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

	if r.publisher != nil {
		if err := r.publisher.Publish(ctx, j); err != nil {
			r.logger.Error("publish job", zap.Error(err))
		}
	}
}

// HandleJobCompletion handles job completion messages from workers
func (r *SchedulerCase) HandleJobCompletion(ctx context.Context, jobID string, status string, finishedAt int64) error {
	// Read the job from repository
	job, err := r.jobsRepo.Read(ctx, jobID)
	if err != nil {
		return fmt.Errorf("read job: %w", err)
	}

	// Update job status and last finished time
	switch status {
	case entity.JobStatusCompleted:
		job.Status = entity.JobStatusCompleted
	case entity.JobStatusFailed:
		job.Status = entity.JobStatusFailed
	default:
		return fmt.Errorf("unknown status: %s", status)
	}

	job.LastFinishedAt = finishedAt

	// Remove from running jobs
	r.mx.Lock()
	defer r.mx.Unlock()

	if runningJob, ok := r.running[jobID]; ok {
		runningJob.Cancel()
		delete(r.running, jobID)
	}

	// Update job in repository
	if err := r.jobsRepo.Upsert(ctx, []*entity.Job{job}); err != nil {
		return fmt.Errorf("upsert job: %w", err)
	}

	r.logger.Info("Job completion handled",
		zap.String("job_id", jobID),
		zap.String("status", status),
		zap.Int64("finished_at", finishedAt))

	return nil
}
