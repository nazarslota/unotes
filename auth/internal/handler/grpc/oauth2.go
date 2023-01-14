package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/nazarslota/unotes/auth/api/proto"
	"github.com/nazarslota/unotes/auth/internal/service"
	"github.com/nazarslota/unotes/auth/internal/service/oauth2"
)

type oAuth2ServiceServer struct {
	logger   Logger
	services *service.Services

	pb.OAuth2ServiceServer
}

func (s *oAuth2ServiceServer) SignUp(ctx context.Context, in *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	request := &oauth2.SignUpRequest{
		Username: in.Username,
		Password: in.Password,
	}

	_, err := s.services.OAuth2Service.SignUpRequestHandler.Handler(ctx, request)
	if errors.Is(err, oauth2.ErrSignUpUserAlreadyExist) {
		return nil, status.Error(codes.AlreadyExists, "A user with this username already exists")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "Internal")
	}

	return &pb.SignUpResponse{}, nil
}

func (s *oAuth2ServiceServer) SignIn(ctx context.Context, in *pb.SignInRequest) (*pb.SignInResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	request := &oauth2.SignInRequest{
		Username: in.Username,
		Password: in.Password,
	}

	response, err := s.services.OAuth2Service.SingInRequestHandler.Handle(ctx, request)
	if errors.Is(err, oauth2.ErrSignInUserNotFound) {
		return nil, status.Error(codes.NotFound, "User with that username was not found")
	} else if errors.Is(err, oauth2.ErrSignInInvalidPassword) {
		return nil, status.Error(codes.InvalidArgument, "Invalid password")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "Internal")
	}

	return &pb.SignInResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}, nil
}

func (s *oAuth2ServiceServer) SignOut(ctx context.Context, in *pb.SignOutRequest) (*pb.SignOutResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	request := &oauth2.SignOutRequest{
		AccessToken: in.AccessToken,
	}

	_, err := s.services.OAuth2Service.SignOutRequestHandler.Handle(ctx, request)
	if errors.Is(err, oauth2.ErrSignOutInvalidOrExpiredToken) {
		return nil, status.Error(codes.InvalidArgument, "Invalid or expired token")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "Internal")
	}

	return &pb.SignOutResponse{}, nil
}

func (s *oAuth2ServiceServer) Refresh(ctx context.Context, in *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	request := &oauth2.RefreshRequest{
		RefreshToken: in.RefreshToken,
	}

	response, err := s.services.OAuth2Service.RefreshRequestHandler.Handle(ctx, request)
	if errors.Is(err, oauth2.ErrRefreshInvalidOrExpiredToken) {
		return nil, status.Error(codes.InvalidArgument, "Invalid or expired token")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "Internal")
	}

	return &pb.RefreshResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}, nil
}
