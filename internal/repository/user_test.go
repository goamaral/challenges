package repository_test

import (
	"context"
	"esl-challenge/internal/entity"
	"esl-challenge/internal/repository"
	"esl-challenge/pkg/providers/postgres"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserRepository_CreateUser(t *testing.T) {
	deploymentFolderPath, _ := filepath.Abs("../../deployment")
	databaseInitSqlBytes, err := os.ReadFile(fmt.Sprintf("%s/database_init.sql", deploymentFolderPath))
	if err != nil {
		t.Fatal(err)
	}

	provider, closeProvider := postgres.NewTestPostgresProvider(t, string(databaseInitSqlBytes))
	defer closeProvider()

	rawPassword := "password"
	newUser := entity.User{}

	r := repository.NewUserRepository(provider)

	// Success
	user, err := r.CreateUser(context.Background(), newUser, rawPassword)
	if assert.NoError(t, err) {
		assert.NotZero(t, user.Id)
		assert.NotZero(t, user.CreatedAt)
		assert.NotZero(t, user.UpdatedAt)
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(rawPassword)))
	}

	// Failure - Nickname not unique
	_, err = r.CreateUser(context.Background(), newUser, rawPassword)
	if assert.Error(t, err) {
		pgErr, ok := err.(interface {
			SQLState() string
		})
		assert.True(t, ok, "Not pg error")
		assert.Equal(t, "23505", pgErr.SQLState()) // unique_violation
	}

	// Failure - Email not unique
	newUser.Nickname = "another"
	_, err = r.CreateUser(context.Background(), newUser, rawPassword)
	if assert.Error(t, err) {
		pgErr, ok := err.(interface {
			SQLState() string
		})
		assert.True(t, ok, "Not pg error")
		assert.Equal(t, "23505", pgErr.SQLState()) // unique_violation
	}
}
