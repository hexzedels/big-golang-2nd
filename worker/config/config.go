package config

import (
	"fmt"
	"os"
)

type Config struct {
	NATSURL string `yaml:"nats_url"`
}

func LoadConfig() (Config, error) {
	cfg := Config{}

	// Parse NATSURL
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		return cfg, fmt.Errorf("NATS_URL environment variable is required")
	}
	cfg.NATSURL = natsURL

	return cfg, nil
}
