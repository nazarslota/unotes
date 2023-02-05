package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	Addr   string
	Server *grpc.Server
}

func newServer(addr string, srv *grpc.Server) *server {
	return &server{Addr: addr, Server: srv}
}

func (s server) Serve() error {
	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	if err := s.Server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

func (s server) Shutdown(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context done: %w", ctx.Err())
	default:
	}
	s.Server.GracefulStop()
	return nil
}
