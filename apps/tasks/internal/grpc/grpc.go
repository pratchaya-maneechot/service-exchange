package grpc

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/grpc/middleware"
	taskServer "github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/grpc/task"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/pkg/bus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	gRPCServer   *grpc.Server
	listener     net.Listener
	config       *config.Config
	logger       *slog.Logger
	healthServer *health.Server
}

func NewServer(
	cfg *config.Config,
	bus *bus.Bus,
	logger *slog.Logger,
) (*Server, error) {

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

	if cfg.Security.EnableTLS {
		creds, err := loadTLSCredentials(cfg.Security)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load TLS credentials")
		}
		opts = append(opts, grpc.Creds(creds))
	}

	opts = append(opts,
		grpc.ChainUnaryInterceptor(
			middleware.UnaryLoggingInterceptor(logger),
			middleware.UnaryMetricsInterceptor(),
			middleware.UnaryRecoveryInterceptor(logger),
			middleware.UnaryRateLimitInterceptor(cfg.Security.RateLimitRPS, cfg.Security.RateLimitBurst),
		),
		grpc.ChainStreamInterceptor(
			middleware.StreamLoggingInterceptor(logger),
			middleware.StreamMetricsInterceptor(),
			middleware.StreamRecoveryInterceptor(logger),
		),
	)

	gRPCServer := grpc.NewServer(opts...)

	taskServer.Register(gRPCServer, bus, logger)

	var healthServer *health.Server
	if cfg.Server.EnableHealthCheck {
		healthServer = health.NewServer()
		grpc_health_v1.RegisterHealthServer(gRPCServer, healthServer)
		healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
		healthServer.SetServingStatus("user.UserService", grpc_health_v1.HealthCheckResponse_SERVING)
		logger.Info("health check service enabled")
	}

	if cfg.Server.EnableReflection {
		reflection.Register(gRPCServer)
		logger.Info("gRPC reflection enabled")
	}

	return &Server{
		gRPCServer:   gRPCServer,
		config:       cfg,
		logger:       logger,
		healthServer: healthServer,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.config.Server.Address)
	if err != nil {
		return errors.Wrapf(err, "failed to listen on %s", s.config.Server.Address)
	}
	s.listener = listener

	s.logger.Info("gRPC server starting",
		"address", s.config.Server.Address,
		"tls_enabled", s.config.Security.EnableTLS)

	go s.handleShutdown(ctx)

	if err := s.gRPCServer.Serve(listener); err != nil {
		s.logger.Error("gRPC server serve error", "error", err)
		return errors.Wrap(err, "failed to serve gRPC server")
	}

	return nil
}

func (s *Server) handleShutdown(ctx context.Context) {
	<-ctx.Done()
	s.logger.Info("initiating gRPC graceful shutdown...")

	if s.healthServer != nil {
		s.healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		s.healthServer.SetServingStatus("user.UserService", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.Server.ShutdownTimeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		s.gRPCServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("gRPC server stopped gracefully")
	case <-shutdownCtx.Done():
		s.logger.Warn("shutdown timeout exceeded, forcing stop")
		s.gRPCServer.Stop()
	}
}

func loadTLSCredentials(cfg config.SecurityConfig) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(cfg.TLSCertFile, cfg.TLSKeyFile)
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.NoClientCert,
	}), nil
}
