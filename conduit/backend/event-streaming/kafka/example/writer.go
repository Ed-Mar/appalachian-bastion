package main

import (
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

// just a play ground to mess around with the go-kafka writer
func main() {
	// Make a writer that publishes messages to topic-A.
	// The topic will be created if it is missing.
	w := &kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"),
		Topic:                  "server-service",
		AllowAutoTopicCreation: true,
		BatchTimeout:           5 * time.Millisecond,
		BatchSize:              5,
	}
	tempUUIDv4, err1 := uuid.NewV4()
	if err1 != nil {
		log.Println(err1)
	}
	tempUUIDv1, err2 := uuid.NewV1()
	if err2 != nil {
		log.Println(err2)
	}
	m := make(map[string]string)
	m["serverid"] = "cb44c9aa-bf16-4af8-85ff-ec3b2f748bc4"
	m["value2"] = "serverUUID2"

	var tempEMessage = model.EventMessage{
		SagaOrchestratorName:    "FakeSaga",
		SagaTransactionOriginID: tempUUIDv4,
		SagaTransactionActionID: tempUUIDv1,
		SagaTransactionName:     "Test",
		ServiceTargetName:       "Get - Singleton",
		ServiceTargetOperation:  "GET ME THAT SERVER-Object",
		MessageCreationTime:     time.Now(),
		SagaTransactionData:     m,
	}
	var tempEMessage1 = model.EventMessage{

		SagaOrchestratorName:    "FakeSaga1",
		SagaTransactionOriginID: tempUUIDv4,
		SagaTransactionActionID: tempUUIDv1,
		SagaTransactionName:     "Test01",
		ServiceTargetName:       "Get - Collection",
		ServiceTargetOperation:  "GET ME THAT SERVER-Object",
		MessageCreationTime:     time.Now(),
		SagaTransactionData:     m,
	}
	laData, err3 := json.Marshal(tempEMessage)
	laData1, err3 := json.Marshal(tempEMessage1)

	/* check the error */
	if err3 != nil {
		fmt.Println("There was an error :[")
		return
	}
	log.Println(" Event Message was created")
	messages := []kafka.Message{
		{
			Key:   []byte("Key-A"),
			Value: laData,
		},
		{
			Key:   []byte("Key-B"),
			Value: laData1,
		},
	}
	var err error
	const retries = 1
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// attempt to create topic prior to publishing the message
		err = w.WriteMessages(ctx, messages...)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			log.Fatalf("unexpected error %v", err)
		}
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
