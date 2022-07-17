package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"log"
)

// Note that this is relevant to the where this is being run
const kakfaConfigPath = "."
const kakfaAbsConfigPath = "/home/kaiser/workspace/tech/appalachian-bastion/conduit/backend/event-streaming/kafka/config/"
const kakfaConfigFilename = "kafkaConfig"
const kakfaconfigFileType = "env"

var ErrLoadingEnvFile = errors.New("[ERROR] [KAFKA] [ENV]: Error Loading ENV")

type kafkaConnectionConfig struct {
	ConnectionType string `mapstructure:"CONNECTION_TYPE"`
	Host           string `mapstructure:"KAFKA_HOST"`
	Port           string `mapstructure:"KAFKA_PORT"`
	HostAndPort    string
}

func LoadKafkaConnectionConfig() (config kafkaConnectionConfig, err error) {
	viper.AddConfigPath(kakfaConfigPath)
	viper.AddConfigPath(kakfaAbsConfigPath)
	viper.SetConfigName(kakfaConfigFilename)
	viper.SetConfigType(kakfaconfigFileType)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Printf(ErrLoadingEnvFile.Error(), err)
	}
	err = viper.Unmarshal(&config)
	config.HostAndPort = config.Host + ":" + config.Port

	//TODO REMOVE THIS LOG ONCE confirmed it working as expected
	fmt.Printf("%+v\n", config)

	return

}

// GetKafkaConnection A wrapper on a wrapper with the loaded connection information for kafka from local file.
func GetKafkaConnection(context context.Context, topic string, partition int) (conn *kafka.Conn, err error) {
	kafkaConfig, err := LoadKafkaConnectionConfig()
	if err != nil {
		return nil, err
	}
	conn, err = kafka.DialLeader(context, kafkaConfig.ConnectionType, kafkaConfig.HostAndPort, topic, partition)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
func GetKafkaHostAndPort() (string, error) {
	kafkaConfig, err := LoadKafkaConnectionConfig()
	if err != nil {
		return "", err
	} else {
		return kafkaConfig.HostAndPort, nil
	}
}
