package server_test

import (
	"context"
	"esl-challenge/api/gen/userpb"
	"esl-challenge/internal/entity"
	"esl-challenge/internal/repository"
	"esl-challenge/mocks"
	"esl-challenge/pkg/grpcclient"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
)

func testUserInit(t *testing.T, userRepository repository.UserRepository) (grpcclient.UserServiceClient, func()) {
	lis, grpcServer := initServer(t, userRepository)
	go grpcServer.Serve(lis)

	testEnd := func() {
		grpcServer.Stop()
	}

	userSvcCli, err := grpcclient.NewUserServiceClient(lis.Addr().String())
	if err != nil {
		testEnd()
		t.Fatal(err)
	}

	return userSvcCli, testEnd
}

type UserRequest interface {
	GetFirstName() string
	GetLastName() string
	GetNickname() string
	GetEmail() string
	GetCountry() string
}

func assertUser(t *testing.T, req UserRequest, user entity.User) {
	assert.Equal(t, req.GetFirstName(), user.FirstName)
	assert.Equal(t, req.GetLastName(), user.LastName)
	assert.Equal(t, req.GetNickname(), user.Nickname)
	assert.Equal(t, req.GetEmail(), user.Email)
	assert.Equal(t, req.GetCountry(), user.Country)
}

func TestUserService_CreateUser(t *testing.T) {
	type Test struct {
		TestName   string
		Request    *userpb.RequestCreateUser
		ExpectedId string
		CreateUser bool
		Validate   func(Test, *userpb.ResponseCreateUser, error)
	}
	tests := []Test{
		{
			TestName: "Success",
			Request: &userpb.RequestCreateUser{
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Password:  "password",
				Email:     "johndoe@email.com",
				Country:   "Germany",
			},
			ExpectedId: ulid.Make().String(),
			CreateUser: true,
			Validate: func(test Test, res *userpb.ResponseCreateUser, err error) {
				if assert.NoError(t, err) {
					assert.Equal(t, test.ExpectedId, res.Id)
				}
			},
		},
		{
			TestName: "Failure - Invalid argument",
			Request:  &userpb.RequestCreateUser{},
			Validate: func(test Test, res *userpb.ResponseCreateUser, err error) {
				assertGrpcErrorCode(t, err, codes.InvalidArgument)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			userRepository := mocks.NewUserRepository(t)
			if test.CreateUser {
				userRepository.On("CreateUser", mock.Anything, mock.Anything, test.Request.Password).
					Return(func(_ context.Context, user entity.User, _ string) entity.User {
						assertUser(t, test.Request, user)
						return entity.User{Id: test.ExpectedId}
					}, nil)
			}

			userSvcCli, testEnd := testUserInit(t, userRepository)
			defer testEnd()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			res, err := userSvcCli.CreateUser(ctx, test.Request)
			test.Validate(test, res, err)
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	req := &userpb.RequestUpdateUser{
		Id:        ulid.Make().String(),
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Password:  "password",
		Email:     "johndoe@email.com",
		Country:   "Germany",
	}

	userRepository := mocks.NewUserRepository(t)
	userRepository.On("UpdateUser", mock.Anything, req.Id, mock.Anything, req.Password).
		Return(func(_ context.Context, _ string, user entity.User, _ string) entity.User {
			assertUser(t, req, user)
			return entity.User{}
		}, nil)

	userSvcCli, testEnd := testUserInit(t, userRepository)
	defer testEnd()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := userSvcCli.UpdateUser(ctx, req)
	assert.NoError(t, err)
}
