package handler

import (
	"context"
	"strings"
	"time"

	"github.com/nazarslota/unotes/auth/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Logger Interceptor

type GRPCLogger interface {
	InfoFields(msg string, fields map[string]any)
}

type loggerInterceptor struct {
	Logger GRPCLogger
}

type loggerInterceptorOptions struct {
	Logger GRPCLogger
}

func newLoggerInterceptor(options loggerInterceptorOptions) *loggerInterceptor {
	return &loggerInterceptor{Logger: options.Logger}
}

func (i *loggerInterceptor) Unary() grpc.UnaryServerInterceptor {
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

func (i *loggerInterceptor) Stream() grpc.StreamServerInterceptor {
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

// Auth Interceptor

type accessTokenValidator interface {
	Validate(token string) (jwt.AccessTokenClaims, error)
}

type authInterceptor struct {
	AccessTokenValidator accessTokenValidator
	Methods              []string
}

type authInterceptorOptions struct {
	AccessTokenValidator accessTokenValidator
	Methods              []string
}

func newAuthInterceptor(options authInterceptorOptions) *authInterceptor {
	return &authInterceptor{
		AccessTokenValidator: options.AccessTokenValidator,
		Methods:              options.Methods,
	}
}

func (i *authInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if !i.allowed(info.FullMethod) {
			return handler(ctx, req)
		}

		ctx, err := i.authorize(ctx)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func (i *authInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if !i.allowed(info.FullMethod) {
			return handler(srv, stream)
		}

		ctx, err := i.authorize(stream.Context())
		if err != nil {
			return err
		}
		return handler(srv, newStreamServerWrapper(stream, ctx))
	}
}

func (i *authInterceptor) allowed(method string) bool { return allowed(method, i.Methods) }

func (i *authInterceptor) authorize(ctx context.Context) (context.Context, error) {
	tokens := metadata.ValueFromIncomingContext(ctx, "authorization")
	if len(tokens) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization header is not provided")
	}

	splits := strings.SplitN(tokens[0], " ", 2)
	if len(splits) < 2 {
		return nil, status.Error(codes.Unauthenticated, "authentication token is not provided")
	}
	token := splits[1]

	claims, err := i.AccessTokenValidator.Validate(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
	}
	return context.WithValue(ctx, "claims", claims), nil
}

// Helpers

type streamServerWrapper struct {
	grpc.ServerStream
	ctx context.Context
}

func newStreamServerWrapper(stream grpc.ServerStream, ctx context.Context) *streamServerWrapper {
	return &streamServerWrapper{ServerStream: stream, ctx: ctx}
}

func (s streamServerWrapper) Context() context.Context { return s.ctx }

func allowed(method string, allowed []string) bool {
	for _, v := range allowed {
		if v == "*" || v == method {
			return true
		}
		if len(v) > 0 && v[len(v)-1] == '*' && strings.HasPrefix(method, v[:len(v)-1]) {
			return true
		}
	}
	return false
}
