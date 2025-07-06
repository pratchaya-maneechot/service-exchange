package infra

import (
	"context"
	"fmt"
	"log/slog"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/google/wire"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/observability"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/observability/metrics"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/postgres"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/readers"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/repositories"
	libPostgres "github.com/pratchaya-maneechot/service-exchange/libs/infra/postgres"
)

type Infra struct {
	dbPool *libPostgres.DBPool
	logger *slog.Logger
	tracer *sdktrace.TracerProvider
}

func NewInfra(
	dbPool *libPostgres.DBPool,
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

var InfraModuleSet = wire.NewSet(
	postgres.NewDBConn,
	repositories.NewPostgresUserRepository,
	readers.NewPostgresRoleReader,
	metrics.NewServer,
	observability.NewLogger,
	observability.NewTracer,
	observability.NewPrometheusMetricsRecorder,
	NewInfra,
)
