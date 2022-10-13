package protobuf

import (
	"esl-challenge/api/gen/userpb"
	"esl-challenge/internal/entity"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func UserToProtobuf(user entity.User) *userpb.User {
	return &userpb.User{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Country:   user.Country,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
