package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type HealthClient interface {
	grpc_health_v1.HealthClient
}

func NewHealthClient(address string) (HealthClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return grpc_health_v1.NewHealthClient(conn), nil
}
