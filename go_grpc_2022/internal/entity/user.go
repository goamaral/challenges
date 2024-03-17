package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id                string `gorm:"primaryKey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	FirstName         string
	LastName          string
	Nickname          string
	EncryptedPassword []byte
	Email             string
	Country           string
}

const bcryptCost = 11

func EncryptPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
}
