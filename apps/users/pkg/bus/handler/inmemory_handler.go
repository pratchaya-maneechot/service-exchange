package handler

import "sync"

type InMemoryBusHandler struct {
	cache sync.Map
}

func NewInMemoryBusHandler() BusHandler {
	return &InMemoryBusHandler{}
}

func (i *InMemoryBusHandler) Load(key any) (value any, ok bool) {
	return i.cache.Load(key)
}

func (i *InMemoryBusHandler) Store(key any, value any) (duplicated bool) {
	_, loaded := i.cache.LoadOrStore(key, value)
	return loaded
}

func (i *InMemoryBusHandler) Delete(key any) {
	i.cache.Delete(key)
}
