package entity

import (
	"time"
)

type JobKind int

const (
	JobUndefined = iota
	JobKindInterval
	JobKindOnce
)

type JobStatus string

const (
	// [queued, running, completed, failed]
	JobStatusQueued    = "queued"
	JobStatusRunning   = "running"
	JobStatusCompleted = "completed"
	JobStatusFailed    = "failed"
)

type Job struct {
	ID             string
	Kind           JobKind
	Status         JobStatus
	Interval       *time.Duration
	Once           *int64
	LastFinishedAt int64
	Payload        any
}
