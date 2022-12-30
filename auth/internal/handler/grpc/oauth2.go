package grpc

import (
	"context"

	"github.com/rs/zerolog"
	pb "github.com/udholdenhed/unotes/auth/api/proto"
	"github.com/udholdenhed/unotes/auth/internal/service"
	"github.com/udholdenhed/unotes/auth/internal/service/oauth2"
)

type oAuth2ServiceServer struct {
	services *service.Service
	logger   *zerolog.Logger

	pb.OAuth2ServiceServer
}

func (h *oAuth2ServiceServer) SignUp(ctx context.Context, in *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	request := &oauth2.SignUpRequest{
		Username: in.Username,
		Password: in.Password,
	}

	_, err := h.services.OAuth2Service.SignUpRequestHandler.Handler(ctx, request)
	if err != nil {
		return nil, err
	}

	return &pb.SignUpResponse{}, nil
}

func (h *oAuth2ServiceServer) SignIn(ctx context.Context, in *pb.SignInRequest) (*pb.SignInResponse, error) {
	request := &oauth2.SignInRequest{
		Username: in.Username,
		Password: in.Password,
	}

	response, err := h.services.OAuth2Service.SingInRequestHandler.Handle(ctx, request)
	if err != nil {
		return nil, err
	}

	return &pb.SignInResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}, nil
}

func (h *oAuth2ServiceServer) SignOut(ctx context.Context, in *pb.SignOutRequest) (*pb.SignOutResponse, error) {
	request := &oauth2.SignOutRequest{
		AccessToken: in.AccessToken,
	}

	_, err := h.services.OAuth2Service.SignOutRequestHandler.Handle(ctx, request)
	if err != nil {
		return nil, err
	}

	return &pb.SignOutResponse{}, nil
}

func (h *oAuth2ServiceServer) Refresh(ctx context.Context, in *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	request := &oauth2.RefreshRequest{
		RefreshToken: in.RefreshToken,
	}

	response, err := h.services.OAuth2Service.RefreshRequestHandler.Handle(ctx, request)
	if err != nil {
		return nil, err
	}

	return &pb.RefreshResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}, nil
}
