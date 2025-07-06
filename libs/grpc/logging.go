package grpc

import (
	"context"
	"log/slog"
	"time"

	"github.com/pratchaya-maneechot/service-exchange/libs/infra/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryTraceInterceptor creates a new span for each gRPC request and propagates context.
func UnaryTraceInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		tracer := otel.Tracer("grpc-server")
		ctx, span := tracer.Start(ctx, info.FullMethod, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		resp, err := handler(ctx, req)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		} else {
			span.SetStatus(codes.Ok, "success")
		}
		return resp, err
	}
}

// UnaryMetricsInterceptor records gRPC request metrics.
func UnaryMetricsInterceptor(metricsRecorder observability.MetricsRecorder) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start).Seconds()

		statusCode := status.Code(err).String() // Get gRPC status code from error
		metricsRecorder.RecordGrpcRequestTotal(info.FullMethod, statusCode)
		metricsRecorder.RecordGrpcRequestDuration(info.FullMethod, duration)

		return resp, err
	}
}

// UnaryRecoveryInterceptor recovers from panics and logs them.
func UnaryRecoveryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Panic in gRPC handler", "panic", r, "method", info.FullMethod)
				err = status.Errorf(grpcCodes.Internal, "panic: %v", r)
			}
		}()
		return handler(ctx, req)
	}
}
