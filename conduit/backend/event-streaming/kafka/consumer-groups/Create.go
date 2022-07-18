package consumer_groups

import (
	"backend/event-streaming/kafka/Brokers"
	"backend/event-streaming/kafka/config"
	"github.com/segmentio/kafka-go"
)

// CreateConsumerGroup create
func CreateConsumerGroupReader(groupID string, topicName string, minBytes int, maxBytes int) (kafka.Reader, error) {
	//Just trying to get ahead of the config not loading, but I think I have this error covered
	if Brokers.GetBrokerString() != "" {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{Brokers.GetBrokerString()},
			GroupID:  groupID,
			Topic:    topicName,
			MinBytes: minBytes,
			MaxBytes: maxBytes,
		})
		return *reader, nil
	}
	return kafka.Reader{}, config.ErrLoadingEnvFile
}
