package server_test

import (
	"challenge/internal/repository"
	"challenge/internal/server"
	"challenge/internal/service"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func initServer(t *testing.T, userRepo repository.UserRepository, rabbitmqSvc service.RabbitmqService) (net.Listener, *server.Server) {
	grpcServer, err := server.NewServer(userRepo, rabbitmqSvc)
	if err != nil {
		t.Fatal(err)
	}
	lis, err := net.Listen("tcp", "")
	if err != nil {
		t.Fatal(err)
	}
	return lis, grpcServer
}

func assertGrpcErrorCode(t *testing.T, err error, c codes.Code) {
	if assert.Error(t, err) {
		sts, ok := status.FromError(err)
		assert.True(t, ok, "Not grpc error")
		assert.Equal(t, c, sts.Code(), "wrong grpc status code")
	}
}
