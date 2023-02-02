package repository

import (
	"challenge/internal/entity"
	"challenge/pkg/gormprovider"
	"context"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	gormprovider.AbstractRepository
}

/* PUBLIC */
type UserRepository interface {
	gormprovider.AbstractRepository
	CreateUser(ctx context.Context, user entity.User, password string) (entity.User, error)
	UpdateUser(ctx context.Context, id string, userUpdates entity.User, passwordUpdate string) (entity.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, paginationToken string, pageSize uint, opts ...gormprovider.Option) ([]entity.User, error)
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{AbstractRepository: gormprovider.NewAbstractRepository(db)}
}

type UserFilterOption struct {
	Country *string
}

func (opt UserFilterOption) Apply(qry *gorm.DB) *gorm.DB {
	if opt.Country != nil {
		qry = qry.Where("country", *opt.Country)
	}

	return qry
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
	err = r.NewQuery(ctx).Create(&user).Error
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
	err := r.NewQuery(ctx).Clauses(clause.Returning{}).Where("id", id).Updates(&userUpdates).Error
	if err != nil {
		return entity.User{}, err
	}

	return userUpdates, nil
}

func (r userRepository) DeleteUser(ctx context.Context, id string) error {
	return r.NewQuery(ctx).Where("id", id).Delete(&entity.User{}).Error
}

func (r userRepository) ListUsers(ctx context.Context, paginationToken string, pageSize uint, opts ...gormprovider.Option) ([]entity.User, error) {
	var users []entity.User

	if pageSize == 0 {
		pageSize = 10
	}

	qry := gormprovider.ApplyOptions(r.NewQuery(ctx), opts...)
	err := qry.Where("id > ?", paginationToken).Limit(int(pageSize)).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
