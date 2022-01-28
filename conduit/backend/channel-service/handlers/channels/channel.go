package channels

import (
	"backend/internal"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type KeyChannel struct{}

// Channels handlers for getting and updating Channels
type Channels struct {
	APILogger *log.Logger
	validator *internal.Validation
}

// NewChannels returns a new Channels handlers with the given APILogger
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

// getURIParmWithMatchingName returns the ID from the URL
// Get the Parameter form the URI that has the matching name
func getURIParmWithMatchingName(r *http.Request, idName string) (uuid.UUID, error) {
	// parse the ids from the uri
	vars := mux.Vars(r)
	log.Println("Here is the Length of the Vars from the mux", len(mux.Vars(r)))
	id, err := uuid.FromString(vars[idName])
	log.Println("From Mux: ", id)
	// this will catch the any incorrect UUID Input
	if err != nil {
		return id, err
	}

	return id, nil
}

func getNumOfURIParms(r *http.Request) int {
	return len(mux.Vars(r))
}
