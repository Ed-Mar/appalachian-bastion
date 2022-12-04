package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

//TODO Need to figure out if I should base64 this prior(so base64 it myself then paste that in the file) or have this do it

// GetTokenIntrospectiveClientBasicInfo returns the header value with "Basic Ze........" Ready for header injection.
func GetTokenIntrospectiveClientBasicInfo(envFileDir string, envFileName string, fileType string) (*string, error) {
	clientAuthentication := ""
	config, err := loadTokenIntrospectiveClientInfo(envFileDir, envFileName, fileType)
	if err != nil {
		err = fmt.Errorf(" %v | %v\n", ErrENVFile, err)
		return &clientAuthentication, err
	}
	clientAuthentication = fmt.Sprintf(
		"Basic %s",
		config.Authentication,
	)

	return &clientAuthentication, err
}

type tokenIntrospectiveClientInfo struct {
	//Currently I am only using the basic authentication for the client which is the base64(<CLIENT>:<PASSWORD>)
	Authentication string `mapstructure:"AUTHORIZATION"`
}

// TODO I need to make this to return/intake a genertic so I don't have to differnt code just to load env vars from file
// TODO Come back and make this only need the file path with file in it, having come in as separate seems annoying to use
// TODO add some fallback options to this cause this is going to be run in every public microservices
func loadTokenIntrospectiveClientInfo(fileDir string, fileName string, fileType string) (config tokenIntrospectiveClientInfo, err error) {
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
