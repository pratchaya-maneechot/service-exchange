package bus

import (
	"github.com/google/wire"
	"github.com/pratchaya-maneechot/service-exchange/libs/bus/command"
	"github.com/pratchaya-maneechot/service-exchange/libs/bus/handler"
	"github.com/pratchaya-maneechot/service-exchange/libs/bus/query"
)

type Bus struct {
	CommandBus command.CommandBus
	QueryBus   query.QueryBus
}

var BusModuleSet = wire.NewSet(
	handler.NewInMemoryQueryBusHandler,
	query.NewQueryBus,
	handler.NewInMemoryCommandBusHandler,
	command.NewCommandBus,
	wire.Struct(new(Bus), "*"),
)
