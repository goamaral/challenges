package server

import (
	"challenge/api/gen/userpb"
	"challenge/internal/entity"
	"challenge/internal/protobuf"
	"challenge/internal/repository"
	"challenge/internal/service"
	"context"

	"gorm.io/gorm/clause"
)

const UserServiceName = "UserService"

type UserServiceServer struct {
	userpb.UnimplementedUserServiceServer
	userRepo  repository.UserRepository
	rabbitSvc service.RabbitmqService
}

func NewUserServiceServer(userRepo repository.UserRepository, rabbitSvc service.RabbitmqService) UserServiceServer {
	return UserServiceServer{userRepo: userRepo, rabbitSvc: rabbitSvc}
}

func (s UserServiceServer) CreateUser(ctx context.Context, req *userpb.RequestCreateUser) (*userpb.ResponseCreateUser, error) {
	var user entity.User
	var err error
	err = s.userRepo.NewTransaction(ctx, func(ctx context.Context) error {
		user, err = s.userRepo.CreateUser(
			ctx,
			entity.User{
				FirstName: req.FirstName,
				LastName:  req.LastName,
				Nickname:  req.Nickname,
				Email:     req.Email,
				Country:   req.Country,
			},
			req.Password,
		)
		if err != nil {
			return err
		}

		err = s.rabbitSvc.Publish(ctx, &userpb.Event_UserCreated{User: protobuf.UserToProtobuf(user)})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		// TODO: Check what this unique violation means
		// if errors.Is(err, gorm_ext.ErrUniqueViolation) {
		// 	log.Warn().Err(err).Msg("failed to list users")
		// 	return nil, status.Error(codes.FailedPrecondition, "TODO")
		// }

		return nil, err
	}

	return &userpb.ResponseCreateUser{Id: user.Id}, nil
}

func (s UserServiceServer) PatchUser(ctx context.Context, req *userpb.RequestPatchUser) (*userpb.ResponsePatchUser, error) {
	err := s.userRepo.NewTransaction(ctx, func(txCtx context.Context) error {
		err := s.userRepo.PatchUser(ctx, req.Id, protobuf.UserPatchFromProtobuf(req.Patch))
		if err != nil {
			return err
		}

		if req.Patch.Password != nil {
			req.Patch.Password.Value = "redacted"
		}

		err = s.rabbitSvc.Publish(txCtx, &userpb.Event_UserPatched{Patch: req.Patch})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &userpb.ResponsePatchUser{}, nil
}

func (s UserServiceServer) DeleteUser(ctx context.Context, req *userpb.RequestDeleteUser) (*userpb.ResponseDeleteUser, error) {
	err := s.userRepo.NewTransaction(ctx, func(ctx context.Context) error {
		err := s.userRepo.DeleteUser(ctx, req.Id)
		if err != nil {
			return err
		}

		err = s.rabbitSvc.Publish(ctx, &userpb.Event_UserDeleted{Id: req.Id})
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

func (s UserServiceServer) ListUsers(ctx context.Context, req *userpb.RequestListUsers) (*userpb.ResponseListUsers, error) {
	var conds []clause.Expression
	if req.Country != "" {
		conds = append(conds, clause.Eq{Column: "country", Value: req.Country})
	}

	// TODO: Validate request
	users, err := s.userRepo.ListUsers(ctx, req.PaginationToken, uint(req.PageSize), conds...)
	if err != nil {
		return nil, err
	}

	return &userpb.ResponseListUsers{Users: protobuf.EntitiesToProtobuf(users, protobuf.UserToProtobuf)}, nil
}
