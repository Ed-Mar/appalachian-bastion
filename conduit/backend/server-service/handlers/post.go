package handlers

import (
	"backend/server-service/model"
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

	srev := r.Context().Value(KeyServer{}).(*model.Server)
	server.severAPILogger.Printf("[DEBUG] Inserting servers: %#validator\n", srev)

	err := model.AddServer(*srev)
	if err != nil {
		return
	}
}
