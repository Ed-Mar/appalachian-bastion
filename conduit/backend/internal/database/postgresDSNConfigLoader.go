package database

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// GetDBPostgresDSN returns a dns to be used for connecting to a database when give the location on the config file
// envFileDir exmaple "../foo/bar" envFileName: "postgres" filetype: "env"
func GetDBPostgresDSN(envFileDir string, envFileName string, fileType string) (string, error) {
	dbConfig, err := loadPostgresConfig(envFileDir, envFileName, fileType)
	if err != nil {
		return "", err
	}

	//TODO Fix this hard coded injection the default on the mapping is not working as expected
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
	return dsn, nil
}

type dbConnectionConfigPostgres struct {
	//Net      string "tcp"
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     string `mapstructure:"POSTGRES_PORT"`
	Username string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	DBName   string `mapstructure:"POSTGRES_DB"`
	SSLMode  string `default:"disable"`
	Timezone string `mapstructure:"PGTZ"`
}

// TODO Come back and make this only need the file path with file in it, having come in as separate seems annoying to use
func loadPostgresConfig(fileDir string, fileName string, fileType string) (config dbConnectionConfigPostgres, err error) {
	viper.AddConfigPath(fileDir)
	viper.SetConfigName(fileName)
	viper.SetConfigType(fileType)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("[ENV FILE ERROR]: ", err)
		return
	}

	err = viper.Unmarshal(&config)
	return
}
