package model

import (
	"github.com/gofrs/uuid"
	"time"
)

// EventMessage a general message structure for kafka messages
type EventMessage struct {
	// Planning on assigning each of the Saga Orchestrators an ID, not sure how I am going to keep track of each
	// of them.
	SagaOrchestratorName string `json:"SagaOrchestratorName"`
	// In the case that of a transaction triggering other event-messages to other services this will stay to allow for the
	// propagates back to it to finish this transaction. Not that positive if this will be needed, but maybe
	SagaTransactionOriginID uuid.UUID `json:"SagaTransactionOriginID"`

	FromServiceTopic   string `json:"FromServiceTopic"`
	TargetServiceTopic string `json:"TargetServiceTopic"`

	// This is the ID of the Transaction sending the request
	SagaTransactionActionID uuid.UUID `json:"SagaTransactionID"`
	// This is the ID of the Transaction sending the request
	SagaTransactionName string `json:"SagaTransactionName"`
	// Which Service is this message intended for
	ServiceTargetName string `json:"ServiceTargetName"`
	// Which Service Operation is message intended for
	ServiceTargetOperation string `json:"ServiceTargetOperation"`
	//Putting the time of the creation of this message in the obj so that  I get better understanding the delay
	MessageCreationTime time.Time `json:"MessageCreationTime"`
	//Is a Response Message Expected
	ReturnMessageExpected bool `json:"ReturnMessageExpected"`
	// For Response Messages to Understand if the request operation was completed. Think I might go with the http codes 200,201,204
	// and error
	ResponseCode string `json:"ResponseCode"`
	// Map to send necessary values between services and or error information to be handled by the Reader.
	SagaTransactionData map[string]string `json:"SagaTransactionData"`
}
