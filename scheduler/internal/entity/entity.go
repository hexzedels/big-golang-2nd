package entity

import "scheduler/scheduler/internal/input/http/gen"

type JobStatus string

const (
	JobStatusQueued    JobStatus = "queued"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

type Job struct {
	ID             string                 `json:"id"`
	Once           *string                `json:"once,omitempty"`
	Interval       *string                `json:"interval,omitempty"`
	Status         JobStatus              `json:"status"`
	CreatedAt      int64                  `json:"createdAt"`
	LastFinishedAt int64                  `json:"lastFinishedAt"`
	Payload        map[string]interface{} `json:"payLoad"`
}

type Execution struct {
	ID         *string     `json:"id,omitempty"`
	JobID      *string     `json:"jobId,omitempty"`
	WorkerID   *string     `json:"workerId,omitempty"`
	Status     *gen.Status `json:"status,omitempty"`
	StartedAt  *int64      `json:"startedAt,omitempty"`
	FinishedAt *int64      `json:"finishedAt,omitempty"`
}
