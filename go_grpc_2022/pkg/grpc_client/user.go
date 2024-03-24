package grpc_client

import (
	"challenge/api/gen/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServiceClient interface {
	userpb.UserServiceClient
}

func NewUserServiceClient(address string) (UserServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return userpb.NewUserServiceClient(conn), nil
}
