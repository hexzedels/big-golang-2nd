package cases

import "errors"

var (
	ErrExecutionInProgress = errors.New("job execution in progress")
)
