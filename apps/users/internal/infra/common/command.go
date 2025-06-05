package common

import (
	"context"
	"fmt"
	"sync"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/common"
)

var ErrHandlerNotFound = fmt.Errorf("handler not found")

type CommandBusImpl struct {
	handlers map[string]common.CommandHandler
	mutex    sync.RWMutex
}

func NewCommandBus() common.CommandBus {
	return &CommandBusImpl{
		handlers: make(map[string]common.CommandHandler),
	}
}

func (b *CommandBusImpl) Register(commandName string, handler common.CommandHandler) {
	if commandName == "" {
		panic("command name cannot be empty")
	}
	if handler == nil {
		panic("handler cannot be nil")
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.handlers[commandName] = handler
}

func (b *CommandBusImpl) Send(ctx context.Context, command common.Command) (interface{}, error) {
	if command == nil {
		return nil, fmt.Errorf("command cannot be nil")
	}

	commandName := command.CommandName()
	if commandName == "" {
		return nil, fmt.Errorf("command name cannot be empty")
	}

	b.mutex.RLock()
	handler, exists := b.handlers[commandName]
	b.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrHandlerNotFound, commandName)
	}

	return handler.Handle(ctx, command)
}
