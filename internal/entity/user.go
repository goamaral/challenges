package entity

import (
	"time"
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
