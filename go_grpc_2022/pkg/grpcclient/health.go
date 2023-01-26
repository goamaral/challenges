package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type HealthServiceClient interface {
	grpc_health_v1.HealthClient
}

func NewHealthServiceClient(address string) (HealthServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return grpc_health_v1.NewHealthClient(conn), nil
}
