//go:build wireinject
// +build wireinject

package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/wire"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/command"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/query"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/grpc"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra"
	"github.com/pratchaya-maneechot/service-exchange/libs/bus"
	libGrpc "github.com/pratchaya-maneechot/service-exchange/libs/grpc"
	"github.com/pratchaya-maneechot/service-exchange/libs/infra/observability"
)

func ProvideConfig() (*config.Config, error) {
	return config.Load()
}

type Internal struct {
	Config       *config.Config
	Server       *libGrpc.GRPCServer
	Logger       *slog.Logger
	App          *app.App
	Bus          *bus.Bus
	MetricServer *observability.MetricServer
	Cleanup      func()
}

func NewInternal(
	cfg *config.Config,
	gs *libGrpc.GRPCServer,
	appModule *app.App,
	bBus *bus.Bus,
	logger *slog.Logger,
	metricServer *observability.MetricServer,
	cleanup func(),
) *Internal {

	bBus.CommandBus.RegisterHandler(command.RegisterUserCommand{}, appModule.RegisterUserCommandHandler)
	bBus.CommandBus.RegisterHandler(command.UpdateUserProfileCommand{}, appModule.UpdateUserProfileCommandHandler)
	bBus.QueryBus.RegisterHandler(query.GetUserProfileQuery{}, appModule.GetUserProfileQueryHandler)

	return &Internal{
		Config:       cfg,
		Server:       gs,
		App:          appModule,
		Bus:          bBus,
		Logger:       logger,
		MetricServer: metricServer,
		Cleanup:      cleanup,
	}
}

func InitializeApp(parentCtx context.Context) (*Internal, error) {
	wire.Build(
		ProvideConfig,
		app.AppModuleSet,
		bus.BusModuleSet,
		infra.InfraModuleSet,
		grpc.NewGRPCServer,
		NewInternal,
		ProvideAppCleanup,
	)
	return &Internal{}, nil
}

func ProvideAppCleanup(
	infraModule *infra.Infra,
	server *libGrpc.GRPCServer,
	metricServer *observability.MetricServer,
	logger *slog.Logger,
	appModule *app.App,
) func() {
	return func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var shutdownErrors []error

		if appModule.RoleCacheService != nil {
			appModule.RoleCacheService.Stop()
			logger.Info("RoleCacheService stopped.")
		}

		if metricServer != nil {
			if err := metricServer.Stop(cleanupCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
				shutdownErrors = append(shutdownErrors, fmt.Errorf("failed to stop metrics server: %w", err))
				logger.Error("Failed to stop metrics server", "error", err)
			} else if err == nil {
				logger.Info("Metrics server stopped.")
			}
		}

		// gRPC server is typically stopped by context cancellation in its Start() method,
		// so an explicit server.Stop() might not be needed here if Start() is blocking.
		// If you had a separate server.Stop() method, you'd call it here.

		if err := infraModule.Close(cleanupCtx); err != nil {
			shutdownErrors = append(shutdownErrors, fmt.Errorf("failed to stop infra: %w", err))
			logger.Error("Failed to stop infrastructure", "error", err)
		} else {
			logger.Info("Infrastructure components closed.")
		}

		if len(shutdownErrors) > 0 {
			logger.Error("Errors encountered during application shutdown", "errors", shutdownErrors)
		} else {
			logger.Info("Application shutdown completed successfully.")
		}
	}
}
