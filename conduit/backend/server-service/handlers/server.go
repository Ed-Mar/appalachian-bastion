package handlers

import (
	"backend/internal"
	"fmt"
	"log"
)

// KeyServer is a key used for the Server object in the context
type KeyServer struct{}

// ServerHandler gerneric-handlers for getting and updating servers
type ServerHandler struct {
	APILogger *log.Logger
	validator *internal.Validation
}

// NewServersHandler returns a new servers generic-handlers with the given APILogger
func NewServersHandler(severAPILogger *log.Logger, validator *internal.Validation) *ServerHandler {
	return &ServerHandler{severAPILogger, validator}
}

// ErrInvalidServerPath is an error message when the servers path is not valid
var ErrInvalidServerPath = fmt.Errorf("invalid Path, path should be /servers/[id]")

// GenericError is a generic error message returned by a ID
type GenericError struct {
	Message string `json:"error-message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"error-messages"`
}
