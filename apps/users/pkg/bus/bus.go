package bus

import "github.com/google/wire"

type Bus struct {
	CommandBus CommandBus
	QueryBus   QueryBus
}

func NewBus(
	cb CommandBus,
	qb QueryBus,
) *Bus {
	return &Bus{
		CommandBus: cb,
		QueryBus:   qb,
	}
}

var BusModuleSet = wire.NewSet(
	NewQueryBus,
	NewCommandBus,
	NewBus,
)
