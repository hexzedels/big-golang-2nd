package handler

import (
	"fmt"
	"scheduler/scheduler/internal/entity"
	"scheduler/scheduler/internal/input/http/gen"
	"time"
)

func toEntityJob(job *gen.JobCreate) (*entity.Job, error) {
	entityJob := new(entity.Job)

	if job.Interval != nil {
		interval, err := time.ParseDuration(*job.Interval)
		if err != nil {
			return nil, fmt.Errorf("parse interval duraton: %w", err)
		}
		entityJob.Interval = &interval
		entityJob.Kind = entity.JobKindInterval
	}

	entityJob.Payload = job.Payload

	return entityJob, nil
}
