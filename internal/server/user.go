package server

import (
	"context"
	"esl-challenge/api/gen/userpb"
	"esl-challenge/internal/entity"
	"esl-challenge/internal/repository"
)

const UserServiceName = "UserService"

type userServiceServer struct {
	userpb.UnimplementedUserServiceServer
	userRepository repository.UserRepository
}

type UserServiceServer struct {
	userpb.UserServiceServer
}

func NewUserServiceServer(userRepository repository.UserRepository) *userServiceServer {
	return &userServiceServer{userRepository: userRepository}
}

func (s userServiceServer) CreateUser(ctx context.Context, req *userpb.RequestCreateUser) (*userpb.ResponseCreateUser, error) {
	user, err := s.userRepository.CreateUser(ctx, entity.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Country:   req.Country,
	}, req.Password)
	if err != nil {
		return nil, err
	}

	return &userpb.ResponseCreateUser{Id: user.Id}, nil
}

func (s userServiceServer) UpdateUser(ctx context.Context, req *userpb.RequestUpdateUser) (*userpb.ResponseUpdateUser, error) {
	_, err := s.userRepository.UpdateUser(ctx, req.Id, entity.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Country:   req.Country,
	}, req.Password)
	if err != nil {
		return nil, err
	}

	return &userpb.ResponseUpdateUser{}, nil
}
