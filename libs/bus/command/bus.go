package command

import (
	"context"
	"fmt"
	"reflect"
)

type Command any
type Result any

type CommandHandler[C Command, R Result] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

type CommandBus interface {
	Dispatch(ctx context.Context, cmd Command) (Result, error)
	RegisterHandler(cmdType Command, handler any) error
}

type ErrNoCommandHandlerFound struct {
	CommandType reflect.Type
}

func (e ErrNoCommandHandlerFound) Error() string {
	return fmt.Sprintf("no command handler found for command type: %s", e.CommandType.String())
}

type ErrInvalidCommandHandler struct {
	HandlerType reflect.Type
	CommandType reflect.Type
	Reason      string
}

func (e ErrInvalidCommandHandler) Error() string {
	return fmt.Sprintf("invalid command handler type %s for command type %s: %s", e.HandlerType.String(), e.CommandType.String(), e.Reason)
}

type ErrCommandAlreadyRegistered struct {
	CommandType reflect.Type
}

func (e ErrCommandAlreadyRegistered) Error() string {
	return fmt.Sprintf("command handler already registered for command type: %s", e.CommandType.String())
}
