package repository

import (
	"context"
	"esl-challenge/internal/entity"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const bcryptCost = 15

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User, rawPassword string) (entity.User, error)
}

func NewUserRepository(provider *gorm.DB) *userRepository {
	return &userRepository{db: provider}
}

func (r userRepository) CreateUser(ctx context.Context, user entity.User, rawPassword string) (entity.User, error) {
	// Encrypt password
	encryptedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcryptCost)
	if err != nil {
		return entity.User{}, err
	}
	user.EncryptedPassword = encryptedPasswordBytes

	// Generate id
	user.Id = ulid.Make().String()

	// Persist
	err = r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
