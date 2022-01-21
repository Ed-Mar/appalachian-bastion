package handlers

import (
	"backend/internal"
	"backend/server-service/models"
	"net/http"
)

// swagger:route DELETE /servers/{id} servers deleteServer
// Update a servers details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  501: errorResponse

// Delete handles DELETE requests and removes items from the database
func (server *Servers) Delete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	id, err := getServerID(r)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			server.severAPILogger.Println("[ERROR] JSON: ", err)
			return
		}
		return
	}

	server.severAPILogger.Println("[DEBUG] deleting record id", id)

	err = models.DeleteServer(id)
	switch err {
	case nil:
	case models.ErrServerNotFound:
		server.severAPILogger.Println("[ERROR] deleting record id does not exist")
		rw.WriteHeader(http.StatusNotFound)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			server.severAPILogger.Println("[ERROR] in JSON encoding: ", err)
			return
		}
		return
	default:
		server.severAPILogger.Println("[ERROR] deleting record", err)
		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			server.severAPILogger.Println("[ERROR] in JSON encoding: ", err)
			return
		}
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
