package port

import (
	"context"
	"scheduler/worker/internal/entity"
)

type CompletionPublisher interface {
	PublishCompletion(ctx context.Context, completion *entity.JobCompletion) error
}
