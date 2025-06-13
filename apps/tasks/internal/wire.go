//go:build wireinject
// +build wireinject

package internal

import (
	"log/slog"

	"github.com/google/wire"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/app"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/app/command"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/app/query"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/config"
	grpc "github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/grpc"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/infra"
	logger "github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/infra/observability/logging"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/pkg/bus"
)

func ProvideConfig() (*config.Config, error) {
	return config.Load()
}

func ProvideLogger(cfg *config.Config) *slog.Logger {
	return logger.New(cfg.Logging)
}

func ProvideShutdownHandler(logger *slog.Logger, cfg *config.Config) *ShutdownHandler {
	return NewShutdownHandler(logger, cfg.Server.ShutdownTimeout)
}

type Internal struct {
	Config   *config.Config
	Server   *grpc.Server
	Infra    *infra.InfraModule
	Logger   *slog.Logger
	App      *app.AppModule
	Bus      *bus.Bus
	Shutdown *ShutdownHandler
}

func NewInternal(
	cf *config.Config,
	gs *grpc.Server,
	inf *infra.InfraModule,
	app *app.AppModule,
	bus *bus.Bus,
	sd *ShutdownHandler,
	lg *slog.Logger,
) *Internal {

	bus.CommandBus.RegisterHandler(command.CreateTaskCommand{}, app.TaskCommand.HandleCreateTaskCommand)

	bus.QueryBus.RegisterHandler(query.GetTaskDetail{}, app.TaskQuery.HandleGetTaskDetail)

	return &Internal{
		Config:   cf,
		Server:   gs,
		Infra:    inf,
		App:      app,
		Bus:      bus,
		Shutdown: sd,
		Logger:   lg,
	}
}

func InitializeApp() (*Internal, error) {
	wire.Build(
		ProvideConfig,
		ProvideLogger,
		ProvideShutdownHandler,
		app.AppModuleSet,
		bus.BusModuleSet,
		grpc.NewServer,
		infra.InfraModuleSet,
		NewInternal,
	)
	return &Internal{}, nil
}
