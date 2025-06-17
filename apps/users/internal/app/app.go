package app

import (
	"github.com/google/wire"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/command"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/query"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
)

type AppModule struct {
	UserQuery   query.UserQueryHandler
	UserCommand command.UserCommandHandler
}

func NewAppModule(
	uq query.UserQueryHandler,
	uc command.UserCommandHandler,
	userRepo user.UserRepository,
) *AppModule {
	return &AppModule{
		UserQuery:   uq,
		UserCommand: uc,
	}
}

var AppModuleSet = wire.NewSet(
	command.NewUserCommandHandler,
	query.NewUserQueryHandler,
	NewAppModule,
)
