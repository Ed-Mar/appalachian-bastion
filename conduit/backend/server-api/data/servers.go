package data

import (
	"fmt"
)

// ErrServerNotFound is an error raised when a server can not be found in the database
var ErrServerNotFound = fmt.Errorf("server not found")

// Server defines the structure for an API server
// swagger:model
type Server struct {
	// the id for the server
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the server

	// the name for this server
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this server
	//
	// required: false
	// max length: 1000
	Description string `json:"description"`
}

// Servers defines a slice of Server
type Servers []*Server

// GetServers returns all servers from the database
func GetServers() Servers {
	return serverList
}

// GetServerByID returns a single server which matches the id from the
// database.
// If a server is not found this function returns a ServerNotFound error
func GetServerByID(id int) (*Server, error) {
	i := findIndexByServerID(id)
	if id == -1 {
		return nil, ErrServerNotFound
	}

	return serverList[i], nil
}

// UpdateServer replaces a server in the database with the given item.
// If a server with the given id does not exist in the database
// this function returns a ServerNotFound error
func UpdateServer(server Server) error {
	i := findIndexByServerID(server.ID)
	if i == -1 {
		return ErrServerNotFound
	}

	// update the server in the DB
	serverList[i] = &server

	return nil
}

// AddServer adds a new server to the database
//TODO interface with keycloak to remove that permissions in the role listing
func AddServer(server Server) {
	// get the next id in sequence
	//maxID := serverList[len(serverList)-1].ID
	//server.ID = maxID + 1

	serverList = append(serverList, &server)
}

// DeleteServer deletes a server from the database
//TODO interface with keycloak to remove that permissions in the role listing
func DeleteServer(id int) error {
	i := findIndexByServerID(id)
	if i == -1 {
		return ErrServerNotFound
	}

	serverList = append(serverList[:i], serverList[i+1])

	return nil
}

// findIndexByServerID finds the index of a server in the database
// internal
// returns -1 when no server can be found
func findIndexByServerID(id int) int {
	for i, server := range serverList {
		if server.ID == id {
			return i
		}
	}

	return -1
}

// serverList local data set for development
var serverList = []*Server{
	{
		ID:          1,
		Name:        "Server 1",
		Description: "Joking and Smoking",
	},
	{
		ID:          2,
		Name:        "Server 2",
		Description: "Smoking and Joking",
	},
}
