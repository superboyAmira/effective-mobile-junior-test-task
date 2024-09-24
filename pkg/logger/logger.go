package logger

import (
	"log/slog"
	"os"
)

const (
	envLocal = "LOCAL"
	envDev   = "DEV"
	envProd  = "PROD"
)

func SetupLogger() *slog.Logger {
	var log *slog.Logger
	env := os.Getenv("LOG_LEVEL")
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	}

	return log
}
