package grpc

import (
	"github.com/rs/zerolog"
	"github.com/udholdenhed/unotes/auth/internal/service"
	"google.golang.org/grpc"

	pb "github.com/udholdenhed/unotes/auth/api/proto"
)

type Handler struct {
	server *oAuth2ServiceServer
}

func NewHandler(services *service.Service, logger *zerolog.Logger) *Handler {
	return &Handler{
		server: &oAuth2ServiceServer{services: services, logger: logger},
	}
}

func (h *Handler) Init(addr string) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterOAuth2ServiceServer(s, h.server)

	return s
}
