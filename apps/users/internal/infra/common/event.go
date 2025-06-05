package common

import (
	"fmt"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/common"
)

type EventBusImpl struct {
}

func NewEventBus() common.EventBus {
	return &EventBusImpl{}
}

func (b *EventBusImpl) Publish(handler common.Event) error {
	fmt.Print(handler.EventName())
	fmt.Print(handler)
	return nil
}
