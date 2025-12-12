package model

import (
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	AccountTypeReferral     AccountType = "1001"
	AccountTypeMain         AccountType = "1002"
	AccountTypeDisbursement AccountType = "1003"
	AccountTypePSP          AccountType = "1004" // Payment Service Provider
)

type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "active"
	AccountStatusInactive AccountStatus = "inactive"
)

type Account struct {
	Id          uuid.UUID     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserId      uuid.UUID     `gorm:"foreignKey:Id;references:UserId;type:uuid;uniqueIndex:idx_account_type_user_id;index:idx_account_type_status"`
	AccountType AccountType   `gorm:"uniqueIndex:idx_account_type_user_id;index:idx_account_type_status"`
	Status      AccountStatus `gorm:"index:idx_account_type_status"`
	CreatedAt   time.Time     `gorm:"autoCreateTime"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime"`
}

func (a *Account) TableName() string {
	return "accounts"
}
