package grpc

import (
	pb "github.com/nazarslota/unotes/auth/api/proto"
	"github.com/nazarslota/unotes/auth/internal/service"
	"google.golang.org/grpc"
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
		grpc.UnaryInterceptor(newUnaryLoggerInterceptor(h.logger)),
		grpc.StreamInterceptor(newStreamLoggerInterceptor(h.logger)),
	)
	pb.RegisterOAuth2ServiceServer(s, h.oAuth2)

	return &Server{address: h.address, server: s}
}

// Logger.
