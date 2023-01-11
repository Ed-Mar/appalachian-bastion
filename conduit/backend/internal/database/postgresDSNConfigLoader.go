package database

import (
	"fmt"
	"github.com/spf13/viper"
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
func GetDBPostgresDSNv2(filename string, possibleFileLocations []string) (string, error) {
	dbConfig, err := loadPostrgresConfigENVv2(filename, possibleFileLocations)
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

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
		}
	}

	err = viper.Unmarshal(&config)
	return
}
func loadPostrgresConfigENVv2(filename string, possibleFileLocations []string) (config dbConnectionConfigPostgres, err error) {
	viper.SetConfigName(filename)
	viper.SetConfigType("env")
	// Since you can add multiple paths for viper to check, and I keep running in issue with relative file locations I am just gonna throw the kitchen sink at it.
	for _, path := range possibleFileLocations {
		viper.AddConfigPath(path)
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return config, err
		} else {
			// Config file was found but another error was produced
			return config, err
		}
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}
	return config, err
}
