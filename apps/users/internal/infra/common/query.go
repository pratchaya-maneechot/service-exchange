package common

import (
	"context"
	"fmt"
	"sync"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/common"
)

type QueryBusImpl struct {
	handlers map[string]common.QueryHandler
	mutex    sync.RWMutex
}

func NewQueryBus() common.QueryBus {
	return &QueryBusImpl{
		handlers: make(map[string]common.QueryHandler),
	}
}

func (b *QueryBusImpl) Register(queryName string, handler common.QueryHandler) {
	if queryName == "" {
		panic("query name cannot be empty")
	}
	if handler == nil {
		panic("handler cannot be nil")
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.handlers[queryName] = handler
}

func (b *QueryBusImpl) Execute(ctx context.Context, query common.Query) (interface{}, error) {
	if query == nil {
		return nil, fmt.Errorf("query cannot be nil")
	}

	queryName := query.QueryName()
	if queryName == "" {
		return nil, fmt.Errorf("query name cannot be empty")
	}

	b.mutex.RLock()
	handler, exists := b.handlers[queryName]
	b.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrHandlerNotFound, queryName)
	}

	return handler.Handle(ctx, query)
}
