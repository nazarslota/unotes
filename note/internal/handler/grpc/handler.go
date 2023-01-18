package grpc

import (
	pb "github.com/nazarslota/unotes/note/api/proto"
	"github.com/nazarslota/unotes/note/internal/service"
	"google.golang.org/grpc"
)

type Handler struct {
	address string
	logger  Logger

	service service.Services
	note    *noteServiceServer
}

func NewHandler(options ...HandlerOption) *Handler {
	h := &Handler{}
	for _, option := range options {
		option(h)
	}
	return h
}

func (h *Handler) S() *Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(newUnaryLoggerInterceptor(h.logger)),
		grpc.StreamInterceptor(newStreamLoggerInterceptor(h.logger)),
	)
	pb.RegisterNoteServiceServer(s, h.note)

	return &Server{server: s}
}
