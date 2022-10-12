package server_test

import (
	"context"
	"esl-challenge/api/gen/userpb"
	"esl-challenge/internal/entity"
	"esl-challenge/mocks"
	"esl-challenge/pkg/grpcclient"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUserService_CreateUser(t *testing.T) {
	userRepository := mocks.NewUserRepository(t)
	lis, grpcServer := initServer(t, userRepository)
	defer grpcServer.Stop()
	go grpcServer.Serve(lis)

	rawPassword := "password"
	expectedId := ulid.Make().String()

	userRepository.On("CreateUser", mock.Anything, mock.Anything, rawPassword).
		Return(entity.User{Id: expectedId}, nil)

	userSvcCli, err := grpcclient.NewUserServiceClient(lis.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	// Success
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := userSvcCli.CreateUser(ctx, &userpb.RequestCreateUser{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Password:  rawPassword,
		Email:     "johndoe@email.com",
		Country:   "Germany",
	})
	if assert.NoError(t, err) {
		assert.Equal(t, expectedId, res.Id)
	}

	// Failure - InvalidArgument
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = userSvcCli.CreateUser(ctx, &userpb.RequestCreateUser{})
	if assert.Error(t, err) {
		grpcErr, ok := err.(interface {
			GRPCStatus() *status.Status
		})
		assert.True(t, ok, "Not grpc error")
		assert.Equal(t, codes.InvalidArgument, grpcErr.GRPCStatus().Code())
	}
}
