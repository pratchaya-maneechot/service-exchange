package bus

import "github.com/google/wire"

type Bus struct {
	CommandBus ICommandBus
	QueryBus   IQueryBus
}

func NewBus(
	cb ICommandBus,
	qb IQueryBus,
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
