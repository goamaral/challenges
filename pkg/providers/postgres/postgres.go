package postgres

import (
	"esl-challenge/pkg/env"
	"fmt"
	"strings"
	"testing"

	"github.com/oklog/ulid/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectToDatabase(dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		env.GetOrDefault("POSTGRES_HOST", "localhost"),
		env.GetOrDefault("POSTGRES_PORT", "5432"),
		env.GetOrDefault("POSTGRES_USER", "postgres"),
		env.GetOrDefault("POSTGRES_PASSWORD", "postgres"),
	)

	if dbName != "" {
		dsn = fmt.Sprintf("%s dbname=%s", dsn, dbName)
	}

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func disconnectFromDatabase(db *gorm.DB) {
	rawDb, err := db.DB()
	if err != nil {
		return
	}
	rawDb.Close()
}

func NewPostgresProvider() (*gorm.DB, error) {
	return connectToDatabase(env.GetOrDefault("POSTGRES_DB", "postgres"))
}

func NewTestPostgresProvider(t *testing.T, databaseInitSql string) (*gorm.DB, func()) {
	db, _ := connectToDatabase("")

	// Create test database
	dbName := strings.ToLower(fmt.Sprintf("test_%s", ulid.Make().String()))
	if err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error; err != nil {
		t.Fatal(err)
	}

	// Connect to test database
	go disconnectFromDatabase(db)
	db, _ = connectToDatabase(dbName)

	// Define close provider function
	closeProvider := func() {
		go disconnectFromDatabase(db)
		db, _ = connectToDatabase("")
		if err := db.Exec(fmt.Sprintf("DROP DATABASE %s", dbName)).Error; err != nil {
			t.Fatal(err)
		}
		go disconnectFromDatabase(db)
	}

	// Load database init
	if err := db.Exec(databaseInitSql).Error; err != nil {
		closeProvider()
		t.Fatal(err)
	}

	return db, closeProvider
}
