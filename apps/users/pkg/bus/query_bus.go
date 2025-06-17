package bus

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/pkg/bus/handler"
)

// QueryBus is a concrete implementation of QueryBus that stores handlers in memory.
type queryBus struct {
	handlers handler.BusHandler
}

// NewQueryBus creates a new QueryBus.
func NewQueryBus() QueryBus {
	return &queryBus{
		handlers: handler.NewInMemoryBusHandler(),
	}
}

// RegisterHandler registers a QueryHandler for a specific Query type.
// It performs strict type checking to ensure the handler correctly implements QueryHandler[Q, R].
func (b *queryBus) RegisterHandler(queryType Query, handler any) error {
	queryReflectType := reflect.TypeOf(queryType)
	handlerReflectType := reflect.TypeOf(handler)

	if handler == nil {
		return fmt.Errorf("query handler for type %s cannot be nil", queryReflectType.String())
	}

	if duplicated := b.handlers.Store(queryReflectType, handler); duplicated {
		return ErrQueryAlreadyRegistered{QueryType: queryReflectType}
	}

	handleMethod, found := handlerReflectType.MethodByName("Handle")
	if !found {
		b.handlers.Delete(queryReflectType)
		return ErrInvalidQueryHandler{
			HandlerType: handlerReflectType,
			QueryType:   queryReflectType,
			Reason:      "handler does not have a 'Handle' method",
		}
	}

	// Method's signature (Input: receiver, ctx, query; Output: R, error)
	expectedIn := []reflect.Type{
		reflect.TypeOf((*context.Context)(nil)).Elem(), // ctx
		queryReflectType, // query
	}
	expectedOut := []reflect.Type{
		nil,                                  // Placeholder for generic result type R
		reflect.TypeOf((*error)(nil)).Elem(), // error
	}

	// Check input parameters
	if handleMethod.Type.NumIn() != len(expectedIn)+1 { // +1 for the receiver
		b.handlers.Delete(queryReflectType)
		return ErrInvalidQueryHandler{
			HandlerType: handlerReflectType,
			QueryType:   queryReflectType,
			Reason:      fmt.Sprintf("Handle method has %d input parameters, expected %d (excluding receiver)", handleMethod.Type.NumIn()-1, len(expectedIn)),
		}
	}
	for i, expected := range expectedIn {
		if handleMethod.Type.In(i+1) != expected { // +1 for the receiver
			b.handlers.Delete(queryReflectType)
			return ErrInvalidQueryHandler{
				HandlerType: handlerReflectType,
				QueryType:   queryReflectType,
				Reason:      fmt.Sprintf("Handle method parameter %d type mismatch: expected %s, got %s", i+1, expected.String(), handleMethod.Type.In(i+1).String()),
			}
		}
	}

	// Check output parameters
	if handleMethod.Type.NumOut() != len(expectedOut) {
		b.handlers.Delete(queryReflectType)
		return ErrInvalidQueryHandler{
			HandlerType: handlerReflectType,
			QueryType:   queryReflectType,
			Reason:      fmt.Sprintf("Handle method has %d output parameters, expected %d", handleMethod.Type.NumOut(), len(expectedOut)),
		}
	}
	// For queries, the first return type (R) is generic, so we can't directly check its type against a concrete `expected` `reflect.Type`.
	// We only ensure it exists and the second is `error`.
	if handleMethod.Type.Out(1) != expectedOut[1] { // Check if the second return value is error
		b.handlers.Delete(queryReflectType)
		return ErrInvalidQueryHandler{
			HandlerType: handlerReflectType,
			QueryType:   queryReflectType,
			Reason:      fmt.Sprintf("Handle method return parameter 2 type mismatch: expected %s, got %s", expectedOut[1].String(), handleMethod.Type.Out(1).String()),
		}
	}

	return nil
}

// Dispatch dispatches a Query to its registered handler and returns the result.
func (b *queryBus) Dispatch(ctx context.Context, query Query) (any, error) {
	queryType := reflect.TypeOf(query)

	handlerUntyped, ok := b.handlers.Load(queryType)
	if !ok {
		return nil, ErrNoQueryHandlerFound{QueryType: queryType}
	}

	handlerVal := reflect.ValueOf(handlerUntyped)
	handleMethod := handlerVal.MethodByName("Handle")
	if !handleMethod.IsValid() {
		return nil, ErrInvalidQueryHandler{
			HandlerType: reflect.TypeOf(handlerUntyped),
			QueryType:   queryType,
		}
	}

	args := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(query),
	}

	results := handleMethod.Call(args)

	var err error
	if len(results) > 1 && !results[1].IsNil() { // Error is the second return value
		err = results[1].Interface().(error)
	}

	var result any
	if len(results) > 0 { // Result is the first return value
		result = results[0].Interface()
	}

	return result, err
}
