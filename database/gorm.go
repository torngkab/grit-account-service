package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(dbUrl string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	return db
}
