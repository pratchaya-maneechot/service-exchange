package app

import (
	"github.com/google/wire"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/command"
	cmdHandles "github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/command/handles"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/query"
	qryHandles "github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/query/handles"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/common"
)

type AppModule struct {
	QueryBus   common.QueryBus
	CommandBus common.CommandBus
}

func NewAppModule(
	queryBus common.QueryBus,
	commandBus common.CommandBus,
	// Query handlers
	getUserProfileHandler *qryHandles.GetUserProfileHandler,

	// Command handlers
	lineRegisterHandler *cmdHandles.LineRegisterHandler,
	updateUserProfileHandler *cmdHandles.UpdateUserProfileHandler,
) *AppModule {
	// Register query handlers
	queryBus.Register(query.GetUserProfileQuery, getUserProfileHandler)

	// Register command handlers
	commandBus.Register(command.LineRegisterCommand, lineRegisterHandler)
	commandBus.Register(command.UpdateUserProfileCommand, updateUserProfileHandler)

	return &AppModule{
		QueryBus:   queryBus,
		CommandBus: commandBus,
	}
}

var AppModuleSet = wire.NewSet(
	qryHandles.NewGetUserProfileHandler,
	cmdHandles.NewLineRegisterHandler,
	cmdHandles.NewUpdateUserProfileHandler,
	NewAppModule,
)
