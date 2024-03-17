package repository

import (
	"challenge/internal/entity"
	"challenge/pkg/gorm_ext"
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/samber/mo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	gorm_ext.Repository
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{Repository: gorm_ext.NewRepository(db)}
}

func (r UserRepository) CreateUser(ctx context.Context, user entity.User, password string) (entity.User, error) {
	// Set password
	encryptedPassword, err := entity.EncryptPassword(password)
	if err != nil {
		return entity.User{}, err
	}
	user.EncryptedPassword = encryptedPassword

	// Generate id
	user.Id = ulid.Make().String() // TODO: Use uuid v7

	// Create user
	err = r.NewQuery(ctx).Create(&user).Error
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

type UserPatch struct {
	FirstName mo.Option[string]
	LastName  mo.Option[string]
	Nickname  mo.Option[string]
	Email     mo.Option[string]
	Country   mo.Option[string]
	Password  mo.Option[string]

	EncryptedPassword mo.Option[[]byte] // Should not be used
}

func (r UserRepository) PatchUser(ctx context.Context, id string, patch UserPatch) error {
	if p, ok := patch.Password.Get(); ok {
		encryptedPassword, err := entity.EncryptPassword(p)
		if err != nil {
			return err
		}
		patch.EncryptedPassword = mo.Some(encryptedPassword)
	}

	// Update user
	return r.NewQuery(ctx).Where("id", id).Updates(&patch).Error
}

func (r UserRepository) DeleteUser(ctx context.Context, id string) error {
	return r.NewQuery(ctx).Where("id", id).Delete(&entity.User{}).Error
}

func (r UserRepository) ListUsers(ctx context.Context, paginationToken string, pageSize uint, cls ...clause.Expression) ([]entity.User, error) {
	var users []entity.User
	qry := r.NewQuery(ctx, cls...)
	return users, qry.Where("id > ?", paginationToken).Limit(int(pageSize)).Find(&users).Error
}
