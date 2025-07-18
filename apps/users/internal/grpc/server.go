package grpc

import (
	"log/slog"
	"time"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/grpc/handlers"
	"github.com/pratchaya-maneechot/service-exchange/libs/bus"
	lg "github.com/pratchaya-maneechot/service-exchange/libs/grpc"
	"github.com/pratchaya-maneechot/service-exchange/libs/infra/observability"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func NewGRPCServer(
	cfg *config.Config,
	bus bus.Bus,
	logger *slog.Logger,
	metricsRecorder observability.MetricsRecorder,
) (*lg.GRPCServer, error) {

	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     cfg.Server.MaxConnectionIdle,
			MaxConnectionAge:      cfg.Server.MaxConnectionAge,
			MaxConnectionAgeGrace: cfg.Server.MaxConnectionAgeGrace,
			Time:                  cfg.Server.KeepaliveTime,
			Timeout:               cfg.Server.KeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.MaxRecvMsgSize(cfg.Server.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(cfg.Server.MaxSendMsgSize),
		grpc.MaxConcurrentStreams(cfg.Server.MaxConcurrentStreams),
	}

	server, err := lg.NewServer(lg.ConfigGRPCServer{
		Address:           cfg.Server.Address,
		EnableHealthCheck: cfg.Server.EnableReflection,
		EnableReflection:  cfg.Server.EnableReflection,
		ShutdownTimeout:   cfg.Server.ShutdownTimeout,
		Options:           opts,
	}, logger)
	if err != nil {
		return nil, err
	}

	server.RegisHandler(func(gs *grpc.Server) {
		handlers.RegisUserGRPCHandler(gs, bus.CommandBus, bus.QueryBus, logger)
	})

	return server, nil
}
