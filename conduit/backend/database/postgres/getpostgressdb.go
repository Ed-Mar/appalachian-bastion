package postgres

import (
	"backend/database/postgres/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func GetPostgresDB() (db *gorm.DB) {
	dsn := config.DbPostgresURL()
	log.Printf("dsn string: %s\n", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database | Error: %s", err))
	}
	return db

}
