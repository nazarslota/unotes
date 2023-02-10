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
	grpcAddr   string
	restAddr   string
	grpcLogger GRPCLogger
	restLogger RESTLogger

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
	logger := newLoggerInterceptor(loggerInterceptorOptions{
		Logger: h.grpcLogger,
	})

	auth := newAuthInterceptor(authInterceptorOptions{
		AccessTokenValidator: h.services.JWTService.AccessTokenValidator,
		Methods:              []string{"/NoteService/*"},
	})

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(logger.Unary(), auth.Unary()),
		grpc.ChainStreamInterceptor(logger.Stream(), auth.Stream()),
	)
	pb.RegisterNoteServiceServer(server, &h.noteServiceServer)
	reflection.Register(server)

	return server
}

func (h *Handler) restServer() *http.Server {
	mux := runtime.NewServeMux()
	_ = pb.RegisterNoteServiceHandlerFromEndpoint(context.Background(), mux, h.grpcAddr, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})

	loggerMiddleware := newLoggerMiddleware(loggerMiddlewareOptions{
		Logger: h.restLogger,
	})
	return &http.Server{Handler: loggerMiddleware.Middleware(mux)}
}
