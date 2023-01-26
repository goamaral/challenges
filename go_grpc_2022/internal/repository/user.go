package repository

import (
	"challenge/internal/entity"
	"context"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	abstractRepository
}

/* PUBLIC */
type UserRepository interface {
	AbstractRepository
	CreateUser(ctx context.Context, user entity.User, password string) (entity.User, error)
	UpdateUser(ctx context.Context, id string, userUpdates entity.User, passwordUpdate string) (entity.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, paginationToken string, pageSize uint, opts *ListUsersOpts) ([]entity.User, error)
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{abstractRepository: abstractRepository{db}}
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
	err = r.newQuery(ctx).Create(&user).Error
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
	err := r.newQuery(ctx).Clauses(clause.Returning{}).Where("id", id).Updates(&userUpdates).Error
	if err != nil {
		return entity.User{}, err
	}

	return userUpdates, nil
}

func (r userRepository) DeleteUser(ctx context.Context, id string) error {
	return r.newQuery(ctx).Where("id", id).Delete(&entity.User{}).Error
}

type ListUsersOpts struct {
	Country string
}

func (opts *ListUsersOpts) Apply(qry *gorm.DB) *gorm.DB {
	if opts != nil {
		if opts.Country != "" {
			qry = qry.Where("country", opts.Country)
		}
	}

	return qry
}

func (r userRepository) ListUsers(ctx context.Context, paginationToken string, pageSize uint, opts *ListUsersOpts) ([]entity.User, error) {
	var users []entity.User

	if pageSize == 0 {
		pageSize = 10
	}

	qry := r.newQuery(ctx).Where("id > ?", paginationToken).Limit(int(pageSize))
	err := opts.Apply(qry).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
