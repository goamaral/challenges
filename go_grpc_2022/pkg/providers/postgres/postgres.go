package postgres

import (
	"challenge/pkg/env"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func connectToDatabase(dbName string, silentLogger bool) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		env.GetOrDefault("POSTGRES_HOST", "localhost"),
		env.GetOrDefault("POSTGRES_PORT", "5432"),
		env.GetOrDefault("POSTGRES_USER", "postgres"),
		env.GetOrDefault("POSTGRES_PASSWORD", "postgres"),
	)

	if dbName != "" {
		dsn = fmt.Sprintf("%s dbname=%s", dsn, dbName)
	}

	cfg := gorm.Config{}
	if silentLogger {
		cfg.Logger = logger.Default.LogMode(logger.Silent)
	}

	return gorm.Open(postgres.Open(dsn), &cfg)
}

func disconnectFromDatabase(db *gorm.DB) {
	rawDb, err := db.DB()
	if err != nil {
		return
	}
	rawDb.Close()
}

func NewPostgresProvider() (db *gorm.DB, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to postgres
	for {
		db, err = connectToDatabase(env.GetOrDefault("POSTGRES_DB", "challenge"), false)
		if err != nil {
			select {
			case <-time.After(time.Second):
				logrus.Info("Waiting for postgres to be ready")
			case <-ctx.Done():
				return nil, err
			}
		} else {
			break
		}
	}

	return db, err
}

func NewTestPostgresProvider(t *testing.T, databaseInitSql string, silentLogger bool) (*gorm.DB, func()) {
	db, _ := connectToDatabase("", silentLogger)

	// Create test database
	dbName := strings.ToLower(fmt.Sprintf("test_%s", ulid.Make().String()))
	if err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error; err != nil {
		t.Fatal(err)
	}

	// Connect to test database
	go disconnectFromDatabase(db)
	db, _ = connectToDatabase(dbName, silentLogger)

	// Define close provider function
	closeProvider := func() {
		go disconnectFromDatabase(db)
		db, _ = connectToDatabase("", silentLogger)
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
