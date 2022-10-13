package repository_test

import (
	"context"
	"esl-challenge/internal/entity"
	"esl-challenge/internal/repository"
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
	if user.Nickname == "" {
		user.Nickname = ulid.Make().String()
	}
	if user.Email == "" {
		user.Email = ulid.Make().String()
	}

	err := db.Create(&user).Error
	if err != nil {
		t.Fatal(err)
	}

	return user
}

func TestUserRepository_CreateUser(t *testing.T) {
	db, testEnd := testInit(t)
	defer testEnd()

	existingUser := addUser(t, db, entity.User{}, "")

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
			User:     entity.User{Nickname: existingUser.Nickname},
			Validate: func(test Test, user entity.User, err error) {
				assertUniqueViolationError(t, err)
			},
		},
		{
			TestName: "Failure - Email not unique",
			User:     entity.User{Email: existingUser.Email},
			Validate: func(test Test, user entity.User, err error) {
				assertUniqueViolationError(t, err)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			r := repository.NewUserRepository(db)
			user, err := r.CreateUser(context.Background(), test.User, test.Password)
			test.Validate(test, user, err)
		})
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {
	db, testEnd := testInit(t)
	defer testEnd()

	existingUser := addUser(t, db, entity.User{}, "")

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
			UserUpdates: entity.User{Nickname: existingUser.Nickname},
			Validate: func(test Test, user entity.User, err error) {
				assertUniqueViolationError(t, err)
			},
		},
		{
			TestName:    "Failure - Email not unique",
			UserUpdates: entity.User{Email: existingUser.Email},
			Validate: func(test Test, user entity.User, err error) {
				assertUniqueViolationError(t, err)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
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
	db, testEnd := testInit(t)
	defer testEnd()

	country := "country"

	firstUser := addUser(t, db, entity.User{Country: country}, "")
	secondUser := addUser(t, db, entity.User{}, "")

	type Test struct {
		TestName        string
		PaginationToken string
		Country         string
		Validate        func(Test, []entity.User, error)
	}
	tests := []Test{
		{
			TestName: "Empty args", // Should return both users
			Validate: func(test Test, users []entity.User, err error) {
				if assert.NoError(t, err) && assert.Len(t, users, 2) {
					assert.Equal(t, firstUser.Id, users[0].Id)
					assert.Equal(t, secondUser.Id, users[1].Id)
				}
			},
		},
		{
			TestName:        "With first user id as paginationToken", // Should only return the second user
			PaginationToken: firstUser.Id,
			Validate: func(test Test, users []entity.User, err error) {
				if assert.NoError(t, err) && assert.Len(t, users, 1) {
					assert.Equal(t, secondUser.Id, users[0].Id)
				}
			},
		},
		{
			TestName:        "With second user id as paginationToken", // Should return no users
			PaginationToken: secondUser.Id,
			Validate: func(test Test, users []entity.User, err error) {
				if assert.NoError(t, err) {
					assert.Len(t, users, 0)
				}
			},
		},
		{
			TestName: "With first user country as country", // Should only return the first user
			Country:  firstUser.Country,
			Validate: func(test Test, users []entity.User, err error) {
				if assert.NoError(t, err) && assert.Len(t, users, 1) {
					assert.Equal(t, firstUser.Id, users[0].Id)
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			var opts *repository.ListUsersOpts
			if test.Country != "" {
				opts = &repository.ListUsersOpts{Country: test.Country}
			}

			r := repository.NewUserRepository(db)
			users, err := r.ListUsers(context.Background(), test.PaginationToken, opts)
			test.Validate(test, users, err)
		})
	}
}
