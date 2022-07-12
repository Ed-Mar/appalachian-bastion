package consumer_groups

import (
	"backend/event-streaming/kafka/Brokers"
	"github.com/segmentio/kafka-go"
)

//ServiceEventBus Need to check if the broker is empty
var ServiceEventBus = kafka.NewReader(kafka.ReaderConfig{
	Brokers:  []string{Brokers.GetBrokerString()},
	GroupID:  "service-event-bus-group",
	Topic:    "service-event-bus",
	MinBytes: 1e3,  // 1KB
	MaxBytes: 10e3, // 10KB
})
