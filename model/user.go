package model

import (
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)

type User struct {
	Id        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Status    UserStatus
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (u *User) TableName() string {
	return "users"
}
