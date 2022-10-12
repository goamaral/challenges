package server_test

import (
	"esl-challenge/internal/repository"
	"esl-challenge/internal/server"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func initServer(t *testing.T, userRepository repository.UserRepository) (net.Listener, *server.Server) {
	grpcServer, err := server.NewServer(userRepository)
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
		grpcErr, ok := err.(interface {
			GRPCStatus() *status.Status
		})
		assert.True(t, ok, "Not grpc error")
		assert.Equal(t, c, grpcErr.GRPCStatus().Code())
	}
}
