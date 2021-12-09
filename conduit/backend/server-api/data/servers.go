package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-playground/validator"
)

// Server  defines structure for an API for Server
type Server struct {
	ID          int    `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	//TODO  implement 'Channel' Data structure for the severs struct
	//Channels 	[]
	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

// Servers a Collection of Server
type Servers []*Server

func (server *Server) Validate() error {
	validate := validator.New()
	return validate.Struct(server)
}

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
func UpdateServer(id int, server *Server) error {
	_, pos, err := findServer(id)
	if err != nil {
		return err
	}
	server.ID = id
	serverList[pos] = server

	return nil
}

var ErrServerNotFound = fmt.Errorf("server not found")

// findServer returns a server when given a matching server id
func findServer(id int) (*Server, int, error) {
	for i, server := range serverList {
		if server.ID == id {
			return server, i, nil
		}
	}
	return nil, -1, ErrServerNotFound
}

// serverList local data set for development
var serverList = []*Server{
	{
		ID:          1,
		Name:        "Server 1",
		Description: "Joking and Smoking",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Server 2",
		Description: "Smoking and Joking",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
