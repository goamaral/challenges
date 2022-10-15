package server_test

import (
	"context"
	"esl-challenge/internal/repository"
	"esl-challenge/internal/server"
	"esl-challenge/internal/service"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		grpcErr, ok := err.(interface {
			GRPCStatus() *status.Status
		})
		assert.True(t, ok, "Not grpc error")
		assert.Equal(t, c, grpcErr.GRPCStatus().Code())
	}
}

func defineRunInTransactionStub(m mock.Mock) mock.Mock {
	m.On("RunInTransaction", mock.Anything, mock.Anything).
		Return(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	return m
}
