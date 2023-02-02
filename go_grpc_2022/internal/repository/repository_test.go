package repository_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gorm.io/gorm"

	"challenge/pkg/gormprovider"
)

func testInit(t *testing.T) (*gorm.DB, func()) {
	deploymentFolderPath, _ := filepath.Abs("../../deployment")
	databaseInitSqlBytes, err := os.ReadFile(fmt.Sprintf("%s/database_init.sql", deploymentFolderPath))
	if err != nil {
		t.Fatal(err)
	}

	return gormprovider.NewTestPostgresProvider(t, string(databaseInitSqlBytes), true)
}
