package consumer_group

import (
	"backend/event-streaming/kafka/Brokers"
	consumer_groups "backend/event-streaming/kafka/consumer-groups"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
)

//ServerServiceConsumerGroup Need to check if the broker is empty before use, but think that error should be handled further down the stack
var ServerServiceConsumerGroup = kafka.NewReader(kafka.ReaderConfig{
	Brokers:  []string{Brokers.GetBrokerString()},
	GroupID:  "server-service-group",
	Topic:    "server-service",
	MinBytes: 1e3,  // 1KB
	MaxBytes: 10e3, // 10KB
})

const (
	groupID  = "server-service-group"
	topic    = "server-service"
	minBytes = 1e3  // 1KB
	maxBytes = 10e3 // 10KB
)

func GetServerServiceConsumerGroup() (kafka.Reader, error) {
	ServerServiceConsumerGroup1, err := consumer_groups.CreateConsumerGroupReader(groupID, topic, minBytes, maxBytes)
	if err != nil {
		var errServerServiceConsumerGroupCreation = errors.New("[ERROR] [EVENT] [READER]: ")
		return ServerServiceConsumerGroup1, fmt.Errorf("%w Error Creating Server Group", errServerServiceConsumerGroupCreation)
	}
	return ServerServiceConsumerGroup1, err
}
