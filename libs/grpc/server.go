package grpc

import (
	"context"
	"log/slog"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/pratchaya-maneechot/service-exchange/libs/infra/observability"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server          *grpc.Server
	logger          *slog.Logger
	healthServer    *health.Server
	address         string
	shutdownTimeout time.Duration
}

type ConfigGRPCServer struct {
	Address           string
	EnableHealthCheck bool
	EnableReflection  bool
	ShutdownTimeout   time.Duration
	MetricsRecorder   *observability.MetricsRecorder
	Options           []grpc.ServerOption
}

func NewServer(cfg ConfigGRPCServer, logger *slog.Logger) (*GRPCServer, error) {
	interceptors := []grpc.UnaryServerInterceptor{
		UnaryRecoveryInterceptor(logger),
		UnaryTraceInterceptor(),
		UnaryLoggerInterceptor(logger),
	}

	if cfg.MetricsRecorder != nil {
		interceptors = append(interceptors, UnaryMetricsInterceptor(*cfg.MetricsRecorder))
	}

	opts := append(cfg.Options, grpc.ChainUnaryInterceptor(interceptors...), grpc.StatsHandler(otelgrpc.NewServerHandler()))
	grpcServer := grpc.NewServer(opts...)

	var healthServer *health.Server
	if cfg.EnableHealthCheck {
		healthServer = health.NewServer()
		grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
		healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
		healthServer.SetServingStatus("user.UserService", grpc_health_v1.HealthCheckResponse_SERVING)
		logger.Info("health check service enabled")
	}

	if cfg.EnableReflection {
		reflection.Register(grpcServer)
		logger.Info("gRPC reflection enabled")
	}

	return &GRPCServer{
		server:          grpcServer,
		logger:          logger,
		healthServer:    healthServer,
		shutdownTimeout: cfg.ShutdownTimeout,
		address:         cfg.Address,
	}, nil
}

func (s *GRPCServer) RegisHandler(fn func(*grpc.Server)) error {
	fn(s.server)
	return nil
}

func (s *GRPCServer) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return errors.Wrapf(err, "failed to listen on %s", s.address)
	}

	go s.handleShutdown(ctx)

	s.logger.Info("gRPC server starting", "address", s.address)
	if err := s.server.Serve(listener); err != nil {
		s.logger.Error("gRPC server serve error", "error", err)
		return errors.Wrap(err, "failed to serve gRPC server")
	}
	return nil
}

func (s *GRPCServer) handleShutdown(ctx context.Context) {
	<-ctx.Done()
	s.logger.Info("initiating graceful shutdown...")

	if s.healthServer != nil {
		s.healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		s.healthServer.SetServingStatus("user.UserService", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("gRPC server stopped gracefully")
	case <-shutdownCtx.Done():
		s.logger.Warn("shutdown timeout exceeded, forcing stop")
		s.server.Stop()
	}
}
