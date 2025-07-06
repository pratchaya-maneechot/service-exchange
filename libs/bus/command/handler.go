package command

import (
	"context"
	"fmt"
	"reflect"
)

type CommandBusHandler interface {
	Load(key any) (value any, ok bool)
	Store(key any, value any) (duplicated bool)
	Delete(key any)
}

type commandBus struct {
	handlers CommandBusHandler
}

func NewCommandBus(h CommandBusHandler) CommandBus {
	return &commandBus{
		handlers: h,
	}
}

func (b *commandBus) RegisterHandler(cmdType Command, handler any) error {
	cmdReflectType := reflect.TypeOf(cmdType)
	handlerReflectType := reflect.TypeOf(handler)

	if handler == nil {
		panic(fmt.Errorf("command handler for type %s cannot be nil", cmdReflectType.String()))
	}

	if duplicated := b.handlers.Store(cmdReflectType, handler); duplicated {
		panic(ErrCommandAlreadyRegistered{CommandType: cmdReflectType})
	}

	handleMethod, found := handlerReflectType.MethodByName("Handle")
	if !found {
		b.handlers.Delete(cmdReflectType)
		panic(ErrInvalidCommandHandler{
			HandlerType: handlerReflectType,
			CommandType: cmdReflectType,
			Reason:      "handler does not have a 'Handle' method",
		})
	}

	expectedIn := []reflect.Type{
		reflect.TypeOf((*context.Context)(nil)).Elem(),
		cmdReflectType,
	}
	expectedOut := []reflect.Type{
		nil,
		reflect.TypeOf((*error)(nil)).Elem(),
	}

	if handleMethod.Type.NumIn() != len(expectedIn)+1 {
		b.handlers.Delete(cmdReflectType)
		panic(ErrInvalidCommandHandler{
			HandlerType: handlerReflectType,
			CommandType: cmdReflectType,
			Reason:      fmt.Sprintf("Handle method has %d input parameters, expected %d (excluding receiver)", handleMethod.Type.NumIn()-1, len(expectedIn)),
		})
	}
	for i, expected := range expectedIn {
		if handleMethod.Type.In(i+1) != expected {
			b.handlers.Delete(cmdReflectType)
			panic(ErrInvalidCommandHandler{
				HandlerType: handlerReflectType,
				CommandType: cmdReflectType,
				Reason:      fmt.Sprintf("Handle method parameter %d type mismatch: expected %s, got %s", i+1, expected.String(), handleMethod.Type.In(i+1).String()),
			})
		}
	}

	if handleMethod.Type.NumOut() != len(expectedOut) {
		b.handlers.Delete(cmdReflectType)
		panic(ErrInvalidCommandHandler{
			HandlerType: handlerReflectType,
			CommandType: cmdReflectType,
			Reason:      fmt.Sprintf("Handle method has %d output parameters, expected %d", handleMethod.Type.NumOut(), len(expectedOut)),
		})
	}
	for i, expected := range expectedOut {
		if expected != nil && handleMethod.Type.Out(i) != expected {
			b.handlers.Delete(cmdReflectType)
			panic(ErrInvalidCommandHandler{
				HandlerType: handlerReflectType,
				CommandType: cmdReflectType,
				Reason:      fmt.Sprintf("Handle method return parameter %d type mismatch: expected %s, got %s", i+1, expected.String(), handleMethod.Type.Out(i).String()),
			})
		}
	}

	return nil
}

func (b *commandBus) Dispatch(ctx context.Context, cmd Command) (Result, error) {
	cmdType := reflect.TypeOf(cmd)

	handlerUntyped, ok := b.handlers.Load(cmdType)
	if !ok {
		return nil, ErrNoCommandHandlerFound{CommandType: cmdType}
	}
	handlerVal := reflect.ValueOf(handlerUntyped)
	handleMethod := handlerVal.MethodByName("Handle")
	if !handleMethod.IsValid() {
		return nil, ErrInvalidCommandHandler{
			HandlerType: reflect.TypeOf(handlerUntyped),
			CommandType: cmdType,
		}
	}

	args := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(cmd),
	}

	results := handleMethod.Call(args)

	var err error
	if len(results) > 1 && !results[1].IsNil() {
		err = results[1].Interface().(error)
	}

	var result any
	if len(results) > 0 {
		result = results[0].Interface()
	}

	return result, err
}
