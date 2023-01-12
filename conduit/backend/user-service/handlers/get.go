package handlers

import (
	"backend/internal"
	"backend/internal/authentication/handlers"
	"backend/user-service/models"
	"fmt"
	"net/http"
)

var ErrUserServiceGetBase = fmt.Errorf("%v [GET] ", ErrUserServiceBaseError)
var WarnUserServiceGetBase = fmt.Errorf("%v [GET] ", WarnUserServiceBaseWarning)

func (uh *UserHandler) GetUserViaExternalId(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	sub, err := handlers.AuthGetExternalIDFromContext(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	uh.StandardHandler.ServiceLogger.Printf("Attempting to get User Information via External Header. sub: ", sub)
	user, err := models.GetUserViaExternalID(sub)
	switch err {
	case nil:
		err = internal.ToJSON(user, rw)
		if err != nil {
			uh.StandardHandler.ServiceLogger.Printf("[ERROR] serializing servers", err)
		}
	default:
		uh.StandardHandler.ServiceLogger.Println(ErrUserServicePostBase, err)
		rw.WriteHeader(http.StatusInternalServerError)
		err = internal.ToJSON(&GenericError{Message: err.Error()}, rw)

	}

}
