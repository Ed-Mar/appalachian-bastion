package servers

import (
	"backend/internal"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

// ErrInvalidServerPath is an error message when the server path is not valid
var ErrInvalidServerPath = fmt.Errorf("invalid Path, path should be /servers/[id]")

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"error-message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"error-messages"`
}

// getServerID returns the server ID from the URL
// Panics if it cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func getServerID(r *http.Request) int {
	// parse the server id from the url
	vars := mux.Vars(r)
	//log.Printf("this is the output of the mux.Var(s):%v",		vars)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	//log.Println("id grab from the URI is: %v", id)
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}
