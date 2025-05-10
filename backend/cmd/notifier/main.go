package main

import (
	"log"
	"log/slog"

	"github.com/lim-bo/calendar/backend/internal/notifier"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	server := notifier.New()
	defer server.Close()
	err := server.Run()
	if err != nil {
		log.Print(err)
	}
}
