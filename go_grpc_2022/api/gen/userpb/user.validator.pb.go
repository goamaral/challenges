// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: user.proto

package userpb

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	_ "google.golang.org/protobuf/types/known/wrapperspb"
	_ "github.com/mwitkow/go-proto-validators"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *User) Validate() error {
	if this.CreatedAt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.CreatedAt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("CreatedAt", err)
		}
	}
	if this.UpdatedAt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.UpdatedAt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("UpdatedAt", err)
		}
	}
	return nil
}
func (this *User_Patch) Validate() error {
	if this.FirstName != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.FirstName); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("FirstName", err)
		}
	}
	if this.LastName != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.LastName); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("LastName", err)
		}
	}
	if this.Nickname != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Nickname); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Nickname", err)
		}
	}
	if this.Password != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Password); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Password", err)
		}
	}
	if this.Email != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Email); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Email", err)
		}
	}
	if this.Country != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Country); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Country", err)
		}
	}
	return nil
}
func (this *RequestCreateUser) Validate() error {
	if this.FirstName == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("FirstName", fmt.Errorf(`value '%v' must not be an empty string`, this.FirstName))
	}
	if this.LastName == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("LastName", fmt.Errorf(`value '%v' must not be an empty string`, this.LastName))
	}
	if this.Nickname == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Nickname", fmt.Errorf(`value '%v' must not be an empty string`, this.Nickname))
	}
	if this.Password == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Password", fmt.Errorf(`value '%v' must not be an empty string`, this.Password))
	}
	if this.Email == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Email", fmt.Errorf(`value '%v' must not be an empty string`, this.Email))
	}
	if this.Country == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Country", fmt.Errorf(`value '%v' must not be an empty string`, this.Country))
	}
	return nil
}
func (this *ResponseCreateUser) Validate() error {
	return nil
}
func (this *RequestPatchUser) Validate() error {
	if this.Id == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Id", fmt.Errorf(`value '%v' must not be an empty string`, this.Id))
	}
	if nil == this.Patch {
		return github_com_mwitkow_go_proto_validators.FieldError("Patch", fmt.Errorf("message must exist"))
	}
	if this.Patch != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Patch); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Patch", err)
		}
	}
	return nil
}
func (this *ResponsePatchUser) Validate() error {
	return nil
}
func (this *RequestDeleteUser) Validate() error {
	if this.Id == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Id", fmt.Errorf(`value '%v' must not be an empty string`, this.Id))
	}
	return nil
}
func (this *ResponseDeleteUser) Validate() error {
	return nil
}
func (this *RequestListUsers) Validate() error {
	return nil
}
func (this *ResponseListUsers) Validate() error {
	for _, item := range this.Users {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Users", err)
			}
		}
	}
	return nil
}
func (this *Event) Validate() error {
	return nil
}
func (this *Event_UserCreated) Validate() error {
	if this.User != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.User); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("User", err)
		}
	}
	return nil
}
func (this *Event_UserPatched) Validate() error {
	if this.Patch != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Patch); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Patch", err)
		}
	}
	return nil
}
func (this *Event_UserDeleted) Validate() error {
	return nil
}
