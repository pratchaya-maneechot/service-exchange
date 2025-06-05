//go:build wireinject
// +build wireinject

package internal

import (
	"log/slog"

	"github.com/google/wire"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra"
	logger "github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/observability/logging"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/interfaces/grpc"
)

// Provider functions
func ProvideConfig() (*config.Config, error) {
	return config.Load()
}

func ProvideLogger(cfg *config.Config) *slog.Logger {
	return logger.New(cfg.Logging)
}

func ProvideShutdownHandler(logger *slog.Logger, cfg *config.Config) *ShutdownHandler {
	return NewShutdownHandler(logger, cfg.Server.ShutdownTimeout)
}

type App struct {
	Config   *config.Config
	Server   *grpc.Server
	Infra    *infra.InfraModule
	Logger   *slog.Logger
	App      *app.AppModule
	Shutdown *ShutdownHandler
}

func NewApp(
	config *config.Config,
	server *grpc.Server,
	infra *infra.InfraModule,
	appModule *app.AppModule,
	shutdown *ShutdownHandler,
	logger *slog.Logger,
) App {
	return App{
		Config:   config,
		Server:   server,
		Infra:    infra,
		App:      appModule,
		Shutdown: shutdown,
		Logger:   logger,
	}
}

// InitializeApp initializes the entire application
func InitializeApp() (App, error) {
	wire.Build(
		ProvideConfig,
		ProvideLogger,
		ProvideShutdownHandler,
		infra.InfraModuleSet,
		app.AppModuleSet,
		grpc.NewServer,
		NewApp,
	)
	return App{}, nil
}
