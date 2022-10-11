package server

import "esl-challenge/api/gen/userpb"

const UserServiceName = "UserService"

type userServiceServer struct {
	userpb.UnimplementedUserServiceServer
}

func NewUserServiceServer() *userServiceServer {
	return &userServiceServer{}
}
