package di

import (
	"challenge/internal/repository"
	"challenge/pkg/gorm_ext"

	"github.com/samber/do"
)

func Setup() *do.Injector {
	i := do.New()

	/* PROVIDERS */
	do.Provide(i, gorm_ext.NewDB)

	/* REPOSITORIES */
	do.Provide(i, repository.NewUserRepository)

	return i
}
