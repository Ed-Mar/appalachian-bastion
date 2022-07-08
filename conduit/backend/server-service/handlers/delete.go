package handlers

import (
	"backend/internal"
	"backend/internal/helper"
	"backend/server-service/models"
	"net/http"
)

// swagger:route DELETE /servers/{id} servers deleteServer
// UpdateSingleton a servers details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  501: errorResponse

// DeleteSingleton handles DELETE requests and removes items from the database
func (server *Servers) DeleteSingleton(rw http.ResponseWriter, r *http.Request) {
	//Checks frist if that id is even the in database before getting to try and delete it
	leServerID, err := helper.GetUUIDFromReqParm(r, "serverID")
	switch err {
	case nil: //Not Error
	case helper.ErrIncorrectUUIDFormat: // Format Error with UUID passed
		rw.WriteHeader(http.StatusBadRequest)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			// if encode to JSON fails just logged from the JSON side
			return
		}
		return
	default: // Catch all error
		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			server.APILogger.Println("[ERROR] in JSON encoding: ", err)
		}
		return
	}

	server.APILogger.Println("[DEBUG] deleting record id", leServerID)
	//Sends the server Id to be deleted to the model function, should erturn
	err = models.DeleteServer(leServerID)
	switch err {
	case nil:
		rw.WriteHeader(http.StatusNoContent)
	case models.ErrServerNotFound:
		server.APILogger.Println("[ERROR] deleting record id does not exist")
		rw.WriteHeader(http.StatusNotFound)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			server.APILogger.Println("[ERROR] in JSON encoding: ", err)
			return
		}
		return
	default:
		server.APILogger.Println("[ERROR] deleting record", err)
		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			server.APILogger.Println("[ERROR] in JSON encoding: ", err)
			return
		}
		return
	}
}
