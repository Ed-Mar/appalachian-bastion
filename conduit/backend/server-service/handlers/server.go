package handlers

import (
	"backend/internal"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// KeyServer is a key used for the Server object in the context
type KeyServer struct{}

// Servers handlers for getting and updating servers
type Servers struct {
	severAPILogger *log.Logger
	validator      *internal.Validation
}

// NewServers returns a new servers handlers with the given severAPILogger
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

// getServerID returns the servers ID from the URL
// converts the id into a UUID
func getServerID(r *http.Request) (uuid.UUID, error) {
	// parse the servers id from the url
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])

	// this will catch the any incorrect UUID Input
	if err != nil {
		return id, err
	}

	return id, nil

}
