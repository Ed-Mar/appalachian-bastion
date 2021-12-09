package handler

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	"server-api/data"
)

// Servers is a http.Handler
type Servers struct {
	logger *log.Logger
}

// NewServers creates a server handler with the given logger
func NewServers(logger *log.Logger) *Servers {
	return &Servers{logger}
}

// GetServers returns the Servers from the data store
func (server *Servers) GetServers(rw http.ResponseWriter, r *http.Request) {
	server.logger.Println("Handle GET Servers")

	// fetch the Servers from the datastore
	listServers := data.GetServers()

	// serialize the list to JSON
	err := listServers.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
func (server Servers) UpdateServers(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	server.logger.Println("Handle PUT Server", id)
	serv := r.Context().Value(KeyServer{}).(data.Server)

	err = data.UpdateServer(id, &serv)
	if err == data.ErrServerNotFound {
		http.Error(rw, "Server not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Server not found", http.StatusInternalServerError)
		return
	}
}

type KeyServer struct{}

func (server Servers) MiddlewareValidateServer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		serv := data.Server{}

		err := serv.FromJSON(r.Body)
		if err != nil {
			server.logger.Println("[ERROR] deserializing server", err)
			http.Error(rw, "Error reading server", http.StatusBadRequest)
			return
		}

		// validate the server
		err = serv.Validate()
		if err != nil {
			server.logger.Println("[ERROR] validating server", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating server: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		// add the Server to the context
		ctx := context.WithValue(r.Context(), KeyServer{}, serv)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
