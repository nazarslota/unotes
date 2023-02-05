package rest

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
)

type server struct {
	addr   string
	server *http.Server
}

func newServer(addr string, srv *http.Server) *server {
	return &server{addr: addr, server: srv}
}

func (s server) Serve() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	if err := s.server.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

func (s server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
