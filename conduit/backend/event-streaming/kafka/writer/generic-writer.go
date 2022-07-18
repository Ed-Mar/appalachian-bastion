package main

import (
	"backend/event-streaming/kafka/config"
	"backend/event-streaming/kafka/messages/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

var errGenericEventMessage = errors.New("[ERROR] [EVENT-STREAMING] [WRITER] ")
var errMarshallingEventMessage = errors.New("[ERROR] [EVENT-STREAMING] [WRITER] [JSON] [MARSHALLING]: ")

func GenericKafkaWriter(messages []model.EventMessage) error {
	if len(messages) <= 0 {
		ErrGenericEventMessage := fmt.Errorf("[INPUT]: Don't send and empty slice")
		log.Println(ErrGenericEventMessage)
	}
	hostAndPort, err := config.GetKafkaHostAndPort()
	if err != nil {
		// Error loading config file. Logged at place of origin
		return err
	}

	//Note that the topic is given, so it must be given in the messages itself
	kWriter := &kafka.Writer{
		Addr:                   kafka.TCP(hostAndPort),
		AllowAutoTopicCreation: true,
		//Shorten Time due to speed concerns plus expected throughput is low
		BatchTimeout: 260 * time.Millisecond,
		BatchSize:    5,
	}

	var kafkaMessages []kafka.Message

	//iterate over all passage messaged to sent via kafka.
	for count, _ := range messages {
		// Creates and unique id for the key of the kafka message to ensure that order is persevered
		key, err1 := uuid.NewV4()
		if err1 != nil {
			ErrGenericEventMessage := fmt.Errorf("%w [UUID]: Error creating UUID for kafka message key | ", err1)
			log.Println(ErrGenericEventMessage)
			return ErrGenericEventMessage
		}
		//convert the single message at slice count into a byte slice for passing onto kafka.
		marshalledMessage, err := json.Marshal(messages[count])
		if err != nil {
			ErrMarshallingEventMessage := fmt.Errorf("%w | ", err)
			return ErrMarshallingEventMessage
		}

		kafkaMessages = append(kafkaMessages, kafka.Message{
			Topic:         messages[count].TargetServiceTopic,
			Partition:     0,
			Offset:        0,
			HighWaterMark: 0,
			Key:           key.Bytes(),
			Value:         marshalledMessage,
			Headers:       nil,
			Time:          time.Time{},
		},
		)
	}

	const retries = 1
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for i := 0; i < retries; i++ {
		// attempt to create topic prior to publishing the message
		err = kWriter.WriteMessages(ctx, kafkaMessages...)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}
		if err != nil {
			log.Fatalf("unexpected error %v", err)
		}
	}

	return nil
}
