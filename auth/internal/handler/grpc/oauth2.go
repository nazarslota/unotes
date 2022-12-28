package grpc

import (
	"context"

	"github.com/rs/zerolog"
	pb "github.com/udholdenhed/unotes/auth/api/proto"
	"github.com/udholdenhed/unotes/auth/internal/service"
)

type oAuth2ServiceServer struct {
	services *service.Service
	logger   *zerolog.Logger

	pb.OAuth2ServiceServer
}

func (h *oAuth2ServiceServer) SignUp(ctx context.Context, in *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	return &pb.SignUpResponse{}, nil
}

func (h *oAuth2ServiceServer) SignIn(ctx context.Context, request *pb.SignInRequest) (*pb.SignInResponse, error) {
	return &pb.SignInResponse{}, nil
}

func (h *oAuth2ServiceServer) SignOut(ctx context.Context, request *pb.SignOutRequest) (*pb.SignOutResponse, error) {
	return &pb.SignOutResponse{}, nil
}

func (h *oAuth2ServiceServer) Refresh(ctx context.Context, request *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	return &pb.RefreshResponse{}, nil
}
