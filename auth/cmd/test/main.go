package main

import (
	"context"
	"fmt"

	pb "github.com/udholdenhed/unotes/auth/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const addr = "0.0.0.0:8091"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer func(conn *grpc.ClientConn) { _ = conn.Close() }(conn)

	client := pb.NewOAuth2ServiceClient(conn)
	DoSum(client)
}

func DoSum(client pb.OAuth2ServiceClient) {
	response, err := client.SignIn(context.Background(), &pb.SignInRequest{})
	if err != nil {
		panic(err)
	}

	fmt.Println(response.RefreshToken)
}
