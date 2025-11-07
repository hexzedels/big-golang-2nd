package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"scheduler/scheduler/config"
	migrations "scheduler/scheduler/pkg/migration/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresPool(ctx context.Context, cfg config.PGConfig) (*pgxpool.Pool, error) {
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := migrations.Migrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	if err = db.Close(); err != nil {
		return nil, fmt.Errorf("failed to close database: %w", err)
	}

	pgxConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	pgxConfig.MaxConns = cfg.MaxOpenConns
	pgxConfig.MaxConnIdleTime = cfg.MaxIdleTime
	pgxConfig.MaxConnLifetime = cfg.MaxLifetime

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	return pool, nil
}
