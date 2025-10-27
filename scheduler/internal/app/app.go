package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"scheduler/scheduler/config"
	"scheduler/scheduler/internal/adapter/repo/postgres"
	"scheduler/scheduler/internal/cases"
	"scheduler/scheduler/internal/input/http/handler"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Start(cfg config.Config) error {
	// TODO: Create jobs repo
	db, err := connectDB(cfg.DB)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	jobsRepo := postgres.NewJobsRepo(db) // TODO: pg config

	schedulerCase := cases.NewSchedulerCase(jobsRepo)

	httpHandler := handler.NewServer(schedulerCase)

	log.Printf("Starting HTTP server on :%s", cfg.HTTP.Port)
	return http.ListenAndServe(":"+cfg.HTTP.Port, httpHandler)
}

func connectDB(cfg config.DBConfig) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL")
	return db, nil
}
