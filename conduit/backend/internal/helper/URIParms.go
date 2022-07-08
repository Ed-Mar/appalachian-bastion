package helper

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var ErrParmFetchingInternalError = fmt.Errorf("[ERROR] [INTERNAL] internal Error Fectching Parms from Request as Occured and has been logged")
var ErrIncorrectUUIDFormat = fmt.Errorf("UUID passed in the request is in inccorrect format")

// GetUUIDFromReqParm returns the ID from the URL
// Get the Parameter form the URI that has the matching name
func GetUUIDFromReqParm(r *http.Request, idName string) (uuid.UUID, error) {
	// parse the ids from the uri
	vars := mux.Vars(r)

	log.Println("Number of Variables in the Request", len(mux.Vars(r)))
	fmt.Println(vars)

	//Check to see if the request parm matches the one we are looking for
	_, ok := vars[idName]
	if !ok {
		//TODO find a way to send a nil or something for the uuid cause it I can't send a nil I think It needs a pointer I'll look at it later
		temp, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")
		return temp, ErrParmFetchingInternalError
	}
	//converting it from the string should catch any error in relation to the validation
	id, err := uuid.FromString(vars[idName])
	if err != nil {
		ErrIncorrectUUIDFormat = fmt.Errorf("%v | %v", ErrIncorrectUUIDFormat, err)
		return id, ErrIncorrectUUIDFormat
	}

	return id, nil
}

func GetNumOfURIParms(r *http.Request) int {
	return len(mux.Vars(r))
}
