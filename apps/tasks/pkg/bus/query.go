package bus

import (
	"context"
	"fmt"
	"reflect"
)

// IQuery is an interface that all queries must implement.
// It's a marker interface.
type IQuery interface{}

// IQueryHandler is an interface for handling a specific query and returning a result.
type IQueryHandler[Q IQuery, R any] interface {
	Handle(ctx context.Context, query Q) (R, error)
}

// IQueryBus is the interface for dispatching queries.
type IQueryBus interface {
	Dispatch(ctx context.Context, query IQuery) (interface{}, error)
	RegisterHandler(queryType IQuery, handler interface{}) error
}

// ErrNoQueryHandlerFound is returned when no handler is registered for a query.
type ErrNoQueryHandlerFound struct {
	QueryType reflect.Type
}

func (e ErrNoQueryHandlerFound) Error() string {
	return fmt.Sprintf("no query handler found for query type: %s", e.QueryType.String())
}

// ErrInvalidQueryHandler is returned when a registered handler does not implement IQueryHandler correctly.
type ErrInvalidQueryHandler struct {
	HandlerType reflect.Type
	QueryType   reflect.Type
	Reason      string // Added reason for more detail
}

func (e ErrInvalidQueryHandler) Error() string {
	return fmt.Sprintf("invalid query handler type %s for query type %s: %s", e.HandlerType.String(), e.QueryType.String(), e.Reason)
}

// ErrQueryAlreadyRegistered is returned when a query handler is registered multiple times.
type ErrQueryAlreadyRegistered struct {
	QueryType reflect.Type
}

func (e ErrQueryAlreadyRegistered) Error() string {
	return fmt.Sprintf("query handler already registered for query type: %s", e.QueryType.String())
}
