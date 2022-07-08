package handlers

import (
	"backend/internal"
	"backend/server-service/models"
	"net/http"
)

// swagger:route PUT /servers servers updateServer
// UpdateSingleton a servers details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// UpdateSingleton handles PUT requests to update servers
func (server *Servers) UpdateSingleton(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	// fetch the servers from the context
	serv := r.Context().Value(KeyServer{}).(*models.Server)
	server.APILogger.Println("[DEBUG] updating record id", serv.ID)

	err := models.UpdateServer(*serv)
	if err == models.ErrServerNotFound {
		server.APILogger.Println("[ERROR] servers not found", err)

		rw.WriteHeader(http.StatusNotFound)
		err := internal.ToJSON(&GenericError{Message: "Server not found in database"}, rw)
		if err != nil {
			server.APILogger.Println(err)
		}
		return
	}

	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)
}
