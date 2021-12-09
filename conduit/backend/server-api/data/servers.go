package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Server  defines structure for an API for Server
type Server struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	//TODO  implement 'Channel' Data structure for the severs struct
	//Channels 	[]
	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

// Servers a Collection of Server
type Servers []*Server

// FromJSON  serializes the JSON of the collection to data
func (server *Server) FromJSON(reader io.Reader) error {
	e := json.NewDecoder(reader)
	return e.Decode(server)
}

// ToJSON serializes the contents of the collection to JSON
func (server *Servers) ToJSON(writer io.Writer) error {
	e := json.NewEncoder(writer)
	return e.Encode(server)
}

// GetServers returns a list of servers
func GetServers() Servers {
	return serverList
}

// findServer returns a server when given a matching server id
func findServer(id int) (*Server, int, error) {
	for i, server := range serverList {
		if server.ID == id {
			return server, i, nil
		}
	}
	var ErrServerNotFound = fmt.Errorf("server not found")
	return nil, -1, ErrServerNotFound
}

// serverList local data set for development
var serverList = []*Server{
	&Server{
		ID:          1,
		Name:        "Server 1",
		Description: "Joking and Smoking",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Server{
		ID:          2,
		Name:        "Server 2",
		Description: "Smoking and Joking",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
