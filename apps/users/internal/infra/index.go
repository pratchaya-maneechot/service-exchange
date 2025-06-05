package infra

import (
	"github.com/google/wire"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/common"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/repositories"
	commonInfra "github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/common"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/observability/metrics"
	repoInfra "github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/repositories"
)

func ProvideRepository(
	userRepo repositories.UserRepository,
	lineProfileRepo repositories.LineProfileRepository,
) domain.Repository {
	return domain.Repository{
		User:        userRepo,
		LineProfile: lineProfileRepo,
	}
}

func ProvideMetricServer(cfg *config.Config) *metrics.MetricServer {
	return metrics.NewServer(cfg.Metrics)
}

type InfraModule struct {
	EventBus   common.EventBus
	QueryBus   common.QueryBus
	CommandBus common.CommandBus
	Repository domain.Repository
	Metric     *metrics.MetricServer
}

func NewInfraModule(
	eventBus common.EventBus,
	queryBus common.QueryBus,
	commandBus common.CommandBus,
	repository domain.Repository,
	metric *metrics.MetricServer,
) *InfraModule {
	return &InfraModule{
		EventBus:   eventBus,
		QueryBus:   queryBus,
		CommandBus: commandBus,
		Repository: repository,
		Metric:     metric,
	}
}

var InfraModuleSet = wire.NewSet(
	commonInfra.NewEventBus,
	commonInfra.NewQueryBus,
	commonInfra.NewCommandBus,
	repoInfra.NewUserRepository,
	repoInfra.NewLineProfileRepository,
	ProvideRepository,
	ProvideMetricServer,
	NewInfraModule,
)
