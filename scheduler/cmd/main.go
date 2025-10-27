package main

import (
	"log"
	"scheduler/scheduler/config"
	"scheduler/scheduler/internal/app"
)

func main() {
	// TODO: config

	// if err := app.Start(config.Config{}); err != nil {
	// 	panic(err)
	// }
	cfg := config.Load()
	if err := app.Start(cfg); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}

}
