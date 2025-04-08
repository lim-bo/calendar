package main

import (
	"log"
	"log/slog"

	"github.com/lim-bo/calendar/backend/internal/api"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	api := api.New()
	api.MountEndpoint()
	log.Fatal(api.Run())
}
