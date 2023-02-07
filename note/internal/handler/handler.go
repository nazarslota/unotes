package handler

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/nazarslota/unotes/note/api/proto"
	"github.com/nazarslota/unotes/note/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
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
	h.noteServiceServer = newNoteServiceServer(h.services)
	return h
}

type Server interface {
	ServeGRPC() error
	ServeREST() error
	ShutdownGRPC(ctx context.Context) error
	ShutdownREST(ctx context.Context) error
}

func (h *Handler) Server() Server {
	return newServer(h.grpcAddr, h.restAddr, h.grpcServer(), h.restServer())
}

func (h *Handler) grpcServer() *grpc.Server {
	loggerInterceptor := newGRPCLoggerInterceptor(h.logger)
	authInterceptor := newAuthInterceptor(h.services.JWTService.JWTVerifier)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggerInterceptor.Unary(),
			authInterceptor.Unary(),
		),
		grpc.ChainStreamInterceptor(
			loggerInterceptor.Stream(),
			authInterceptor.Stream(),
		),
	)
	pb.RegisterNoteServiceServer(grpcServer, &h.noteServiceServer)
	reflection.Register(grpcServer)

	return grpcServer
}

func (h *Handler) restServer() *http.Server {
	mux := runtime.NewServeMux()
	_ = pb.RegisterNoteServiceHandlerFromEndpoint(context.Background(), mux, h.grpcAddr, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})

	loggerMiddleware := newRESTLoggerMiddleware(h.logger)
	return &http.Server{Handler: loggerMiddleware.Middleware(mux, h.logger)}
}
