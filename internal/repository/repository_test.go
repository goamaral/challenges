package repository_test

import (
	"esl-challenge/pkg/providers/postgres"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func testInit(t *testing.T) (*gorm.DB, func()) {
	deploymentFolderPath, _ := filepath.Abs("../../deployment")
	databaseInitSqlBytes, err := os.ReadFile(fmt.Sprintf("%s/database_init.sql", deploymentFolderPath))
	if err != nil {
		t.Fatal(err)
	}

	return postgres.NewTestPostgresProvider(t, string(databaseInitSqlBytes))
}

func assertUniqueViolationError(t *testing.T, err error) {
	if assert.Error(t, err) {
		pgErr, ok := err.(interface {
			SQLState() string
		})
		assert.True(t, ok, "Not pg error")
		assert.Equal(t, "23505", pgErr.SQLState()) // unique_violation
	}
}
