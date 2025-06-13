package bus

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/pkg/bus/handler"
)

// CommandBus is a concrete implementation of CommandBus that stores handlers in memory.
type CommandBus struct {
	handlers handler.BusHandler
}

// NewCommandBus creates a new CommandBus.
func NewCommandBus() ICommandBus {
	return &CommandBus{
		handlers: handler.NewInMemoryBusHandler(),
	}
}

// RegisterHandler registers a CommandHandler for a specific Command type.
// It performs strict type checking to ensure the handler correctly implements CommandHandler[C].
func (b *CommandBus) RegisterHandler(cmdType ICommand, handler interface{}) error {
	cmdReflectType := reflect.TypeOf(cmdType)
	handlerReflectType := reflect.TypeOf(handler)

	if handler == nil {
		return fmt.Errorf("command handler for type %s cannot be nil", cmdReflectType.String())
	}

	// Check if a handler is already registered for this command type
	if duplicated := b.handlers.Store(cmdReflectType, handler); duplicated {
		return ErrCommandAlreadyRegistered{CommandType: cmdReflectType}
	}

	// Now perform stricter validation of the handler's type
	// The handler must have a method named "Handle" that matches the CommandHandler[C] signature.
	handleMethod, found := handlerReflectType.MethodByName("Handle")
	if !found {
		b.handlers.Delete(cmdReflectType) // Remove partially registered handler
		return ErrInvalidCommandHandler{
			HandlerType: handlerReflectType,
			CommandType: cmdReflectType,
			Reason:      "handler does not have a 'Handle' method",
		}
	}

	// Expected signature: func(ctx context.Context, cmd C) error
	// Method's signature (Input: receiver, ctx, cmd; Output: error)
	expectedIn := []reflect.Type{
		reflect.TypeOf((*context.Context)(nil)).Elem(), // ctx
		cmdReflectType, // cmd
	}
	expectedOut := []reflect.Type{
		reflect.TypeOf((*error)(nil)).Elem(), // error
	}

	// Check input parameters
	if handleMethod.Type.NumIn() != len(expectedIn)+1 { // +1 for the receiver
		b.handlers.Delete(cmdReflectType)
		return ErrInvalidCommandHandler{
			HandlerType: handlerReflectType,
			CommandType: cmdReflectType,
			Reason:      fmt.Sprintf("Handle method has %d input parameters, expected %d (excluding receiver)", handleMethod.Type.NumIn()-1, len(expectedIn)),
		}
	}
	for i, expected := range expectedIn {
		// handleMethod.Type.In(0) is the receiver itself
		if handleMethod.Type.In(i+1) != expected {
			b.handlers.Delete(cmdReflectType)
			return ErrInvalidCommandHandler{
				HandlerType: handlerReflectType,
				CommandType: cmdReflectType,
				Reason:      fmt.Sprintf("Handle method parameter %d type mismatch: expected %s, got %s", i+1, expected.String(), handleMethod.Type.In(i+1).String()),
			}
		}
	}

	// Check output parameters
	if handleMethod.Type.NumOut() != len(expectedOut) {
		b.handlers.Delete(cmdReflectType)
		return ErrInvalidCommandHandler{
			HandlerType: handlerReflectType,
			CommandType: cmdReflectType,
			Reason:      fmt.Sprintf("Handle method has %d output parameters, expected %d", handleMethod.Type.NumOut(), len(expectedOut)),
		}
	}
	for i, expected := range expectedOut {
		if handleMethod.Type.Out(i) != expected {
			b.handlers.Delete(cmdReflectType)
			return ErrInvalidCommandHandler{
				HandlerType: handlerReflectType,
				CommandType: cmdReflectType,
				Reason:      fmt.Sprintf("Handle method return parameter %d type mismatch: expected %s, got %s", i+1, expected.String(), handleMethod.Type.Out(i).String()),
			}
		}
	}

	return nil
}

// Dispatch dispatches a Command to its registered handler.
func (b *CommandBus) Dispatch(ctx context.Context, cmd ICommand) error {
	cmdType := reflect.TypeOf(cmd)

	handlerUntyped, ok := b.handlers.Load(cmdType)
	if !ok {
		return ErrNoCommandHandlerFound{CommandType: cmdType}
	}

	handlerVal := reflect.ValueOf(handlerUntyped)
	handleMethod := handlerVal.MethodByName("Handle")
	// The following check is now less critical as it's primarily done in RegisterHandler,
	// but kept for robustness against potential future misuses.
	if !handleMethod.IsValid() {
		return ErrInvalidCommandHandler{
			HandlerType: reflect.TypeOf(handlerUntyped),
			CommandType: cmdType,
			Reason:      "handler's 'Handle' method is invalid (should have been caught during registration)",
		}
	}

	args := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(cmd),
	}

	results := handleMethod.Call(args)

	if len(results) > 0 && !results[0].IsNil() {
		return results[0].Interface().(error)
	}

	return nil
}
