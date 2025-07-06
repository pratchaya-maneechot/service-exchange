package query

import (
	"context"
	"fmt"
	"reflect"
)

type queryBus struct {
	handlers QueryBusHandler
}

func NewQueryBus(h QueryBusHandler) QueryBus {
	return &queryBus{
		handlers: h,
	}
}

func (b *queryBus) RegisterHandler(queryType Query, handler any) error {
	queryReflectType := reflect.TypeOf(queryType)
	handlerReflectType := reflect.TypeOf(handler)

	if handler == nil {
		return fmt.Errorf("query handler for type %s cannot be nil", queryReflectType.String())
	}

	if duplicated := b.handlers.Store(queryReflectType, handler); duplicated {
		panic(ErrQueryAlreadyRegistered{QueryType: queryReflectType})
	}

	handleMethod, found := handlerReflectType.MethodByName("Handle")
	if !found {
		b.handlers.Delete(queryReflectType)
		panic(ErrInvalidQueryHandler{
			HandlerType: handlerReflectType,
			QueryType:   queryReflectType,
			Reason:      "handler does not have a 'Handle' method",
		})
	}

	expectedIn := []reflect.Type{
		reflect.TypeOf((*context.Context)(nil)).Elem(),
		queryReflectType,
	}
	expectedOut := []reflect.Type{
		nil,
		reflect.TypeOf((*error)(nil)).Elem(),
	}

	if handleMethod.Type.NumIn() != len(expectedIn)+1 {
		b.handlers.Delete(queryReflectType)
		panic(ErrInvalidQueryHandler{
			HandlerType: handlerReflectType,
			QueryType:   queryReflectType,
			Reason:      fmt.Sprintf("Handle method has %d input parameters, expected %d (excluding receiver)", handleMethod.Type.NumIn()-1, len(expectedIn)),
		})
	}
	for i, expected := range expectedIn {
		if handleMethod.Type.In(i+1) != expected {
			b.handlers.Delete(queryReflectType)
			panic(ErrInvalidQueryHandler{
				HandlerType: handlerReflectType,
				QueryType:   queryReflectType,
				Reason:      fmt.Sprintf("Handle method parameter %d type mismatch: expected %s, got %s", i+1, expected.String(), handleMethod.Type.In(i+1).String()),
			})
		}
	}

	if handleMethod.Type.NumOut() != len(expectedOut) {
		b.handlers.Delete(queryReflectType)
		panic(ErrInvalidQueryHandler{
			HandlerType: handlerReflectType,
			QueryType:   queryReflectType,
			Reason:      fmt.Sprintf("Handle method has %d output parameters, expected %d", handleMethod.Type.NumOut(), len(expectedOut)),
		})
	}

	if handleMethod.Type.Out(1) != expectedOut[1] {
		b.handlers.Delete(queryReflectType)
		panic(ErrInvalidQueryHandler{
			HandlerType: handlerReflectType,
			QueryType:   queryReflectType,
			Reason:      fmt.Sprintf("Handle method return parameter 2 type mismatch: expected %s, got %s", expectedOut[1].String(), handleMethod.Type.Out(1).String()),
		})
	}

	return nil
}

func (b *queryBus) Dispatch(ctx context.Context, query Query) (Result, error) {
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
	if len(results) > 1 && !results[1].IsNil() {
		err = results[1].Interface().(error)
	}

	var result any
	if len(results) > 0 {
		result = results[0].Interface()
	}

	return result, err
}
