package gorm_ext

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectToDatabase(dialector gorm.Dialector, silentLogger bool) (db *gorm.DB, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg := gorm.Config{}
	if silentLogger {
		cfg.Logger = logger.Default.LogMode(logger.Silent)
	}

	for {
		db, err = gorm.Open(dialector, &cfg)
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

func DisconnectFromDatabase(db *gorm.DB) {
	rawDb, err := db.DB()
	if err != nil {
		return
	}
	rawDb.Close()
}
