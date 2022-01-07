package handlers

import (
	"backend/internal"
	"backend/server-service/model"
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
	id := getServerID(r)

	server.severAPILogger.Println("[DEBUG] deleting record id", id)

	err := model.DeleteServer(uint(id))
	if err == model.ErrServerNotFound {
		server.severAPILogger.Println("[ERROR] deleting record id does not exist")

		rw.WriteHeader(http.StatusNotFound)
		internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		server.severAPILogger.Println("[ERROR] deleting record", err)

		rw.WriteHeader(http.StatusInternalServerError)
		internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
