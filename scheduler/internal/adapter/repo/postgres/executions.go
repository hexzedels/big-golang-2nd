package postgres

import (
	"context"
	"scheduler/scheduler/internal/entity"
	"scheduler/scheduler/internal/port/repo"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	upsertOnConflit = `ON CONFLICT (id) DO UPDATE SET
			worker_id = EXCLUDED.worker_id,
			status = EXCLUDED.status,
			started_at = EXCLUDED.started_at,
			finished_at = EXCLUDED.finished_at,
			error = EXCLUDED.error
			`

	tableExecutions = "executions"

	columnID         = "id"
	columnJobID      = "job_id"
	columnWorkerID   = "worker_id"
	columnStatus     = "status"
	columnQueuedAt   = "queued_at"
	columnStartedAt  = "started_at"
	columnFinishedAt = "finished_at"
	columnError      = "error"
)

var (
	executionsCols = []string{
		columnID,
		columnJobID,
		columnWorkerID,
		columnStatus,
		columnQueuedAt,
		columnStartedAt,
		columnFinishedAt,
		columnError,
	}
)

var _ repo.Executions = (*ExecutionsRepo)(nil)

type ExecutionsRepo struct {
	pool *pgxpool.Pool
}

func NewExecutionsRepo(pool *pgxpool.Pool) *ExecutionsRepo {
	return &ExecutionsRepo{
		pool: pool,
	}
}

func (r *ExecutionsRepo) Upsert(ctx context.Context, exec *entity.Execution) error {
	query, args := sqlbuilder.PostgreSQL.NewInsertBuilder().InsertInto(tableExecutions).Cols(executionsCols...).
		Values(
			exec.ID,
			exec.JobID,
			exec.WorkerID,
			exec.Status,
			exec.QueuedAt,
			exec.StartedAt,
			exec.FinishedAt,
			exec.Error, // just for now
		).SQL(upsertOnConflit).Build()

	_, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *ExecutionsRepo) Read(ctx context.Context, filter *entity.ReadExecutionFilter) (*entity.Execution, error) {
	panic("not implemented") // TODO: Implement
}

func (r *ExecutionsRepo) List(ctx context.Context, filter *entity.ListExecutionFilter) ([]*entity.Execution, error) {
	panic("not implemented") // TODO: Implement
}
