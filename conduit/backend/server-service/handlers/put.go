package handlers

import (
	"backend/internal"
	"backend/server-service/model"
	"net/http"
)

// swagger:route PUT /servers servers updateServer
// Update a servers details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// Update handles PUT requests to update servers
func (server *Servers) Update(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	// fetch the servers from the context
	serv := r.Context().Value(KeyServer{}).(*model.Server)
	server.severAPILogger.Println("[DEBUG] updating record id", serv.ID)

	err := model.UpdateServer(*serv)
	if err == model.ErrServerNotFound {
		server.severAPILogger.Println("[ERROR] servers not found", err)

		rw.WriteHeader(http.StatusNotFound)
		err := internal.ToJSON(&GenericError{Message: "Server not found in database"}, rw)
		if err != nil {
			server.severAPILogger.Println(err)
		}
		return
	}

	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)
}
