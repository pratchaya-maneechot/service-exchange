package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/wire"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/command"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/query"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	"github.com/pratchaya-maneechot/service-exchange/libs/infra/observability"
)

type App struct {
	GetUserProfileQueryHandler      *query.GetUserProfileQueryHandler
	RegisterUserCommandHandler      *command.RegisterUserCommandHandler
	UpdateUserProfileCommandHandler *command.UpdateUserProfileCommandHandler
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
	query.NewGetUserProfileQueryHandler,
	command.NewRegisterUserCommandHandler,
	command.NewUpdateUserProfileCommandHandler,
	ProvideRoleCacheService,
	wire.Struct(new(App), "*"),
)
