package server_test

import (
	"challenge/api/gen/userpb"
	"challenge/internal/entity"
	"challenge/internal/repository"
	"challenge/internal/service"
	"challenge/mocks"
	"challenge/pkg/gorm_ext"
	"challenge/pkg/grpcclient"
	"context"
	"testing"
	"time"

	"github.com/jackc/pgconn"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm/clause"
)

func testUserInit(t *testing.T, userRepo repository.UserRepository, rabbitmqSvc service.RabbitmqService) (grpcclient.UserServiceClient, func()) {
	lis, grpcServer := initServer(t, userRepo, rabbitmqSvc)
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
	ctx := context.Background()

	validReq := &userpb.RequestCreateUser{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Password:  "password",
		Email:     "johndoe@email.com",
		Country:   "Germany",
	}

	t.Run("Success", func(t *testing.T) {
		userId := ulid.Make().String()

		userRepo := mocks.NewUserRepository(t)
		userRepo.On("CreateUser", mock.Anything, mock.Anything, validReq.Password).
			Return(func(_ context.Context, user entity.User, _ string) entity.User {
				assertUser(t, validReq, user)
				user.Id = userId
				return user
			}, nil)

		rabbitmqSvc := mocks.NewRabbitmqService(t)
		rabbitmqSvc.On("PublishChanges", mock.Anything, mock.Anything, service.EntityType_USER, service.Action_CREATE).
			Return(func(_ context.Context, user any, _ service.EntityType, _ service.Action) error {
				assert.Equal(t, userId, user.(entity.User).Id)
				return nil
			})

		userSvcCli, testEnd := testUserInit(t, userRepo, rabbitmqSvc)
		defer testEnd()

		res, err := userSvcCli.CreateUser(ctx, validReq)
		if assert.NoError(t, err) {
			assert.Equal(t, userId, res.Id)
		}
	})

	t.Run("Invalid argument - Missing params", func(t *testing.T) {
		userSvcCli, testEnd := testUserInit(t, repository.UserRepository{}, nil)
		defer testEnd()

		_, err := userSvcCli.CreateUser(ctx, &userpb.RequestCreateUser{})
		assertGrpcErrorCode(t, err, codes.InvalidArgument)
	})

	t.Run("FailedPrecondition - Duplicate key", func(t *testing.T) {
		userRepo := mocks.NewUserRepository(t)
		userRepo.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
			Return(entity.User{}, &pgconn.PgError{Code: string(gorm_ext.PostgresErrorCode_UNIQUE_VIOLATION)})

		userSvcCli, testEnd := testUserInit(t, userRepo, nil)
		defer testEnd()

		_, err := userSvcCli.CreateUser(ctx, validReq)
		assertGrpcErrorCode(t, err, codes.FailedPrecondition)
	})
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

	updatedUser := entity.User{}

	userRepo := mocks.NewUserRepository(t)
	rabbitmqSvc := mocks.NewRabbitmqService(t)
	userRepo.On("UpdateUser", mock.Anything, req.Id, mock.Anything, req.Password).
		Return(func(_ context.Context, _ string, user entity.User, _ string) entity.User {
			assertUser(t, req, user)
			return updatedUser
		}, nil)
	rabbitmqSvc.On("PublishChanges", mock.Anything, updatedUser, service.EntityType_USER, service.Action_UPDATE).Return(nil)

	userSvcCli, testEnd := testUserInit(t, userRepo, rabbitmqSvc)
	defer testEnd()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := userSvcCli.UpdateUser(ctx, req)
	assert.NoError(t, err)
}

func TestUserService_DeleteUser(t *testing.T) {
	type Test struct {
		TestName   string
		Request    *userpb.RequestDeleteUser
		DeleteUser bool
		Validate   func(Test, error)
	}
	tests := []Test{
		{
			TestName: "Success",
			Request: &userpb.RequestDeleteUser{
				Id: ulid.Make().String(),
			},
			DeleteUser: true,
			Validate: func(test Test, err error) {
				assert.NoError(t, err)
			},
		},
		{
			TestName: "Invalid argument",
			Request:  &userpb.RequestDeleteUser{},
			Validate: func(test Test, err error) {
				assertGrpcErrorCode(t, err, codes.InvalidArgument)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			userRepo := mocks.NewUserRepository(t)
			rabbitmqSvc := mocks.NewRabbitmqService(t)
			if test.DeleteUser {
				userRepo.On("DeleteUser", mock.Anything, test.Request.Id).Return(nil)
				rabbitmqSvc.On("PublishChanges", mock.Anything, test.Request.Id, service.EntityType_USER, service.Action_DELETE).Return(nil)
			}

			userSvcCli, testEnd := testUserInit(t, userRepo, rabbitmqSvc)
			defer testEnd()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := userSvcCli.DeleteUser(ctx, test.Request)
			test.Validate(test, err)
		})
	}
}

func TestUserService_ListUsers(t *testing.T) {
	paginationToken := "paginationToken"
	var pageSize uint = 5
	country := "country"
	user := entity.User{Country: country}

	userRepo := mocks.NewUserRepository(t)
	userRepo.On("ListUsers", mock.Anything, paginationToken, pageSize, mock.Anything).
		Return(func(_ context.Context, _ string, _ uint, conds ...clause.Expression) []entity.User {
			if assert.Len(t, conds, 1) {
				filterOpt, ok := opts[0].(repository.UserFilterOption)
				if assert.True(t, ok, "not repository.UserFilterOption") && assert.NotNil(t, filterOpt.Country) {
					assert.Equal(t, country, *filterOpt.Country)
				}
			}
			return []entity.User{{Country: country}}
		}, nil)

	userSvcCli, testEnd := testUserInit(t, userRepo, nil)
	defer testEnd()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := userSvcCli.ListUsers(ctx, &userpb.RequestListUsers{PaginationToken: paginationToken, PageSize: uint32(pageSize), Country: country})
	if assert.NoError(t, err) && assert.Len(t, res.Users, 1) {
		assertUser(t, res.Users[0], user)
	}
}
