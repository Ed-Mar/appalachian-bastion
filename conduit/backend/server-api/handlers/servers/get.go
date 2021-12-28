package servers

import (
	"backend/internal"
	"net/http"

	"backend/server-api/data"
)

// swagger:route GET /servers servers listServers
// Return a list of servers from the database
// responses:
//	200: ServersResponse

// ListAll handles GET requests and returns all current servers
func (server *Servers) ListAll(rw http.ResponseWriter, r *http.Request) {
	server.severAPILogger.Println("[DEBUG] get all records")
	rw.Header().Add("Content-Type", "application/json")

	servs, err := data.GetServers()
	if err != nil {
		server.severAPILogger.Println("[ERROR]: ", err)
	}

	err = internal.ToJSON(servs, rw)
	if err != nil {
		// we should never be here but log the error just incase
		server.severAPILogger.Println("[ERROR] serializing server", err)
	}
}

// swagger:route GET /servers/{id} servers listSingleServer
// Return a list of servers from the database
// responses:
//	200: ServerResponse
//	404: errorResponse

// ListSingle handles GET requests
func (server *Servers) ListSingle(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	id := getServerID(r)

	server.severAPILogger.Println("[DEBUG] get record id", id)

	serv, err := data.GetServerByID(uint(id))

	switch err {
	case nil:
	case data.ErrServerNotFound:
		server.severAPILogger.Println("[ERROR] fetching server", err)

		rw.WriteHeader(http.StatusNotFound)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			server.severAPILogger.Println("[ERROR] in JSON encoding: ", err)
		}
		return
	default:
		server.severAPILogger.Println("[ERROR] fetching server", err)

		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			server.severAPILogger.Println("[ERROR] in JSON encoding: ", err)
		}
		return
	}

	err = internal.ToJSON(serv, rw)
	if err != nil {
		// we should never be here but log the error just incase
		server.severAPILogger.Println("[ERROR] serializing server", err)
	}
}
