package handler

import (
	"conduit/data"
	"log"
	"net/http"
)

// Severs is a http.Handler
type Sever struct {
	logger *log.Logger
}

// NewServers creates a server handler with the given logger
func NewServers(logger *log.Logger) *Servers {
	return &Servers{logger}
}

func (server *Servers) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle the request for a list of Servers
	if r.Method == http.MethodGet {
		server.getServers(rw, r)
		return
	}

	// if no method is satisfied return an error
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// getServers returns the Servers from the data store
func (p *Servers) getServers(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Servers")

	// fetch the Servers from the datastore
	lp := data.GetServers()

	// serialize the list to JSON
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
