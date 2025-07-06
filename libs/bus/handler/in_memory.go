package handler

import (
	"sync"

	"github.com/pratchaya-maneechot/service-exchange/libs/bus/command"
	"github.com/pratchaya-maneechot/service-exchange/libs/bus/query"
)

type inMemoryBusHandler struct {
	cache sync.Map
}

func NewInMemoryQueryBusHandler() query.QueryBusHandler {
	return &inMemoryBusHandler{}
}

func NewInMemoryCommandBusHandler() command.CommandBusHandler {
	return &inMemoryBusHandler{}
}

func (i *inMemoryBusHandler) Load(key any) (value any, ok bool) {
	return i.cache.Load(key)
}

func (i *inMemoryBusHandler) Store(key any, value any) (duplicated bool) {
	_, loaded := i.cache.LoadOrStore(key, value)
	return loaded
}

func (i *inMemoryBusHandler) Delete(key any) {
	i.cache.Delete(key)
}
