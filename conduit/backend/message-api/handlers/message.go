package handlers

import (
	"backend/internal"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type KeyMessage struct{}

// Messages handlers for getting and updating Messages
type Messages struct {
	messageAPILogger *log.Logger
	validator        *internal.Validation
}

// NewMessages returns a new Messages handlers with the given MessageAPILogger
func NewMessages(messageAPILogger *log.Logger, validator *internal.Validation) *Messages {
	return &Messages{messageAPILogger, validator}
}

// ErrInvalidMessagePath is an error message when the ID path is not valid
var ErrInvalidMessagePath = fmt.Errorf("invalid Path, path should be /servers/[id]/channels/[id]")

// GenericError is a generic error message returned by a ID
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

// getMessageID returns the Message ID from the URL
// Panics if it cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func getMessageID(r *http.Request) int {
	// parse the ID id from the url
	vars := mux.Vars(r)
	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}
