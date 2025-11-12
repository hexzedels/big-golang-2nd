package port

import (
	"context"
	"scheduler/scheduler/internal/entity"
)

type JobPublisher interface {
	Publish(
		ctx context.Context,
		kind entity.JobKind,
		exec *entity.Execution,
		jobPayload any,
	) error
}
