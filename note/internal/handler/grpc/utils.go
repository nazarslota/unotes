package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

type Logger interface {
	InfoFields(msg string, fields map[string]any)
}

func newUnaryLoggerInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		fields := map[string]any{
			"method":   info.FullMethod,
			"duration": duration.String(),
		}
		logger.InfoFields("gRPC, unary request handled.", fields)

		return resp, err
	}
}

func newStreamLoggerInterceptor(logger Logger) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		err := handler(srv, stream)
		duration := time.Since(start)

		fields := map[string]any{
			"method":   info.FullMethod,
			"duration": duration.String(),
		}
		logger.InfoFields("gRPC, stream request handled.", fields)

		return err
	}
}
