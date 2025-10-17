package main

import (
	"scheduler/scheduler/config"
	"scheduler/scheduler/internal/app"
)

func main() {
	// TODO: config

	if err := app.Start(config.Config{}); err != nil {
		panic(err)
	}
}
