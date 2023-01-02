package rest

import (
	"context"
	"net/http"
)

type Server struct {
	server *http.Server
}

func (s *Server) Serve() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
