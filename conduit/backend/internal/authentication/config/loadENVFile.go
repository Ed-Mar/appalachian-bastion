package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

var ErrENVFile = fmt.Errorf("[ERROR][ENV-FILE] ")

func GetPublicRSASecret(envFileDir string, envFileName string, fileType string) (string, error) {
	config, err := loadPublicRSASecret(envFileDir, envFileName, fileType)
	if err != nil {
		err = fmt.Errorf(" %v | %v\n", ErrENVFile, err)
		return "", err
	}

	return config.PublicRSASecret, err
}

type authPublicRSASecret struct {
	PublicRSASecret string `mapstructure:"PUBLIC_RSA_SECRET"`
}

// TODO Come back and make this only need the file path with file in it, having come in as separate seems annoying to use
// TODO add some fallback options to this cause this is going to be run in every public microservices
func loadPublicRSASecret(fileDir string, fileName string, fileType string) (config authPublicRSASecret, err error) {
	viper.AddConfigPath(fileDir)
	viper.SetConfigName(fileName)
	viper.SetConfigType(fileType)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		err = fmt.Errorf(" %v | %v\n", ErrENVFile, err)
		log.Println(err)
		return
	}

	err = viper.Unmarshal(&config)
	return
}
