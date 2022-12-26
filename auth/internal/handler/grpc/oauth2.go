package grpc

import (
	"context"

	pb "github.com/udholdenhed/unotes/auth/api/proto"
)

type oAuth2ServerGRPC struct {
	pb.OAuth2ServiceServer
}

func (oAuth2ServerGRPC) SignUp(ctx context.Context, request *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	return &pb.SignUpResponse{}, nil
}

func (oAuth2ServerGRPC) SignIn(ctx context.Context, request *pb.SignInRequest) (*pb.SignInResponse, error) {
	return &pb.SignInResponse{
		AccessToken:  "OK",
		RefreshToken: "Ok!",
	}, nil
}
