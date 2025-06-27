package observability

import (
	"context"
	"log/slog"
	"os"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"go.opentelemetry.io/otel/trace"
)

func NewLogger(cfg *config.Config) *slog.Logger {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level:     parseLevel(cfg.Logging.Level),
		AddSource: !cfg.IsDevelopment(),
	}
	switch cfg.Logging.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

func withTracer(ctx context.Context, logger *slog.Logger) *slog.Logger {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		return logger.With(
			"trace_id", spanCtx.TraceID().String(),
			"span_id", spanCtx.SpanID().String(),
		)
	}
	return logger
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

func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	var bLogger *slog.Logger
	if logger, ok := ctx.Value("logger").(*slog.Logger); ok {
		bLogger = logger
	} else {
		bLogger = slog.Default()
	}
	return withTracer(ctx, bLogger)
}
