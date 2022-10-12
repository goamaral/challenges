package repository

import (
	"context"
	"esl-challenge/internal/entity"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *gorm.DB
}

/* PUBLIC */
type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User, password string) (entity.User, error)
	UpdateUser(ctx context.Context, id string, userUpdates entity.User, passwordUpdate string) (entity.User, error)
	DeleteUser(ctx context.Context, id string) error
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r userRepository) CreateUser(ctx context.Context, user entity.User, password string) (entity.User, error) {
	// Set password
	err := user.SetPassword(password)
	if err != nil {
		return entity.User{}, err
	}

	// Generate id
	user.Id = ulid.Make().String()

	// Create user
	err = r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r userRepository) UpdateUser(ctx context.Context, id string, userUpdates entity.User, passwordUpdate string) (entity.User, error) {
	if passwordUpdate != "" {
		// Set password
		err := userUpdates.SetPassword(passwordUpdate)
		if err != nil {
			return entity.User{}, err
		}
	}

	// Update user
	err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Where("id", id).Updates(&userUpdates).Error
	if err != nil {
		return entity.User{}, err
	}

	return userUpdates, nil
}

func (r userRepository) DeleteUser(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id", id).Delete(&entity.User{}).Error
}
