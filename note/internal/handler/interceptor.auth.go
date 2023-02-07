package handler

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type JWTVerifier interface {
	Verify(token string) (jwt.MapClaims, error)
}

type authInterceptor struct {
	JWTVerifier JWTVerifier
}

func newAuthInterceptor(jwtVerifier JWTVerifier) *authInterceptor {
	return &authInterceptor{JWTVerifier: jwtVerifier}
}

func (i *authInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := i.authorize(ctx); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func (i *authInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, ss)
	}
}

func (i *authInterceptor) authorize(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	values, ok := md["authorization"]
	if !ok {
		return status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	if len(values) == 0 {
		return status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	token := values[0]
	if strings.Contains(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}

	_, err := i.JWTVerifier.Verify(token)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}
	return nil
}
