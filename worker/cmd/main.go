package main

import (
	"scheduler/worker/config"
	"scheduler/worker/internal/app"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	if err := app.Start(cfg); err != nil {
		panic(err)
	}
}
