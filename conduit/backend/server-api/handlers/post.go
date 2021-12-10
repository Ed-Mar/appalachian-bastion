package handlers

import (
	"net/http"

	"server-api/data"
)

// swagger:route POST /servers servers createServer
// Create a new server
//
// responses:
//	200: ServerResponse
//  422: errorValidation
//  501: errorResponse

// Create handles POST requests to add new servers
func (server *Servers) Create(rw http.ResponseWriter, r *http.Request) {
	// fetch the server from the context
	srev := r.Context().Value(KeyServer{}).(data.Server)

	server.logger.Printf("[DEBUG] Inserting server: %#validator\n", srev)
	data.AddServer(srev)
}
