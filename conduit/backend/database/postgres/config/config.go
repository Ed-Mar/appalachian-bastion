package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

var envFilePath = "../database/postgres/config/"

// DBConfPostgres struct is to pull the config from the env from the local postgres.env file
type dbConfigPostgres struct {
	//Net      string "tcp"

	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     string `mapstructure:"POSTGRES_PORT"`
	Username string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	DBName   string `mapstructure:"POSTGRES_DB"`
	SSLMode  string `default:"disable"`
	Timezone string `mapstructure:"PGTZ"`
}

// DbPostgresURL will return the dBConfPostgres to open a connection with gorm to the postgres conduit database
func DbPostgresURL() string {
	dbConfig, err := loadPostgresConfig(envFilePath)
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	//	"host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	log.Printf(dbConfig.SSLMode)
	//TEMP FIX
	sslmode := "disable"
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s Timezone=%s",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.DBName,
		sslmode,
		dbConfig.Timezone,
	)
	log.Printf(dsn)
	return dsn
}
func loadPostgresConfig(path string) (config dbConfigPostgres, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("postgres")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to load Posgres env %s", err)
		return
	}

	err = viper.Unmarshal(&config)
	return
}
