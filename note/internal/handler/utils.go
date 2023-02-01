package handler

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

type Logger interface {
	InfoFields(msg string, fields map[string]any)
}

func newGRPCLoggerUnaryInterceptor(logger Logger) grpc.UnaryServerInterceptor {
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

func newGRPCLoggerStreamInterceptor(logger Logger) grpc.StreamServerInterceptor {
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

func newRESTLogger(handler http.Handler, logger Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		duration := time.Since(start)

		fields := map[string]any{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": duration.String(),
		}
		logger.InfoFields("HTTP request handled.", fields)
	})
}
