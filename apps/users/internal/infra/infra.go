package infra

import (
	"github.com/google/wire"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/observability/metrics"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/postgres"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/readers"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/repositories"
)

type Repository struct {
	User user.UserRepository
}
type Reader struct {
	Role role.RoleReader
}

func ProvideReader(
	role role.RoleReader,
) *Reader {
	return &Reader{
		Role: role,
	}
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
	Reader     *Reader
	Repository *Repository
	Metric     *metrics.MetricServer
}

func NewInfraModule(
	repository *Repository,
	reader *Reader,
	metric *metrics.MetricServer,
) *InfraModule {
	return &InfraModule{
		Repository: repository,
		Metric:     metric,
		Reader:     reader,
	}
}

var InfraModuleSet = wire.NewSet(
	postgres.NewDBConn,
	repositories.NewPostgresUserRepository,
	readers.NewPostgresRoleReader,
	ProvideReader,
	ProvideRepository,
	ProvideMetricServer,
	NewInfraModule,
)
