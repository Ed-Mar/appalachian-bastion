package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"log"
)

// Note that this is relevant to the where this is being run
const kakfaConfigPath = "../config/"
const kakfaConfigFilename = "kafkaConfig"
const kakfaconfigFileType = "env"

type kafkaConnectionConfig struct {
	ConnectionType string `mapstructure:"CONNECTION_TYPE"`
	Host           string `mapstructure:"KAFKA_HOST"`
	Port           string `mapstructure:"KAFKA_PORT"`
	HostAndPort    string
}

func loadKafkaConfig(fileDir string, fileName string, fileType string) (config kafkaConnectionConfig, err error) {
	viper.AddConfigPath(fileDir)
	viper.SetConfigName(fileName)
	viper.SetConfigType(fileType)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("[ENV FILE ERROR] [Kafka]: Error Loading config env file", err)
		return
	}

	err = viper.Unmarshal(&config)
	config.HostAndPort = config.Host + ":" + config.Port

	//TODO REMOVE THIS LOG ONCE confirmed it working as expected
	fmt.Printf("%+v\n", config)

	return
}

// GetKafkaConnection A wrapper on a wrapper with the loaded connection information for kafka from local file.
func GetKafkaConnection(context context.Context, topic string, partition int) (conn *kafka.Conn, err error) {
	kafkaConfig, err := loadKafkaConfig(kakfaConfigPath, kakfaConfigFilename, kakfaconfigFileType)
	if err != nil {
		return nil, err
	}
	conn, err = kafka.DialLeader(context, kafkaConfig.ConnectionType, kafkaConfig.HostAndPort, topic, partition)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
