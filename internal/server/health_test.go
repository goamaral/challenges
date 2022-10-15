package server_test

import (
	"context"
	"esl-challenge/internal/server"
	"esl-challenge/pkg/grpcclient"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func testHealthInit(t *testing.T) (grpcclient.HealthServiceClient, *server.Server, func()) {
	lis, grpcServer := initServer(t, nil, nil)
	go grpcServer.Serve(lis)

	testEnd := func() {
		grpcServer.Stop()
	}

	healthSvcCli, err := grpcclient.NewHealthServiceClient(lis.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	return healthSvcCli, grpcServer, testEnd
}

func TestHealthService_Check(t *testing.T) {
	healthSvcCli, grpcServer, testEnd := testHealthInit(t)
	defer testEnd()

	t.Run("Check serving server health", func(t *testing.T) {
		res, err := healthSvcCli.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
		if assert.NoError(t, err) {
			assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, res.Status)
		}
	})

	t.Run("Check serving service health", func(t *testing.T) {
		res, err := healthSvcCli.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: server.UserServiceName})
		if assert.NoError(t, err) {
			assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, res.Status)
		}
	})

	t.Run("Check not serving service health", func(t *testing.T) {
		grpcServer.SetServingStatus(server.UserServiceName, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		res, err := healthSvcCli.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: server.UserServiceName})
		if assert.NoError(t, err) {
			assert.Equal(t, grpc_health_v1.HealthCheckResponse_NOT_SERVING, res.Status)
		}
	})

	t.Run("Check unknown service health", func(t *testing.T) {
		_, err := healthSvcCli.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "UnknownService"})
		assert.Error(t, err)
	})

	t.Run("Server down", func(t *testing.T) {
		downHealthSvcCli, err := grpcclient.NewHealthServiceClient("")
		if err != nil {
			t.Fatal(err)
		}
		_, err = downHealthSvcCli.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
		assert.Error(t, err)
	})
}

func TestHealthService_Watch(t *testing.T) {
	lis, grpcServer := initServer(t, nil, nil)
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
