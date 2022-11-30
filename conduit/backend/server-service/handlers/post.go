package handlers

import (
	"backend/server-service/models"
	"net/http"
)

// swagger:route POST /servers servers createServer
// CreateSingleton a new servers
//
// responses:
//	200: ServerResponse
//  422: errorValidation
//  501: errorResponse

// CreateSingleton handles POST requests to add new servers
func (server *ServerHandler) CreateSingleton(rw http.ResponseWriter, r *http.Request) {
	// fetch the servers from the context
	rw.Header().Add("Content-Type", "application/json")

	//Gets the passed JSON and coverts it's into a Server Struct
	leserver := r.Context().Value(KeyServer{}).(*models.Server)
	server.APILogger.Println("[DEBUG] Inserting Server")

	//Overwrites the income server.Status cause all item need to be put in pending until no issues can be confirmed then it will be updated to READY
	leserver.Status = "PENDING"

	err := models.AddServer(*leserver)
	if err != nil {
		return
	}
}
