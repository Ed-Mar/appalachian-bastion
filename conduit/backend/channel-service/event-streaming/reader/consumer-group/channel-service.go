package consumer_group

import (
	"backend/event-streaming/kafka/Brokers"
	"github.com/segmentio/kafka-go"
)

//ChannelServiceConsumerGroup Need to check if the broker is empty before use, but think that error should be handled further down the stack
var ChannelServiceConsumerGroup = kafka.NewReader(kafka.ReaderConfig{
	Brokers:  []string{Brokers.GetBrokerString()},
	GroupID:  "channel-service-group",
	Topic:    "channel-service",
	MinBytes: 1e3,  // 1KB
	MaxBytes: 10e3, // 10KB
})
