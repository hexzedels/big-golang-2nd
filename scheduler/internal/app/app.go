package app

import (
	"scheduler/scheduler/config"
	"scheduler/scheduler/internal/adapter/repo/postgres"
	"scheduler/scheduler/internal/cases"
	"scheduler/scheduler/internal/input/http/handler"
)

func Start(cfg config.Config) error {
	// TODO: Create jobs repo
	jobsRepo := postgres.NewJobsRepo() // TODO: pg config

	scheduler := cases.NewSchedulerCase(jobsRepo)
	handler.NewServer(scheduler)

	// TODO: Add http server.
	return nil
}
