package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/wire"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/handler"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
)

type AppModule struct {
	GetUserProfileQueryHandler      *handler.GetUserProfileQueryHandler
	RegisterUserCommandHandler      *handler.RegisterUserCommandHandler
	UpdateUserProfileCommandHandler *handler.UpdateUserProfileCommandHandler
}

func NewAppModule(
	gpq *handler.GetUserProfileQueryHandler,
	ruc *handler.RegisterUserCommandHandler,
	upc *handler.UpdateUserProfileCommandHandler,
) *AppModule {
	return &AppModule{
		RegisterUserCommandHandler:      ruc,
		GetUserProfileQueryHandler:      gpq,
		UpdateUserProfileCommandHandler: upc,
	}
}

func ProvideRoleCacheService(
	reader role.RoleReader,
	logger *slog.Logger,
	cfg *config.Config,
	ctx context.Context,
) *role.RoleCacheService {
	refreshInterval := time.Duration(cfg.Server.CacheRefreshIntervalDay) * 24 * time.Hour
	rcm := role.NewRoleCacheService(reader, logger, refreshInterval)
	rcm.InitAndStartRefresh(ctx)
	return rcm
}

var AppModuleSet = wire.NewSet(
	handler.NewGetUserProfileQueryHandler,
	handler.NewRegisterUserCommandHandler,
	handler.NewUpdateUserProfileCommandHandler,
	ProvideRoleCacheService,
	NewAppModule,
)
