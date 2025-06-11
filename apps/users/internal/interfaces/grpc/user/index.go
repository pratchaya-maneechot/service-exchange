package server

import (
	"log/slog"

	userpb "github.com/pratchaya-maneechot/service-exchange/apps/users/api/proto/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/common"
	"google.golang.org/grpc"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
	commandBus common.CommandBus
	queryBus   common.QueryBus
	logger     *slog.Logger
}

func Register(gRPCServer *grpc.Server, appModule *app.AppModule, logger *slog.Logger) {
	server := &UserServer{
		commandBus: appModule.CommandBus,
		queryBus:   appModule.QueryBus,
		logger:     logger,
	}
	userpb.RegisterUserServiceServer(gRPCServer, server)
}
