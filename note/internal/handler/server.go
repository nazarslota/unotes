package handler

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

type server struct {
	GRPCAddr   string
	RESTAddr   string
	GRPCServer *grpc.Server
	RESTServer *http.Server
}

func newServer(grpcAddr, restAddr string, grpcServer *grpc.Server, restServer *http.Server) *server {
	return &server{
		GRPCAddr:   grpcAddr,
		RESTAddr:   restAddr,
		GRPCServer: grpcServer,
		RESTServer: restServer,
	}
}

func (s server) ServeGRPC() error {
	lis, err := net.Listen("tcp", s.GRPCAddr)
	if err != nil {
		return fmt.Errorf("failed to create a listener: %w", err)
	}

	if err := s.GRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("internal error: %w", err)
	}
	return nil
}

func (s server) ServeREST() error {
	lis, err := net.Listen("tcp", s.RESTAddr)
	if err != nil {
		return fmt.Errorf("failed to create a listener: %w", err)
	}

	if err := s.RESTServer.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("internal error: %w", err)
	}
	return nil
}

func (s server) ShutdownGRPC(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context done: %w", ctx.Err())
	default:
	}

	s.GRPCServer.GracefulStop()
	return nil
}

func (s server) ShutdownREST(ctx context.Context) error {
	if err := s.RESTServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	return nil
}
