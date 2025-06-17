package infra

import (
	"github.com/google/wire"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/observability/metrics"
	repository "github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/inmemory/repositories"
)

type Repository struct {
	User user.UserRepository
}

func ProvideRepository(
	tr user.UserRepository,
) *Repository {
	return &Repository{
		User: tr,
	}
}

func ProvideMetricServer(cfg *config.Config) *metrics.MetricServer {
	return metrics.NewServer(cfg.Metrics)
}

type InfraModule struct {
	Repository *Repository
	Metric     *metrics.MetricServer
}

func NewInfraModule(
	repository *Repository,
	metric *metrics.MetricServer,
) *InfraModule {
	return &InfraModule{
		Repository: repository,
		Metric:     metric,
	}
}

var InfraModuleSet = wire.NewSet(
	repository.NewUserRepository,
	ProvideRepository,
	ProvideMetricServer,
	NewInfraModule,
)
