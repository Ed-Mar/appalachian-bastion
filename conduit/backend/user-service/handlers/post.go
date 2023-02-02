package handlers

import (
	"backend/internal"
	"backend/internal/authentication/handlers"
	"backend/user-service/models"
	"fmt"
	"net/http"
)

const ACCEPTING_NEW_USERS = true

var ErrUserServicePostBase = fmt.Errorf("%v [POST] ", ErrUserServiceBaseError)
var WarnUserServicePostBase = fmt.Errorf("%v [POST] ", WarnUserServiceBaseWarning)

var WarnUserServicePostNotAcceptingNewUser = fmt.Errorf("%v: Application is not accpecting new users at this time", WarnUserServicePostBase)

func (uh *UserHandler) CreateNewUserViaJSON(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	if !ACCEPTING_NEW_USERS {
		http.Error(rw, WarnUserServicePostNotAcceptingNewUser.Error(), http.StatusForbidden)
		uh.StandardHandler.ServiceLogger.Println(WarnUserServicePostNotAcceptingNewUser)
		return
	}
	//Grabs the Users form the incoming JSON
	passedUser := r.Context().Value(UserHandlerKey{}).(*models.User)
	uh.StandardHandler.ServiceLogger.Println("[DEBUG] Attempting to insert new User")
	err := models.AddUser(*passedUser)
	if err != nil {
		uh.StandardHandler.ServiceLogger.Println(ErrUserServicePostBase, err)
		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			// if encode to JSON fails just logged from the JSON side
			return
		}

	}

}
func (uh *UserHandler) CreateNewUserViaHeader(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	if !ACCEPTING_NEW_USERS {
		http.Error(rw, WarnUserServicePostNotAcceptingNewUser.Error(), http.StatusForbidden)
		uh.StandardHandler.ServiceLogger.Println(WarnUserServicePostNotAcceptingNewUser)
		return
	}
	sub, err := handlers.AuthGetExternalIDFromContext(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	uh.StandardHandler.ServiceLogger.Printf("sub from of current logged in User %v", sub)

	iss, err := handlers.AuthGetJWTIssuer(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	uh.StandardHandler.ServiceLogger.Printf("iss from of current logged in User %v", iss)

	authParty, err := handlers.AuthGetJWTAuthorizationParty(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	uh.StandardHandler.ServiceLogger.Printf("auth_party from of current logged in User %v", authParty)
	upn, err := handlers.AuthGetUserPrincipalName(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	uh.StandardHandler.ServiceLogger.Printf("upn from of current logged in User %v", upn)
	newUser := models.User{
		ExternalID:           sub,
		ExternalAuthProvider: iss,
		ExternalAuthClientID: authParty,
		ExternalUserName:     upn,
		DisplayUserName:      upn,
		UserType:             "default",
		Servers:              nil,
		Status:               "PENDING",
	}
	err = models.AddUser(newUser)
	if err != nil {
		uh.StandardHandler.ServiceLogger.Println(ErrUserServicePostBase, err)
		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			// if encode to JSON fails just logged from the JSON side
			return
		}

	}

}
