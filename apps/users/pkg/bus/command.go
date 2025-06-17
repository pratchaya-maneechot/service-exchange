package bus

import (
	"context"
	"fmt"
	"reflect"
)

// Command is an interface that all commands must implement.
// It's a marker interface, no methods are required.
type Command any

// CommandHandler is an interface for handling a specific command.
type CommandHandler[C Command, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

// CommandBus is the interface for dispatching commands.
type CommandBus interface {
	Dispatch(ctx context.Context, cmd Command) (any, error)
	RegisterHandler(cmdType Command, handler any) error
}

// ErrNoCommandHandlerFound is returned when no handler is registered for a command.
type ErrNoCommandHandlerFound struct {
	CommandType reflect.Type
}

func (e ErrNoCommandHandlerFound) Error() string {
	return fmt.Sprintf("no command handler found for command type: %s", e.CommandType.String())
}

// ErrInvalidCommandHandler is returned when a registered handler does not implement CommandHandler correctly.
type ErrInvalidCommandHandler struct {
	HandlerType reflect.Type
	CommandType reflect.Type
	Reason      string
}

func (e ErrInvalidCommandHandler) Error() string {
	return fmt.Sprintf("invalid command handler type %s for command type %s: %s", e.HandlerType.String(), e.CommandType.String(), e.Reason)
}

// ErrCommandAlreadyRegistered is returned when a command handler is registered multiple times.
type ErrCommandAlreadyRegistered struct {
	CommandType reflect.Type
}

func (e ErrCommandAlreadyRegistered) Error() string {
	return fmt.Sprintf("command handler already registered for command type: %s", e.CommandType.String())
}
