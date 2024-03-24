package gorm_ext

import (
	"github.com/samber/do"
	"gorm.io/gorm"
)

func NewDB(i *do.Injector) (*DB, error) {
	// TODO
	return nil, nil
}

type DB struct {
	*gorm.DB
}

func (db *DB) HealthCheck() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (db *DB) Shutdown() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
