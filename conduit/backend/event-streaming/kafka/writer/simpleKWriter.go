package main

import (
	"backend/event-streaming/kafka/messages/model"
	"github.com/gofrs/uuid"
	"log"
	"time"
)

// This was just test use of the generic Writer
func main() {
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

	var tempMessage0 = model.EventMessage{
		SagaOrchestratorName:    "FakeSaga",
		SagaTransactionOriginID: tempUUIDv4,
		SagaTransactionActionID: tempUUIDv1,
		SagaTransactionName:     "Test",
		ServiceTargetName:       "Get - Singleton",
		ServiceTargetOperation:  "GET ME THAT SERVER-Object",
		MessageCreationTime:     time.Now(),
		SagaTransactionData:     m,
		TargetServiceTopic:      "server-service",
	}
	var tempMessage1 = model.EventMessage{

		SagaOrchestratorName:    "FakeSaga1",
		SagaTransactionOriginID: tempUUIDv4,
		SagaTransactionActionID: tempUUIDv1,
		SagaTransactionName:     "Test01",
		ServiceTargetName:       "Get - Collection",
		ServiceTargetOperation:  "GET ME THAT SERVER-Object",
		MessageCreationTime:     time.Now(),
		SagaTransactionData:     m,
		TargetServiceTopic:      "server-service",
	}
	var testMessages []model.EventMessage
	testMessages = append(testMessages, tempMessage0, tempMessage1)
	err := GenericKafkaWriter(testMessages)
	if err != nil {
		log.Println(err)
	}
}
