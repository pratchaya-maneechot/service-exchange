package common

import "context"

type Command interface {
	CommandName() string
}

type CommandHandler interface {
	Handle(ctx context.Context, command Command) (interface{}, error)
}

type CommandBus interface {
	Register(commandName string, handler CommandHandler)
	Send(ctx context.Context, command Command) (interface{}, error)
}
