package middleware

import (
	"context"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryRateLimitInterceptor(rps, burst int) grpc.UnaryServerInterceptor {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := limiter.Wait(ctx); err != nil {
			return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded: %v", err)
		}

		return handler(ctx, req)
	}
}
