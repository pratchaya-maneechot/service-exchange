package query

import (
	"context"
	"fmt"
	"reflect"
)

type Query any
type Result any

type QueryBusHandler interface {
	Load(key any) (value any, ok bool)
	Store(key any, value any) (duplicated bool)
	Delete(key any)
}

type QueryHandler[Q Query, R any] interface {
	Handle(ctx context.Context, query Q) (R, error)
}

type QueryBus interface {
	Dispatch(ctx context.Context, query Query) (Result, error)
	RegisterHandler(queryType Query, handler any) error
}

type ErrNoQueryHandlerFound struct {
	QueryType reflect.Type
}

func (e ErrNoQueryHandlerFound) Error() string {
	return fmt.Sprintf("no query handler found for query type: %s", e.QueryType.String())
}

type ErrInvalidQueryHandler struct {
	HandlerType reflect.Type
	QueryType   reflect.Type
	Reason      string
}

func (e ErrInvalidQueryHandler) Error() string {
	return fmt.Sprintf("invalid query handler type %s for query type %s: %s", e.HandlerType.String(), e.QueryType.String(), e.Reason)
}

type ErrQueryAlreadyRegistered struct {
	QueryType reflect.Type
}

func (e ErrQueryAlreadyRegistered) Error() string {
	return fmt.Sprintf("query handler already registered for query type: %s", e.QueryType.String())
}
