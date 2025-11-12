package entity

type ExecutionStatus string

const (
	ExecutionStatusQueued    ExecutionStatus = "queued"
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusCompleted ExecutionStatus = "completed"
	ExecutionStatusFailed    ExecutionStatus = "failed"
)

type ExecutionDTO struct {
	ID       string          `json:"id"`
	JobID    string          `json:"job_id"`
	Status   ExecutionStatus `json:"status"`
	QueuedAt int64           `json:"queued_at"`
	Payload  any             `json:"payload,omitempty"`
}
