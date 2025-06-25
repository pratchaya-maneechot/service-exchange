package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/wire"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/handler"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/observability"
)

type App struct {
	GetUserProfileQueryHandler      *handler.GetUserProfileQueryHandler
	RegisterUserCommandHandler      *handler.RegisterUserCommandHandler
	UpdateUserProfileCommandHandler *handler.UpdateUserProfileCommandHandler
	RoleCacheService                *role.RoleCacheService
}

func ProvideRoleCacheService(
	reader role.RoleReader,
	logger *slog.Logger,
	cfg *config.Config,
	parentCtx context.Context,
	metricsRecorder observability.MetricsRecorder,
) *role.RoleCacheService {
	refreshInterval := time.Duration(cfg.Server.CacheRefreshIntervalDay) * 24 * time.Hour
	rcm := role.NewRoleCacheService(reader, logger, metricsRecorder, refreshInterval)
	rcm.InitAndStartRefresh(parentCtx)
	return rcm
}

var AppModuleSet = wire.NewSet(
	handler.NewGetUserProfileQueryHandler,
	handler.NewRegisterUserCommandHandler,
	handler.NewUpdateUserProfileCommandHandler,
	ProvideRoleCacheService,
	wire.Struct(new(App), "*"),
)
