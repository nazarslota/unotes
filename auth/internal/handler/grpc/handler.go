package grpc

import (
	"context"

	pb "github.com/nazarslota/unotes/auth/api/proto"
	"github.com/nazarslota/unotes/auth/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Handler struct {
	addr   string
	logger Logger

	services            service.Services
	oAuth2ServiceServer oAuth2ServiceServer
}

func NewHandler(options ...HandlerOption) *Handler {
	h := new(Handler)
	for _, option := range options {
		option(h)
	}
	h.oAuth2ServiceServer = newOAuth2ServiceServer(h.services)
	return h
}

type Server interface {
	Serve() error
	Shutdown(ctx context.Context) error
}

func (h *Handler) Server() Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(newGRPCLoggerUnaryInterceptor(h.logger)),
		grpc.StreamInterceptor(newGRPCLoggerStreamInterceptor(h.logger)),
	)
	pb.RegisterOAuth2ServiceServer(server, &h.oAuth2ServiceServer)
	reflection.Register(server)

	return newServer(h.addr, server)
}
