package internal

import (
	"backend/database/postgres"
	"backend/server-api/data"
	"fmt"
)

func AutoMigrateDB() {
	db := postgres.GetPostgresDB()
	err := db.AutoMigrate(&data.Server{})
	if err != nil {
		panic(fmt.Sprintf("failed to Automigrate Table | Error: %s", err))
	}
}
