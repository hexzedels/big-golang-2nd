package port

import "context"

type ExecutionsChecker interface {
	Check(ctx context.Context, executionID string) (needRun bool, err error)
}
