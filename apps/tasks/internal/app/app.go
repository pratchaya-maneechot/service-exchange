package app

import (
	"github.com/google/wire"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/app/command"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/app/query"
)

type AppModule struct {
	TaskQuery   *query.TaskQueryService
	TaskCommand *command.TaskCommandService
}

func NewAppModule(
	tq *query.TaskQueryService,
	tc *command.TaskCommandService,
) *AppModule {
	return &AppModule{
		TaskQuery:   tq,
		TaskCommand: tc,
	}
}

var AppModuleSet = wire.NewSet(
	command.NewTaskCommandService,
	query.NewTaskQueryService,
	NewAppModule,
)
