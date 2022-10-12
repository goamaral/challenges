package server_test

import (
	"context"
	"esl-challenge/internal/server"
	"esl-challenge/pkg/grpcclient"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestIntegrationHealthService_Check(t *testing.T) {
	lis, grpcServer := initServer(t, nil)
	defer grpcServer.Stop()
	go grpcServer.Serve(lis)

	// Initialize health client
	healthSvcCli, err := grpcclient.NewHealthServiceClient(lis.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	// Check serving server health
	res, err := healthSvcCli.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	if assert.NoError(t, err) {
		assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, res.Status)
	}

	// Check serving service health
	res, err = healthSvcCli.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: server.UserServiceName})
	if assert.NoError(t, err) {
		assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, res.Status)
	}

	// Check not serving service health
	grpcServer.SetServingStatus(server.UserServiceName, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	res, err = healthSvcCli.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: server.UserServiceName})
	if assert.NoError(t, err) {
		assert.Equal(t, grpc_health_v1.HealthCheckResponse_NOT_SERVING, res.Status)
	}

	// Check unknown service health
	_, err = healthSvcCli.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "UnknownService"})
	assert.Error(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_NOT_SERVING, res.Status)

	// Server down
	downHealthSvcCli, err := grpcclient.NewHealthServiceClient("")
	if err != nil {
		t.Fatal(err)
	}
	res, err = downHealthSvcCli.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestIntegrationHealthService_Watch(t *testing.T) {
	lis, grpcServer := initServer(t, nil)
	defer grpcServer.Stop()
	go grpcServer.Serve(lis)

	// Initialize health client
	healthSvcCli, err := grpcclient.NewHealthServiceClient(lis.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	// Check server health
	stream, err := healthSvcCli.Watch(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	if assert.NoError(t, err) {
		res, err := stream.Recv()
		if assert.NoError(t, err) {
			assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, res.Status)
		}
	}

	// Check service health
	stream, err = healthSvcCli.Watch(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: server.UserServiceName})
	if assert.NoError(t, err) {
		res, err := stream.Recv()
		if assert.NoError(t, err) {
			assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, res.Status)
		}
	}

	// Check not serving service health
	grpcServer.SetServingStatus(server.UserServiceName, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	stream, err = healthSvcCli.Watch(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: server.UserServiceName})
	if assert.NoError(t, err) {
		res, err := stream.Recv()
		if assert.NoError(t, err) {
			assert.Equal(t, grpc_health_v1.HealthCheckResponse_NOT_SERVING, res.Status)
		}
	}

	// Check unknown service health
	stream, err = healthSvcCli.Watch(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "UnknownService"})
	if assert.NoError(t, err) {
		res, err := stream.Recv()
		if assert.NoError(t, err) {
			assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN, res.Status)
		}
	}

	// Server down
	downHealthSvcCli, err := grpcclient.NewHealthServiceClient("")
	if err != nil {
		t.Fatal(err)
	}
	stream, err = downHealthSvcCli.Watch(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	assert.Error(t, err)
	assert.Nil(t, stream)
}
