package handlers

import (
	"backend/internal"
	"backend/internal/helper"
	"backend/server-service/models"
	"net/http"
)

// swagger:route GET /servers servers listServers
// Return a list of servers from the database
// responses:
//	200: ServersResponse

// ListCollection handles GET requests and returns all current servers
func (server *ServerHandler) ListCollection(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Add("Content-Type", "application/json")

	server.APILogger.Println("[DEBUG] get all records")
	servers, err := models.GetServers()
	switch err {
	case nil:
		err = internal.ToJSON(servers, rw)
		if err != nil {
			server.APILogger.Println("[ERROR] serializing channel", err)
		}
	case models.ErrServerNotFound:
		rw.WriteHeader(http.StatusNotFound)
		err := internal.ToJSON(GenericError{Message: err.Error()}, rw)
		if err != nil {
			server.APILogger.Println("[ERROR] in JSON encoding: ", err)
		}
		return
	default:
		server.APILogger.Println("[ERROR] fetching channels", err)

		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(GenericError{Message: err.Error()}, rw)
		if err != nil {
			server.APILogger.Println("[ERROR] in JSON encoding: ", err)
		}
		return
	}
}

// swagger:route GET /servers/{id} servers listSingleServer
// Return a list of servers from the database
// responses:
//	200: ServerResponse
//	404: errorResponse

// ListSingleton handles GET requests
func (server *ServerHandler) ListSingleton(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	serverID, err := helper.GetUUIDFromReqParm(r, "serverID")
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			return
		}
		return
	}

	server.APILogger.Println("[DEBUG] get record id", serverID)
	serv, err := models.GetServerByID(serverID)
	switch err {
	case nil:
		err = internal.ToJSON(serv, rw)
		if err != nil {
			// we should never be here but log the error just incase
			server.APILogger.Println("[ERROR] serializing servers", err)
		}
	case models.ErrServerNotFound:
		server.APILogger.Println("[DEBUG] fetching servers", err)

		rw.WriteHeader(http.StatusNotFound)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			server.APILogger.Println("[ERROR] in JSON encoding: ", err)
		}
		return
	default:
		server.APILogger.Println("[ERROR] fetching servers", err)

		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			server.APILogger.Println("[ERROR] in JSON encoding: ", err)
		}
		return
	}

}
