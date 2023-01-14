package grpc

import (
	"context"

	pb "github.com/nazarslota/unotes/auth/api/proto"
	"github.com/nazarslota/unotes/auth/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type Handler struct {
	address  string
	logger   Logger
	services *service.Services

	oAuth2 *oAuth2ServiceServer
}

func NewHandler(options ...HandlerOption) *Handler {
	h := &Handler{}
	for _, option := range options {
		option(h)
	}

	h.oAuth2 = &oAuth2ServiceServer{
		services: h.services,
		logger:   h.logger,
	}
	return h
}

func (h *Handler) S() *Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(newRequestLoggerUnaryInterceptor(h.logger)),
	)
	pb.RegisterOAuth2ServiceServer(s, h.oAuth2)

	return &Server{address: h.address, server: s}
}

// Logger.

type Logger interface {
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

func newRequestLoggerUnaryInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)

		method, _ := grpc.Method(ctx)
		logger.Infof("gRPC: \"%s %s %s\" %s", "UNARY", method, "HTTP/2.0", status.Code(err).String())

		return resp, err
	}
}
