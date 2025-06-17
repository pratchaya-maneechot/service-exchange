package middleware

import (
	"context"
	"strconv"
	"time"

	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/infra/observability/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryMetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		code := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			}
		}

		metrics.GrpcRequestsTotal.WithLabelValues(info.FullMethod, strconv.Itoa(int(code))).Inc()
		metrics.GrpcRequestDuration.WithLabelValues(info.FullMethod).Observe(duration.Seconds())

		return resp, err
	}
}

func StreamMetricsInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		err := handler(srv, stream)

		duration := time.Since(start)
		code := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			} else {
				code = codes.Unknown
			}
		}

		metrics.GrpcRequestsTotal.WithLabelValues(info.FullMethod, strconv.Itoa(int(code))).Inc()
		metrics.GrpcRequestDuration.WithLabelValues(info.FullMethod).Observe(duration.Seconds())

		return err
	}
}
