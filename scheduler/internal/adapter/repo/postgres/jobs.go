package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"scheduler/scheduler/internal/entity"
	"scheduler/scheduler/internal/port/repo"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	upsertJobsQuery = `
		INSERT INTO jobs (id, kind, status, interval_seconds, once_timestamp, last_finished_at, payload)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			interval_seconds = EXCLUDED.interval_seconds,
			last_finished_at = EXCLUDED.last_finished_at
	`
)

var _ repo.Jobs = (*JobsRepo)(nil)

type JobsRepo struct {
	pool *pgxpool.Pool
}

func NewJobsRepo(pool *pgxpool.Pool) *JobsRepo {
	return &JobsRepo{
		pool: pool,
	}
}

func (r *JobsRepo) Create(ctx context.Context, job *entity.Job) error {
	payloadJSON, err := json.Marshal(job.Payload)
	if err != nil {
		return err
	}

	var intervalSeconds *int64
	if job.Interval != nil {
		seconds := int64(job.Interval.Seconds())
		intervalSeconds = &seconds
	}

	query := `
		INSERT INTO jobs (id, kind, status, interval_seconds, once_timestamp, last_finished_at, payload)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO NOTHING
	`

	_, err = r.pool.Exec(ctx, query,
		job.ID,
		job.Kind,
		job.Status,
		sql.NullInt64{
			Int64: func() int64 {
				if intervalSeconds != nil {
					return *intervalSeconds
				}
				return 0
			}(),
			Valid: intervalSeconds != nil,
		},
		sql.NullInt64{
			Int64: func() int64 {
				if job.Once != nil {
					return *job.Once
				}
				return 0
			}(),
			Valid: job.Once != nil,
		},
		job.LastFinishedAt,
		payloadJSON,
	)

	return err
}

func (r *JobsRepo) Read(ctx context.Context, jobID string) (*entity.Job, error) {
	var (
		id              string
		kind            int
		status          string
		intervalSeconds sql.NullInt64
		onceTimestamp   sql.NullInt64
		lastFinishedAt  int64
		payloadJSON     []byte
	)

	query := `
		SELECT id, kind, status, interval_seconds, once_timestamp, last_finished_at, payload
		FROM jobs
		WHERE id = $1
	`

	err := r.pool.QueryRow(ctx, query, jobID).Scan(
		&id,
		&kind,
		&status,
		&intervalSeconds,
		&onceTimestamp,
		&lastFinishedAt,
		&payloadJSON,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repo.ErrJobNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query row: %w", err)
	}

	var payload any
	if len(payloadJSON) > 0 {
		if err := json.Unmarshal(payloadJSON, &payload); err != nil {
			return nil, fmt.Errorf("unmarshal payload: %w", err)
		}
	}

	var interval *time.Duration
	if intervalSeconds.Valid {
		dur := time.Duration(intervalSeconds.Int64) * time.Second
		interval = &dur
	}

	var once *int64
	if onceTimestamp.Valid {
		once = &onceTimestamp.Int64
	}

	return &entity.Job{
		ID:             id,
		Kind:           entity.JobKind(kind),
		Status:         entity.JobStatus(status),
		Interval:       interval,
		Once:           once,
		LastFinishedAt: lastFinishedAt,
		Payload:        payload,
	}, nil
}

func (r *JobsRepo) List(ctx context.Context) ([]*entity.Job, error) {
	query := `
		SELECT id, kind, status, interval_seconds, once_timestamp, last_finished_at, payload
		FROM jobs
		ORDER BY id
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*entity.Job
	for rows.Next() {
		var (
			id              string
			kind            int
			status          string
			intervalSeconds sql.NullInt64
			onceTimestamp   sql.NullInt64
			lastFinishedAt  int64
			payloadJSON     []byte
		)

		if err := rows.Scan(
			&id,
			&kind,
			&status,
			&intervalSeconds,
			&onceTimestamp,
			&lastFinishedAt,
			&payloadJSON,
		); err != nil {
			return nil, err
		}

		var payload any
		if len(payloadJSON) > 0 {
			if err := json.Unmarshal(payloadJSON, &payload); err != nil {
				return nil, fmt.Errorf("unmarshal payload: %w", err)
			}
		}

		var interval *time.Duration
		if intervalSeconds.Valid {
			dur := time.Duration(intervalSeconds.Int64) * time.Second
			interval = &dur
		}

		var once *int64
		if onceTimestamp.Valid {
			once = &onceTimestamp.Int64
		}

		jobs = append(jobs, &entity.Job{
			ID:             id,
			Kind:           entity.JobKind(kind),
			Status:         entity.JobStatus(status),
			Interval:       interval,
			Once:           once,
			LastFinishedAt: lastFinishedAt,
			Payload:        payload,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if len(jobs) == 0 {
		return nil, repo.ErrJobNotFound
	}

	return jobs, nil
}

func (r *JobsRepo) Upsert(ctx context.Context, jobs []*entity.Job) error {
	if len(jobs) == 0 {
		return nil
	}

	for _, job := range jobs {
		payloadJSON, err := json.Marshal(job.Payload)
		if err != nil {
			return err
		}

		var intervalSeconds *int64
		if job.Interval != nil {
			seconds := int64(job.Interval.Seconds())
			intervalSeconds = &seconds
		}

		_, err = r.pool.Exec(ctx, upsertJobsQuery,
			job.ID,
			int(job.Kind),
			string(job.Status),
			intervalSeconds,
			job.Once,
			job.LastFinishedAt,
			payloadJSON,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
