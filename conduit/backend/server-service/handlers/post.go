package handlers

import (
	"backend/server-service/models"
	"net/http"
)

// swagger:route POST /servers servers createServer
// Create a new servers
//
// responses:
//	200: ServerResponse
//  422: errorValidation
//  501: errorResponse

// Create handles POST requests to add new servers
func (server *Servers) Create(rw http.ResponseWriter, r *http.Request) {
	// fetch the servers from the context
	rw.Header().Add("Content-Type", "application/json")

	srev := r.Context().Value(KeyServer{}).(*models.Server)
	server.severAPILogger.Println("[DEBUG] Inserting Server")
	//TODO remove this later
	if len(srev.Status) < 0 {
		srev.Status = "FAKE STATUS"
	}
	/////
	err := models.AddServer(*srev)
	if err != nil {
		return
	}
}
