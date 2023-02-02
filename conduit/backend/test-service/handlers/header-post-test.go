package handlers

import (
	"backend/internal/authentication/handlers"
	"net/http"
)

func (DevObj *DevTest) PostTest(rw http.ResponseWriter, r *http.Request) {
	// fetch the servers from the uuid
	rw.Header().Add("Content-Type", "application/json")

	uuid, err := handlers.AuthGetExternalIDFromContext(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	DevObj.GenericHandler.ServiceLogger.Printf("sid from of current logged in User %v", uuid)

	upn, err := handlers.AuthGetUserPrincipalName(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	DevObj.GenericHandler.ServiceLogger.Printf("upn from of current logged in User %v", upn)
	iss, err := handlers.AuthGetJWTIssuer(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	DevObj.GenericHandler.ServiceLogger.Printf("iss from of current logged in User %v", iss)
	authP, err := handlers.AuthGetJWTAuthorizationParty(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	DevObj.GenericHandler.ServiceLogger.Printf("auth_party from of current logged in User %v", authP)
	groups, err := handlers.AuthGetJWTIssuerGroups(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	DevObj.GenericHandler.ServiceLogger.Printf("groups from of current logged in User %v", groups)

	auth_time, err := handlers.AuthGetJWTAuthorizationTime(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	DevObj.GenericHandler.ServiceLogger.Printf("auth_time from of current logged in User %v", auth_time)
	exp, err := handlers.AuthGetJWTExpiration(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	DevObj.GenericHandler.ServiceLogger.Printf("exp from of current logged in User %v", exp)
	lastV, err := handlers.AuthGetLastTokenIntrospective(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	DevObj.GenericHandler.ServiceLogger.Printf("last verification time from of current logged in User %v", lastV)

	authVtype, err := handlers.AuthGetLastTokenVerificationType(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	DevObj.GenericHandler.ServiceLogger.Printf("last verification type from of current logged in User %v", authVtype)

	return
}
