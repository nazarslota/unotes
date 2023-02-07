package handler

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

type Logger interface {
	InfoFields(msg string, fields map[string]any)
}

type grpcLoggerInterceptor struct {
	Logger Logger
}

func newGRPCLoggerInterceptor(logger Logger) *grpcLoggerInterceptor {
	return &grpcLoggerInterceptor{Logger: logger}
}

func (i *grpcLoggerInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		fields := map[string]any{
			"method":   info.FullMethod,
			"duration": duration.String(),
		}
		i.Logger.InfoFields("gRPC, unary request handled.", fields)

		return resp, err
	}
}

func (i *grpcLoggerInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		err := handler(srv, stream)
		duration := time.Since(start)

		fields := map[string]any{
			"method":   info.FullMethod,
			"duration": duration.String(),
		}
		i.Logger.InfoFields("gRPC, stream request handled.", fields)

		return err
	}
}
