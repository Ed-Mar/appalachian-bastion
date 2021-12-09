package handler

import (
	"context"
	"log"
	"net/http"

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

type KeyServer struct{}

func (server Servers) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Server{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			server.logger.Println("[ERROR] deserializing server", err)
			http.Error(rw, "Error reading server", http.StatusBadRequest)
			return
		}

		// add the Server to the context
		ctx := context.WithValue(r.Context(), KeyServer{}, prod)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
