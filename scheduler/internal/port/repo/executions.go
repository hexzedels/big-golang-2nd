package repo

import (
	"context"
	"scheduler/scheduler/internal/entity"
)

type Executions interface {
	Upsert(ctx context.Context, exec *entity.Execution) error
	Read(ctx context.Context, filter *entity.ReadExecutionFilter) (*entity.Execution, error)
	List(ctx context.Context, filter *entity.ListExecutionFilter) ([]*entity.Execution, error)
}
