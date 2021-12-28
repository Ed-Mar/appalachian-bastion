package channels

import (
	"backend/internal"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type KeyChannel struct{}

// Channels handlers for getting and updating Channels
type Channels struct {
	channelAPILogger *log.Logger
	validator        *internal.Validation
}

// NewChannels returns a new Channels handlers with the given channelAPILogger
func NewChannels(channelAPILogger *log.Logger, validator *internal.Validation) *Channels {
	return &Channels{channelAPILogger, validator}
}

// ErrInvalidChannelPath is an error message when the ID path is not valid
var ErrInvalidChannelPath = fmt.Errorf("invalid Path, path should be /servers/[id]/channels/[id]")

// GenericError is a generic error message returned by a ID
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

// getChannelID returns the ID ID from the URL
// Panics if it cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func getChannelID(r *http.Request) int {
	// parse the ID id from the url
	vars := mux.Vars(r)
	log.Printf("this is the output of the mux.Var(s):%v",
		vars)
	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}
