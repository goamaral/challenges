package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 11

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

func (u *User) SetPassword(password string) error {
	encryptedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	u.EncryptedPassword = encryptedPasswordBytes

	return nil
}
