package handlers

import (
	"backend/internal"
	"fmt"
	"log"
)

// KeyServer is a key used for the Server object in the context
type KeyServer struct{}

// Servers handlers for getting and updating servers
type Servers struct {
	APILogger *log.Logger
	validator *internal.Validation
}

// NewServers returns a new servers handlers with the given APILogger
func NewServers(severAPILogger *log.Logger, validator *internal.Validation) *Servers {
	return &Servers{severAPILogger, validator}
}

// ErrInvalidServerPath is an error message when the servers path is not valid
var ErrInvalidServerPath = fmt.Errorf("invalid Path, path should be /servers/[id]")

// GenericError is a generic error message returned by a servers
type GenericError struct {
	Message string `json:"error-message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"error-messages"`
}
