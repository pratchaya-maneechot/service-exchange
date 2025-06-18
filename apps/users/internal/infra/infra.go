package infra

import (
	"github.com/google/wire"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/observability/metrics"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/postgres"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/postgres/repositories"
)

type Repository struct {
	User user.UserRepository
}

func ProvideRepository(
	ur user.UserRepository,
) *Repository {
	return &Repository{
		User: ur,
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
	postgres.NewDBConn,
	repositories.NewPostgresUserRepository,
	ProvideRepository,
	ProvideMetricServer,
	NewInfraModule,
)
