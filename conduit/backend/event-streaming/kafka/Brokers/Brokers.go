package Brokers

import (
	"backend/event-streaming/kafka/config"
	"github.com/segmentio/kafka-go"
	"log"
	"strconv"
)

// GetBroker returns the a go-kafka broker object
func GetBroker() (broker kafka.Broker, err error) {
	// Loads an object used for Kafka other types of connection then converts it to the kafka-go Broker Obj
	kafkaConn, err := config.LoadKafkaConnectionConfig()
	if err != nil {
		return broker, err
	}

	broker.Host = kafkaConn.Host
	temp, err1 := strconv.Atoi(kafkaConn.Port)
	if err1 != nil {
		//Error converting the String value of the loaded env config file struct
		return broker, err1
	}
	broker.Port = temp
	// I do not know what ID nor Racket means really I think I know, but the
	// documentation just leaves think the easy answer

	return broker, nil
}
func GetBrokerString() (hostport string) {
	kafkaConn, err := config.LoadKafkaConnectionConfig()
	if err != nil {
		log.Println(err)
		return ""
	}
	return kafkaConn.HostAndPort

}
