package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	DefaultMaxOpenConns      int32         = 25
	DefaultMaxLifetime       time.Duration = 5 * time.Minute
	DefaultMaxIdleTime       time.Duration = 1 * time.Minute
	DefaultSchedulerInterval time.Duration = 1 * time.Minute
	DefaultHTTPPort          string        = ":8090"
)

type Config struct {
	PG                PGConfig      `yaml:"pg"`
	SchedulerInterval time.Duration `yaml:"scheduler_interval"`
	NATSURL           string        `yaml:"nats_url"`
	HTTPPort          string        `yaml:"http_port"`
}

type PGConfig struct {
	DSN          string        `yaml:"dsn"`
	MaxOpenConns int32         `yaml:"max_open_conns"`
	MaxLifetime  time.Duration `yaml:"max_lifetime"`
	MaxIdleTime  time.Duration `yaml:"max_idle_time"`
}

func LoadConfig() (Config, error) {
	cfg := Config{}

	// Load PostgreSQL configuration
	pgDSN := os.Getenv("PG_DSN")
	if pgDSN == "" {
		pgDSN = os.Getenv("DATABASE_URL")
	}
	if pgDSN == "" {
		return cfg, fmt.Errorf("PG_DSN or DATABASE_URL environment variable is required")
	}
	cfg.PG.DSN = pgDSN

	// Parse MaxOpenConns
	maxOpenConnsStr := os.Getenv("PG_MAX_OPEN_CONNS")
	if maxOpenConnsStr == "" {
		cfg.PG.MaxOpenConns = DefaultMaxOpenConns
	} else {
		maxOpenConns, err := strconv.ParseInt(maxOpenConnsStr, 10, 32)
		if err != nil {
			return cfg, fmt.Errorf("invalid PG_MAX_OPEN_CONNS: %w", err)
		}
		cfg.PG.MaxOpenConns = int32(maxOpenConns)
	}

	// Parse MaxLifetime
	maxLifetimeStr := os.Getenv("PG_MAX_LIFETIME")
	if maxLifetimeStr == "" {
		cfg.PG.MaxLifetime = DefaultMaxLifetime
	} else {
		maxLifetime, err := time.ParseDuration(maxLifetimeStr)
		if err != nil {
			return cfg, fmt.Errorf("invalid PG_MAX_LIFETIME: %w", err)
		}
		cfg.PG.MaxLifetime = maxLifetime
	}

	// Parse MaxIdleTime
	maxIdleTimeStr := os.Getenv("PG_MAX_IDLE_TIME")
	if maxIdleTimeStr == "" {
		cfg.PG.MaxIdleTime = DefaultMaxIdleTime
	} else {
		maxIdleTime, err := time.ParseDuration(maxIdleTimeStr)
		if err != nil {
			return cfg, fmt.Errorf("invalid PG_MAX_IDLE_TIME: %w", err)
		}
		cfg.PG.MaxIdleTime = maxIdleTime
	}

	// Parse SchedulerInterval
	rawInterval := os.Getenv("SCHEDULER_INTERVAL")
	if rawInterval == "" {
		cfg.SchedulerInterval = DefaultSchedulerInterval
	} else {
		interval, err := time.ParseDuration(rawInterval)
		if err != nil {
			return cfg, fmt.Errorf("invalid SCHEDULER_INTERVAL: %w", err)
		}
		cfg.SchedulerInterval = interval
	}

	// Parse NATSURL
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		return cfg, fmt.Errorf("NATS_URL environment variable is required")
	}
	cfg.NATSURL = natsURL

	// Parse HTTPPort
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		cfg.HTTPPort = DefaultHTTPPort
	} else {
		// Ensure port starts with ':' if not already present
		if len(httpPort) > 0 && httpPort[0] != ':' {
			cfg.HTTPPort = ":" + httpPort
		} else {
			cfg.HTTPPort = httpPort
		}
	}

	return cfg, nil
}
