package users

import (
	"backend/internal"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// KeyUser is a key used for the User object in the context
type KeyUser struct{}

// Users handlers for getting and updating Users
type Users struct {
	userAPILogger *log.Logger
	validator     *internal.Validation
}

// NewUsers returns a new Users handlers with the given
func NewUsers(userAPILogger *log.Logger, validator *internal.Validation) *Users {
	return &Users{userAPILogger, validator}
}

// GenericError is a generic error message returned by a user
type GenericError struct {
	Message string `json:"error-message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"error-messages"`
}

func getUserID(r *http.Request) int {
	// parse the servers id from the url
	vars := mux.Vars(r)
	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	//log.Println("id grab from the URI is: %v", id)
	if err != nil {
		// should never happen
		panic(err)
	}
	return id
}
