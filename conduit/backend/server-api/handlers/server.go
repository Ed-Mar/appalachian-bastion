package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"server-api/data"
)

// KeyServer is a key used for the Server object in the context
type KeyServer struct{}

// Servers handlers for getting and updating servers
type Servers struct {
	logger    *log.Logger
	validator *data.Validation
}

// NewServers returns a new servers handlers with the given logger
func NewServers(logger *log.Logger, validator *data.Validation) *Servers {
	return &Servers{logger, validator}
}

// ErrInvalidServerPath is an error message when the server path is not valid
var ErrInvalidServerPath = fmt.Errorf("invalid Path, path should be /servers/[id]")

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

// getServerID returns the server ID from the URL
// Panics if it cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func getServerID(r *http.Request) int {
	// parse the server id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}
