package handler

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/nazarslota/unotes/note/api/proto"
	"github.com/nazarslota/unotes/note/internal/service"
	"google.golang.org/grpc"
)

type Handler struct {
	grpcAddr string
	restAddr string
	logger   Logger

	services          service.Services
	noteServiceServer noteServiceServer
}

func NewHandler(options ...Option) *Handler {
	h := &Handler{}
	for _, option := range options {
		option(h)
	}

	h.noteServiceServer = noteServiceServer{services: h.services}
	return h
}

func (h *Handler) S() *Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(newGRPCLoggerUnaryInterceptor(h.logger)),
		grpc.StreamInterceptor(newGRPCLoggerStreamInterceptor(h.logger)),
	)
	pb.RegisterNoteServiceServer(grpcServer, &h.noteServiceServer)

	mux := runtime.NewServeMux()
	_ = pb.RegisterNoteServiceHandlerServer(context.Background(), mux, h.noteServiceServer)

	restServer := &http.Server{
		Handler: newRESTLogger(mux, h.logger),
	}

	return NewServer(h.grpcAddr, h.restAddr, grpcServer, restServer)
}
