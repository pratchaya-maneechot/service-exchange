package bus

import (
	"context"
	"fmt"
	"reflect"
)

// ICommand is an interface that all commands must implement.
// It's a marker interface, no methods are required.
type ICommand interface{}

// ICommandHandler is an interface for handling a specific command.
type ICommandHandler[C ICommand] interface {
	Handle(ctx context.Context, cmd C) error
}

// ICommandBus is the interface for dispatching commands.
type ICommandBus interface {
	Dispatch(ctx context.Context, cmd ICommand) error
	RegisterHandler(cmdType ICommand, handler interface{}) error
}

// ErrNoCommandHandlerFound is returned when no handler is registered for a command.
type ErrNoCommandHandlerFound struct {
	CommandType reflect.Type
}

func (e ErrNoCommandHandlerFound) Error() string {
	return fmt.Sprintf("no command handler found for command type: %s", e.CommandType.String())
}

// ErrInvalidCommandHandler is returned when a registered handler does not implement ICommandHandler correctly.
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
