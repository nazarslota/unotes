package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

type Logger interface {
	Trace(v ...any)
	Tracef(format string, v ...any)
	TraceFields(msg string, fields map[string]any)

	Debug(v ...any)
	Debugf(format string, v ...any)
	DebugFields(msg string, fields map[string]any)

	Info(v ...any)
	Infof(format string, v ...any)
	InfoFields(msg string, fields map[string]any)

	Warn(v ...any)
	Warnf(format string, v ...any)
	WarnFields(msg string, fields map[string]any)

	Error(v ...any)
	Errorf(format string, v ...any)
	ErrorFields(msg string, fields map[string]any)

	Fatal(v ...any)
	Fatalf(format string, v ...any)
	FatalFields(msg string, fields map[string]any)

	Panic(v ...any)
	Panicf(format string, v ...any)
	PanicFields(msg string, fields map[string]any)
}

func newUnaryLoggerInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		fields := map[string]interface{}{
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

		fields := map[string]interface{}{
			"method":   info.FullMethod,
			"duration": duration.String(),
		}
		logger.InfoFields("gRPC, stream request handled.", fields)

		return err
	}
}
