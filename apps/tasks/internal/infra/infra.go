package infra

import (
	"github.com/google/wire"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/domain/task"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/infra/observability/metrics"
	repository "github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/infra/persistence/inmemory/repositories"
)

type Repository struct {
	Task task.TaskRepository
}

func ProvideRepository(
	tr task.TaskRepository,
) *Repository {
	return &Repository{
		Task: tr,
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
	repository.NewTaskRepository,
	ProvideRepository,
	ProvideMetricServer,
	NewInfraModule,
)
