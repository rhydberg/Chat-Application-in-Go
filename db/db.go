package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


const dsn = "postgres://rhydberg:rhydbpass@localhost:5432/chat_app"

func Init() *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	return db
}
	