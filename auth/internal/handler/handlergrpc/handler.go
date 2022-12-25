package handlergrpc

import (
	"github.com/rs/zerolog"
	"github.com/udholdenhed/unotes/auth/internal/service"
	"google.golang.org/grpc"

	pb "github.com/udholdenhed/unotes/auth/api/proto"
)

type Handler struct {
	services *service.Service
	logger   *zerolog.Logger

	server *oAuth2ServerGRPC
}

func NewHandler(services *service.Service, logger *zerolog.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
		server:   &oAuth2ServerGRPC{},
	}
}

func (h *Handler) Init(addr string) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterOAuth2ServiceServer(s, h.server)

	return s
}

// protoc -Iapi/proto --go_out=. --go_opt=module=github.com/udholdenhed/unotes/auth --go-grpc_out=. --go-grpc_opt=module=github.com/udholdenhed/unotes/auth api/proto/oauth2.proto
