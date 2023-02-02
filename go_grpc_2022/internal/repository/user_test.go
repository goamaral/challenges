package repository_test

import (
	"challenge/internal/entity"
	"challenge/internal/repository"
	"challenge/pkg/gormprovider"
	"context"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func addUser(t *testing.T, db *gorm.DB, user entity.User, password string) entity.User {
	if password == "" {
		user.SetPassword("password")
	} else {
		user.SetPassword(password)
	}
	if user.Id == "" {
		user.Id = ulid.Make().String()
	}

	err := db.Create(&user).Error
	if err != nil {
		t.Fatal(err)
	}

	return user
}

func TestUserRepository_CreateUser(t *testing.T) {
	nickname := "nickname"
	email := "user@email.com"

	type Test struct {
		TestName string
		User     entity.User
		Password string
		Validate func(Test, entity.User, error)
	}
	tests := []Test{
		{
			TestName: "Success",
			User:     entity.User{},
			Password: "password",
			Validate: func(test Test, user entity.User, err error) {
				if assert.NoError(t, err) {
					assert.NotZero(t, user.Id)
					assert.NotZero(t, user.CreatedAt)
					assert.NotZero(t, user.UpdatedAt)
					assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(test.Password)))
				}
			},
		},
		{
			TestName: "Failure - Nickname not unique",
			User:     entity.User{Nickname: nickname},
			Validate: func(test Test, user entity.User, err error) {
				if assert.Error(t, err) {
					assert.Truef(t, gormprovider.IsUniqueViolationError(err), "not pg unique_violation: %s", err.Error())
				}
			},
		},
		{
			TestName: "Failure - Email not unique",
			User:     entity.User{Email: email},
			Validate: func(test Test, user entity.User, err error) {
				if assert.Error(t, err) {
					assert.Truef(t, gormprovider.IsUniqueViolationError(err), "not pg unique_violation: %s", err.Error())
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			db, testEnd := testInit(t)
			defer testEnd()

			addUser(t, db, entity.User{Nickname: nickname, Email: email}, "")

			r := repository.NewUserRepository(db)
			user, err := r.CreateUser(context.Background(), test.User, test.Password)
			test.Validate(test, user, err)
		})
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {
	nickname := "nickname"
	email := "user@email.com"

	type Test struct {
		TestName       string
		User           entity.User
		Password       string
		UserUpdates    entity.User
		PasswordUpdate string
		Validate       func(Test, entity.User, error)
	}
	tests := []Test{
		{
			TestName:    "Success - Update single field (not password)",
			User:        entity.User{FirstName: "John", LastName: "Doe"},
			UserUpdates: entity.User{FirstName: "Joe"},
			Validate: func(test Test, user entity.User, err error) {
				if assert.NoError(t, err) {
					assert.Equal(t, test.UserUpdates.FirstName, user.FirstName) // Updates
					assert.Equal(t, test.User.LastName, user.LastName)          // Does not update
				}
			},
		},
		{
			TestName:       "Success - Update password",
			Password:       "password",
			PasswordUpdate: "new_password",
			Validate: func(test Test, user entity.User, err error) {
				if assert.NoError(t, err) {
					assert.NoError(t, bcrypt.CompareHashAndPassword(user.EncryptedPassword, []byte(test.PasswordUpdate)))
				}
			},
		},
		{
			TestName:    "Failure - Nickname not unique",
			UserUpdates: entity.User{Nickname: nickname},
			Validate: func(test Test, user entity.User, err error) {
				if assert.Error(t, err) {
					assert.Truef(t, gormprovider.IsUniqueViolationError(err), "not pg unique_violation: %s", err.Error())
				}
			},
		},
		{
			TestName:    "Failure - Email not unique",
			UserUpdates: entity.User{Email: email},
			Validate: func(test Test, user entity.User, err error) {
				if assert.Error(t, err) {
					assert.Truef(t, gormprovider.IsUniqueViolationError(err), "not pg unique_violation: %s", err.Error())
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			db, testEnd := testInit(t)
			defer testEnd()

			addUser(t, db, entity.User{Nickname: nickname, Email: email}, "")
			user := addUser(t, db, test.User, test.Password)

			r := repository.NewUserRepository(db)
			user, err := r.UpdateUser(context.Background(), user.Id, test.UserUpdates, test.PasswordUpdate)
			test.Validate(test, user, err)
		})
	}
}

func TestUserRepository_DeleteUser(t *testing.T) {
	db, testEnd := testInit(t)
	defer testEnd()

	type Test struct {
		TestName string
		Setup    func() string
		Validate func(Test, error)
	}
	tests := []Test{
		{
			TestName: "Delete existing user",
			Setup: func() string {
				user := addUser(t, db, entity.User{}, "")
				return user.Id
			},
			Validate: func(test Test, err error) {
				assert.NoError(t, err)
			},
		},
		{
			TestName: "Delete non existing user",
			Setup: func() string {
				return ulid.Make().String()
			},
			Validate: func(test Test, err error) {
				assert.NoError(t, err)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			id := test.Setup()

			r := repository.NewUserRepository(db)
			err := r.DeleteUser(context.Background(), id)
			test.Validate(test, err)
		})
	}
}

func TestUserRepository_ListUsers(t *testing.T) {
	country := "country"
	firstUserId := ulid.Make().String()
	secondUserId := ulid.Make().String()

	type Test struct {
		TestName        string
		PaginationToken string
		PageSize        uint
		Country         string
		Validate        func(Test, []entity.User, error)
	}
	tests := []Test{
		{
			TestName: "Empty args", // Should return both users
			Validate: func(test Test, users []entity.User, err error) {
				if assert.NoError(t, err) && assert.Len(t, users, 2) {
					assert.Equal(t, firstUserId, users[0].Id)
					assert.Equal(t, secondUserId, users[1].Id)
				}
			},
		},
		{
			TestName:        "With first user id as paginationToken", // Should only return the second user
			PaginationToken: firstUserId,
			Validate: func(test Test, users []entity.User, err error) {
				if assert.NoError(t, err) && assert.Len(t, users, 1) {
					assert.Equal(t, secondUserId, users[0].Id)
				}
			},
		},
		{
			TestName:        "With second user id as paginationToken", // Should return no users
			PaginationToken: secondUserId,
			Validate: func(test Test, users []entity.User, err error) {
				if assert.NoError(t, err) {
					assert.Len(t, users, 0)
				}
			},
		},
		{
			TestName: "With first user country as country", // Should only return the first user
			Country:  country,
			Validate: func(test Test, users []entity.User, err error) {
				if assert.NoError(t, err) && assert.Len(t, users, 1) {
					assert.Equal(t, firstUserId, users[0].Id)
				}
			},
		},
		{
			TestName: "With pageSize set to 1", // Should only return the first user
			PageSize: 1,
			Validate: func(test Test, users []entity.User, err error) {
				if assert.NoError(t, err) && assert.Len(t, users, 1) {
					assert.Equal(t, firstUserId, users[0].Id)
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			var filterOpt repository.UserFilterOption
			if test.Country != "" {
				filterOpt = repository.UserFilterOption{Country: &test.Country}
			}

			db, testEnd := testInit(t)
			defer testEnd()

			addUser(t, db, entity.User{Id: firstUserId, Nickname: "nickname1", Email: "email1", Country: country}, "")
			addUser(t, db, entity.User{Id: secondUserId, Nickname: "nickname2", Email: "email2"}, "")

			r := repository.NewUserRepository(db)
			users, err := r.ListUsers(context.Background(), test.PaginationToken, test.PageSize, filterOpt)
			test.Validate(test, users, err)
		})
	}
}
