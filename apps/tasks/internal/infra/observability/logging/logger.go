package logger

import (
	"log/slog"
	"os"

	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/config"
)

func New(cfg config.LoggingConfig) *slog.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:     parseLevel(cfg.Level),
		AddSource: true,
	}

	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
