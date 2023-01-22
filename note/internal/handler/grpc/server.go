package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	address string
	server  *grpc.Server
}

func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to create a listener: %w", err)
	}
	return s.server.Serve(lis)
}

func (s *Server) Shutdown(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to shut down the server: %w", ctx.Err())
	default:
	}

	s.server.GracefulStop()
	return nil
}
