package port

import (
	"context"
	"scheduler/scheduler/internal/entity"
)

type JobPublisher interface {
	Publish(ctx context.Context, job *entity.Job) error
}
