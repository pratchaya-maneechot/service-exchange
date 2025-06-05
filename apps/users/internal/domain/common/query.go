package common

import "context"

type Query interface {
	QueryName() string
}

type QueryHandler interface {
	Handle(ctx context.Context, query Query) (interface{}, error)
}

type QueryBus interface {
	Register(commandName string, handler QueryHandler)
	Execute(ctx context.Context, query Query) (interface{}, error)
}
