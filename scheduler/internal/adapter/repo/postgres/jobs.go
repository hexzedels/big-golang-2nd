package postgres

import (
	"context"
	"errors"
	"fmt"
	"scheduler/scheduler/internal/port/repo"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repo.Jobs = (*JobsRepo)(nil)

type JobsRepo struct {
	db *pgxpool.Pool
}

func NewJobsRepo(db *pgxpool.Pool) *JobsRepo {
	return &JobsRepo{
		db: db,
	}
}

func (r *JobsRepo) Create(ctx context.Context, job *repo.JobDTO) error {
	query := `
		INSERT INTO jobs (id, once, interval, status, createdAt, lastFinishedAt, payLoad)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		job.ID,
		job.Once,
		job.Interval,
		job.Status,
		job.CreatedAt,
		job.LastFinishedAt,
		job.Payload,
	)

	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	return nil

}

func (r *JobsRepo) Read(ctx context.Context, jobID string) (*repo.JobDTO, error) {
	query := `
		SELECT id, once, interval, status, createdAt, lastFinishedAt, payLoad FROM jobs
		WHERE id = $1
	`
	job := &repo.JobDTO{}
	err := r.db.QueryRow(ctx, query, jobID).Scan(
		&job.ID,
		&job.Once,
		&job.Interval,
		&job.Status,
		&job.CreatedAt,
		&job.LastFinishedAt,
		&job.Payload,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("job not found: %s", jobID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read job: %w", err)
	}

	return job, nil
}

func (r *JobsRepo) Update(ctx context.Context, jobID string, job *repo.JobDTO) error {
	query := `
		UPDATE jobs 
		SET once = $1, interval = $2, status = $3, lastFinishedAt = $4, payLoad = $5
		WHERE id = $6
	`
	res, err := r.db.Exec(ctx, query,
		job.Once,
		job.Interval,
		job.Status,
		job.LastFinishedAt,
		job.Payload,
		jobID,
	)
	if err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("job not found: %s", jobID)
	}

	return nil
}

func (r *JobsRepo) Delete(ctx context.Context, jobID string) error {
	query := `
		DELETE FROM jobs WHERE id = $1
	`
	res, err := r.db.Exec(ctx, query, jobID)
	if err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("job not found: %s", jobID)
	}

	return nil
}
