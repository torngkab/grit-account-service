package model

import (
	"time"

	"github.com/google/uuid"
)

type Balance struct {
	AccountId       uuid.UUID `gorm:"unique;type:uuid"`
	Balance         float64
	LatestUpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (b *Balance) TableName() string {
	return "balances"
}
