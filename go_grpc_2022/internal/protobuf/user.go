package protobuf

import (
	"challenge/api/gen/userpb"
	"challenge/internal/entity"
	"challenge/internal/repository"
	"challenge/pkg/protobuf_ext"

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

func UserPatchFromProtobuf(patch *userpb.User_Patch) repository.UserPatch {
	return repository.UserPatch{
		FirstName: protobuf_ext.WrappedValueToOption(patch.FirstName),
		LastName:  protobuf_ext.WrappedValueToOption(patch.LastName),
		Nickname:  protobuf_ext.WrappedValueToOption(patch.Nickname),
		Email:     protobuf_ext.WrappedValueToOption(patch.Email),
		Country:   protobuf_ext.WrappedValueToOption(patch.Country),
		Password:  protobuf_ext.WrappedValueToOption(patch.Password),
	}
}
