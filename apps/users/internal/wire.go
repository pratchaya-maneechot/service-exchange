//go:build wireinject
// +build wireinject

package internal

import (
	"context"
	"log/slog"

	"github.com/google/wire"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/command"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/query"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	grpc "github.com/pratchaya-maneechot/service-exchange/apps/users/internal/grpc"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra"
	logger "github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/observability/logging"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/pkg/bus"
)

func ProvideConfig() (*config.Config, error) {
	return config.Load()
}

func ProvideLogger(cfg *config.Config) *slog.Logger {
	return logger.New(cfg)
}

type Internal struct {
	Config   *config.Config
	Server   *grpc.Server
	Infra    *infra.InfraModule
	Logger   *slog.Logger
	App      *app.AppModule
	Bus      *bus.Bus
	Shutdown *ShutdownHandler
	BContext context.Context
}

func NewInternal(
	cf *config.Config,
	gs *grpc.Server,
	inf *infra.InfraModule,
	app *app.AppModule,
	bus *bus.Bus,
	sd *ShutdownHandler,
	lg *slog.Logger,
	bCtx context.Context,
) *Internal {

	bus.CommandBus.RegisterHandler(command.RegisterUserCommand{}, app.RegisterUserCommandHandler)

	bus.QueryBus.RegisterHandler(query.GetUserProfileQuery{}, app.GetUserProfileQueryHandler)

	return &Internal{
		Config:   cf,
		Server:   gs,
		Infra:    inf,
		App:      app,
		Bus:      bus,
		Shutdown: sd,
		Logger:   lg,
		BContext: bCtx,
	}
}

func InitializeApp(ctx context.Context) (*Internal, error) {
	wire.Build(
		ProvideConfig,
		ProvideLogger,
		NewShutdownHandler,
		app.AppModuleSet,
		bus.BusModuleSet,
		grpc.NewServer,
		infra.InfraModuleSet,
		NewInternal,
	)
	return &Internal{}, nil
}
