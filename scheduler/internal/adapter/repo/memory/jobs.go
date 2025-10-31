package memory

import (
	"context"
	"scheduler/scheduler/internal/entity"
	"scheduler/scheduler/internal/port/repo"
	"sync"
)

var _ repo.Jobs = (*JobsRepo)(nil)

type JobsRepo struct {
	jobs map[string]*entity.Job
	mu   sync.RWMutex
}

func NewJobsRepo() *JobsRepo {
	return &JobsRepo{
		jobs: make(map[string]*entity.Job),
	}
}

func (r *JobsRepo) Create(ctx context.Context, job *entity.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.jobs[job.ID] = job

	return nil
}

func (r *JobsRepo) Read(ctx context.Context, jobID string) (*entity.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, ok := r.jobs[jobID]
	if !ok {
		return nil, repo.ErrJobNotFound
	}

	return job, nil
}

func (r *JobsRepo) List(ctx context.Context) ([]*entity.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]*entity.Job, 0, len(r.jobs))
	for _, j := range r.jobs {
		out = append(out, j)
	}

	if len(out) == 0 {
		return nil, repo.ErrJobNotFound
	}

	return out, nil
}

func (r *JobsRepo) Upsert(ctx context.Context, jobs []*entity.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, j := range jobs {
		r.jobs[j.ID] = j
	}

	return nil
}

