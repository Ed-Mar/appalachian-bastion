package helper

import (
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// GetURIParmWithMatchingName returns the ID from the URL
// Get the Parameter form the URI that has the matching name
func GetURIParmWithMatchingName(r *http.Request, idName string) (uuid.UUID, error) {
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

func GetNumOfURIParms(r *http.Request) int {
	return len(mux.Vars(r))
}
