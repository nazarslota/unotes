package handler

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

type Server struct {
	grpcAddr   string
	restAddr   string
	grpcServer *grpc.Server
	restServer *http.Server
}

func NewServer(grpcAddr, restAddr string, grpcServer *grpc.Server, restServer *http.Server) *Server {
	return &Server{
		grpcAddr:   grpcAddr,
		restAddr:   restAddr,
		grpcServer: grpcServer,
		restServer: restServer,
	}
}

func (s *Server) ServeGRPC() error {
	lis, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to create a listener: %w", err)
	}

	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("internal error: %w", err)
	}
	return nil
}

func (s *Server) ServeREST() error {
	lis, err := net.Listen("tcp", s.restAddr)
	if err != nil {
		return fmt.Errorf("failed to create a listener: %w", err)
	}

	if err := s.restServer.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("internal error: %w", err)
	}
	return nil
}

func (s *Server) ShutdownGRPC(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context done: %w", ctx.Err())
	default:
	}

	s.grpcServer.GracefulStop()
	return nil
}

func (s *Server) ShutdownREST(ctx context.Context) error {
	if err := s.restServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	return nil
}
