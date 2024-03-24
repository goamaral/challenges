package helper

import (
	"challenge/internal/di"
	"challenge/internal/entity"
	"challenge/pkg/gorm_ext"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/samber/do"
	"github.com/stretchr/testify/require"
)

func RunTest(t *testing.T, name string, fn func(Tester)) {
	t.Run(name, func(t *testing.T) {
		fn(newTester(t))
	})
}

type Tester struct {
	T *testing.T
	I *do.Injector
}

func newTester(t *testing.T) Tester {
	i := di.Setup()
	t.Cleanup(func() { require.NoError(t, i.Shutdown()) })

	return Tester{T: t, I: i}
}

/* DB */
func (t *Tester) AddUser(user entity.User, password string) entity.User {
	if password == "" {
		password = "password"
	}
	encryptedPassword, err := entity.EncryptPassword(password)
	require.NoError(t.T, err)
	user.EncryptedPassword = encryptedPassword

	if user.Id == "" {
		user.Id = ulid.Make().String()
	}

	db := do.MustInvoke[*gorm_ext.DB](t.I)
	require.NoError(t.T, db.Create(&user).Error)

	return user
}
