package data

import (
	"backend/database/postgres"
	"backend/internal"
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
	ID uint `json:"id" gorm:"primaryKey,autoIncrement,unique,not null"` // Unique identifier for the server

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

	// collection of any channels inside this server
	//
	// required: false
	Channels []*Channel `json:"channels,omitempty" gorm:"ForeignKey:ServerID"`
	//gorm:"foreignKey:ID;references:ID"`
	// This for database use not to be returned
	internal.CustomGromModel `json:"-"`
}

// Servers defines a slice of Server
type Servers []*Server

// GetServers returns all servers from the database
func GetServers() (Servers, error) {
	db := postgres.GetPostgresDB()
	var servers []*Server
	if err := db.Find(&servers).Error; err != nil {
		return nil, err
	}
	return servers, nil

}

// GetServerByID returns a single server which matches the id from the
// database.
// If a server is not found this function returns a ServerNotFound error
func GetServerByID(id uint) (*Server, error) {
	if id <= 0 {
		return nil, ErrServerNotFound
	}

	db := postgres.GetPostgresDB()
	var server *Server

	db.First(&server, id)
	if err := db.First(&server, id).Error; err != nil {
		return nil, err
	}
	return server, nil
}

// UpdateServer replaces a server in the database with the given item.
// If a server with the given id does not exist in the database
// this function returns a ServerNotFound error
func UpdateServer(server Server) error {
	i := findIndexByServerID(server.ID)
	if i == -1 {
		return ErrServerNotFound
	}

	db := postgres.GetPostgresDB()
	// update the server in the DB
	// Update attributes with `struct `, will only update non-zero fields
	db.Model(&server).Updates(Server{ID: server.ID, Name: server.Name, Description: server.Description, Channels: server.Channels})
	// UPDATE servers SET ID=server.id, name=18, ....

	return nil
}

// AddServer adds a new server to the database
//TODO interface with keycloak to remove that permissions in the role listing
func AddServer(server Server) error {
	db := postgres.GetPostgresDB()
	// If the server passed to be created has a channel(s) passed along with
	// then those channel(s) are pass along the ID (C)RUD operation for the creations of those channels

	//if server.Channels != nil {
	//	for _, channel := range server.Channels {
	//		err := AddChannel(channel, server.ID)
	//		if err != nil {
	//			return err
	//		}
	//	}

	if err := db.Create(&Server{ID: server.ID, Name: server.Name, Description: server.Description, Channels: server.Channels}).Error; err != nil {
		return err
	}

	//} else {
	//// If there is no channel passed with the server then it will just create a server
	//if err := db.Create(&Server{ID: server.ID, Name: server.Name, Description: server.Description}).Error; err != nil {
	//return err
	//}
	//}
	//
	return nil
}

// DeleteServer deletes a server from the database
//TODO interface with keycloak to remove that permissions in the role listing
func DeleteServer(id uint) error {
	i := findIndexByServerID(id)
	if i == -1 {
		return ErrServerNotFound
	}

	db := postgres.GetPostgresDB()
	var server *Server

	db.Delete(&server, id)
	// DELETE FROM servers WHERE id =i;

	return nil
}

// findIndexByServerID finds the index of a server in the database
// internal
// returns -1 when no server can be found
func findIndexByServerID(id uint) int {
	// Copy pasta from the List all Severs, will leave this for now
	// This is now dead code maybe will be used for future use if needed
	db := postgres.GetPostgresDB()
	var servers []*Server
	db.Find(&servers)

	for i, server := range servers {
		if server.ID == id {
			return i
		}
	}

	return -1
}
