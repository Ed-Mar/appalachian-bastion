package postgres

import (
	"backend/database/postgres/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetPostgresDB() (db *gorm.DB) {
	dsn := config.DbPostgresURL()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database | Error: %s", err))
	}
	return db

}
