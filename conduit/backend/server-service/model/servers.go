package model

import (
	"backend/database/postgres"
	"backend/internal"
	"fmt"
	//"github.com/lib/pq"
)

// ErrServerNotFound is an error raised when a servers can not be found in the database
var ErrServerNotFound = fmt.Errorf("servers not found")

// Server defines the structure for an API servers
// swagger:model
type Server struct {
	// the id for the servers
	//
	// required: false
	// min: 1
	ID uint `json:"id" gorm:"primaryKey,autoIncrement,unique,not null"` // Unique identifier for the servers

	// the name for this servers
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this servers
	//
	// required: false
	// max length: 1000
	Description string `json:"description"`

	// collection of any channels inside this server
	//
	// required: false
	Channels []uint `json:"channels" gorm:"type:integer[],foreignKey:channel_id"`
	// collection of users in this server
	//
	// required: false
	Users []uint `json:"users" gorm:"type:integer[]"`

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

// GetServerByID returns a single servers which matches the id from the
// database.
// If a servers is not found this function returns a ServerNotFound error
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

// UpdateServer replaces a servers in the database with the given item.
// If a servers with the given id does not exist in the database
// this function returns a ServerNotFound error
func UpdateServer(server Server) error {
	i := findIndexByServerID(server.ID)
	if i == -1 {
		return ErrServerNotFound
	}

	db := postgres.GetPostgresDB()

	db.Model(&server).Updates(Server{
		ID:          server.ID,
		Name:        server.Name,
		Description: server.Description})
	// UPDATE servers SET ID=servers.id, name=18, ....

	return nil
}

// AddServer adds a new servers to the database
//TODO interface with keycloak to remove that permissions in the role listing
func AddServer(server Server) error {
	db := postgres.GetPostgresDB()

	if err := db.Create(&Server{
		ID:          server.ID,
		Name:        server.Name,
		Description: server.Description,
	}).Error; err != nil {
		return err
	}

	return nil
}

// DeleteServer deletes a servers from the database
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

// findIndexByServerID finds the index of a servers in the database
// internal
// returns -1 when no servers can be found
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
