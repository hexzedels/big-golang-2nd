package entity

type ExecutionStatus string

const (
	// [queued, running, completed, failed]
	ExecutionStatusQueued    ExecutionStatus = "queued"
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusCompleted ExecutionStatus = "completed"
	ExecutionStatusFailed    ExecutionStatus = "failed"
)

type Execution struct {
	ID         string
	JobID      string
	WorkerID   *string
	Status     ExecutionStatus
	QueuedAt   int64
	StartedAt  *int64
	FinishedAt *int64
	Error      *ExecutionError
}

type ExecutionError struct {
	// Code    ErrorCode
	Details string `json:"details"`
}

type ListExecutionFilter struct {
	IDs    []string
	JobIDs []string
	Status *ExecutionStatus
}

type ReadExecutionFilter struct {
	ID     string
	JobID  string
	Status *ExecutionStatus
}
