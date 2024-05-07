package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
// const dsn = "postgres://rhydberg:fDsR9Kq0GZLXM2Sc37nqFriqBTFCvMuU@dpg-conephq1hbls73fhl92g-a.oregon-postgres.render.com/rhydberg"
const dsn = "postgres://rhydberg:fDsR9Kq0GZLXM2Sc37nqFriqBTFCvMuU@dpg-conephq1hbls73fhl92g-a.oregon-postgres.render.com/rhydberg"
// const devdsn = "postgres://rhydberg:rhydbpass@localhost:5432/chat_app"

func Init() *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	return db
}
	// 