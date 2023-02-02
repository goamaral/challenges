package server

import (
	"challenge/api/gen/userpb"
	"challenge/internal/entity"
	"challenge/internal/protobuf"
	"challenge/internal/repository"
	"challenge/internal/service"
	"challenge/pkg/gormprovider"
	"context"

	"github.com/jackc/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const UserServiceName = "UserService"

type userServiceServer struct {
	userpb.UnimplementedUserServiceServer
	userRepo    repository.UserRepository
	rabbitmqSvc service.RabbitmqService
}

type UserServiceServer struct {
	userpb.UserServiceServer
}

func NewUserServiceServer(userRepo repository.UserRepository, rabbitmqSvc service.RabbitmqService) *userServiceServer {
	return &userServiceServer{userRepo: userRepo, rabbitmqSvc: rabbitmqSvc}
}

func (s userServiceServer) CreateUser(ctx context.Context, req *userpb.RequestCreateUser) (*userpb.ResponseCreateUser, error) {
	var user entity.User
	var err error
	err = s.userRepo.RunInTransaction(ctx, func(txCtx context.Context) error {
		// Create user
		user, err = s.userRepo.CreateUser(txCtx, entity.User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Nickname:  req.Nickname,
			Email:     req.Email,
			Country:   req.Country,
		}, req.Password)
		if err != nil {
			return err
		}

		// Publish user changes
		err = s.rabbitmqSvc.PublishChanges(txCtx, user, service.EntityType_USER, service.Action_CREATE)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if gormprovider.IsUniqueViolationError(err) {
			pgErr := err.(*pgconn.PgError)
			return nil, status.Error(codes.FailedPrecondition, pgErr.Detail)
		}

		return nil, err
	}

	return &userpb.ResponseCreateUser{Id: user.Id}, nil
}

func (s userServiceServer) UpdateUser(ctx context.Context, req *userpb.RequestUpdateUser) (*userpb.ResponseUpdateUser, error) {
	var user entity.User
	var err error
	err = s.userRepo.RunInTransaction(ctx, func(txCtx context.Context) error {
		user, err = s.userRepo.UpdateUser(ctx, req.Id, entity.User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Nickname:  req.Nickname,
			Email:     req.Email,
			Country:   req.Country,
		}, req.Password)
		if err != nil {
			return err
		}

		// Publish user changes
		err = s.rabbitmqSvc.PublishChanges(txCtx, user, service.EntityType_USER, service.Action_UPDATE)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &userpb.ResponseUpdateUser{}, nil
}

func (s userServiceServer) DeleteUser(ctx context.Context, req *userpb.RequestDeleteUser) (*userpb.ResponseDeleteUser, error) {
	var err error
	err = s.userRepo.RunInTransaction(ctx, func(txCtx context.Context) error {
		// Delete user
		err := s.userRepo.DeleteUser(ctx, req.Id)
		if err != nil {
			return err
		}

		// Publish user changes
		err = s.rabbitmqSvc.PublishChanges(txCtx, req.Id, service.EntityType_USER, service.Action_DELETE)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &userpb.ResponseDeleteUser{}, nil
}

func (s userServiceServer) ListUsers(ctx context.Context, req *userpb.RequestListUsers) (*userpb.ResponseListUsers, error) {
	users, err := s.userRepo.ListUsers(ctx, req.PaginationToken, uint(req.PageSize), &repository.ListUsersOpts{
		Country: req.Country,
	})
	if err != nil {
		return nil, err
	}

	return &userpb.ResponseListUsers{Users: protobuf.EntitiesToProtobuf(users, protobuf.UserToProtobuf)}, nil
}
