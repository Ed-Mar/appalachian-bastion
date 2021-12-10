package handlers

import (
	"net/http"

	"server-api/data"
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

	server.logger.Println("[DEBUG] deleting record id", id)

	err := data.DeleteServer(id)
	if err == data.ErrServerNotFound {
		server.logger.Println("[ERROR] deleting record id does not exist")

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		server.logger.Println("[ERROR] deleting record", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
