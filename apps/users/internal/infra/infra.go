package infra

import (
	"context"
	"fmt"
	"log/slog"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/google/wire"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/postgres"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/readers"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/repositories"
	"github.com/pratchaya-maneechot/service-exchange/libs/infra/observability"
	lp "github.com/pratchaya-maneechot/service-exchange/libs/infra/postgres"
)

type Infra struct {
	dbPool *lp.DBPool
	logger *slog.Logger
	tracer *sdktrace.TracerProvider
}

func NewInfra(
	dbPool *lp.DBPool,
	logger *slog.Logger,
	tracer *sdktrace.TracerProvider,
) *Infra {
	return &Infra{
		dbPool,
		logger,
		tracer,
	}
}

func (i *Infra) Close(ctx context.Context) error {
	var errs []error
	if i.tracer != nil {
		if err := i.tracer.Shutdown(ctx); err != nil {
			i.logger.Error("Failed to shutdown OpenTelemetry TracerProvider", "error", err)
			errs = append(errs, fmt.Errorf("tracer shutdown failed: %w", err))
		}
	}
	if i.dbPool != nil {
		i.dbPool.Close()
		i.logger.Info("PostgreSQL database pool closed.")
	}
	i.logger.Info("Infrastructure components closed.")

	if len(errs) > 0 {
		return fmt.Errorf("errors during infra close: %v", errs)
	}
	return nil
}

func ProvideMetricServer(cfg *config.Config) *observability.MetricServer {
	return observability.NewMetricServer(observability.MetricConfig{
		Path:    cfg.Metrics.Path,
		Addr:    cfg.Metrics.Address,
		Enabled: cfg.Metrics.Enabled,
	})
}

func ProvideLogger(cfg *config.Config) *slog.Logger {
	return observability.NewLogger(observability.LoggerConfig{
		Level:     cfg.Logging.Level,
		Format:    cfg.Logging.Format,
		AddSource: !cfg.IsDevelopment(),
	})
}

func ProvideTracer(ctx context.Context, cfg *config.Config) (*sdktrace.TracerProvider, error) {
	return observability.NewTracer(ctx, observability.TracerConfig{
		Name:        cfg.Name,
		Version:     cfg.Version,
		Environment: cfg.Environment,
	})
}

func ProvideMetricRecorder() observability.MetricsRecorder {
	return observability.NewPrometheusMetricsRecorder()
}

var InfraModuleSet = wire.NewSet(
	postgres.NewDBConn,
	repositories.NewPostgresUserRepository,
	readers.NewPostgresRoleReader,
	ProvideMetricServer,
	ProvideMetricRecorder,
	ProvideLogger,
	ProvideTracer,
	NewInfra,
)
