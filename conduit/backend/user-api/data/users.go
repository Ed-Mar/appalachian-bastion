package data

import (
	"backend/database/postgres"
	"backend/internal"
	server "backend/server-api/data"
)

//User defines the structure for the User API
//swagger:model
type User struct {

	//User.ID
	//conduit id for the
	//required:true
	//min:1
	ID uint `json:"id" gorm:"primaryKey,unique"`

	// User.ABUserID
	// Application Bastion UserID
	// required:true
	// max length: 254
	ABUserID uint `json:"-"`

	// User.DefaultUserName
	// When a user joins a Server the default username for that server can be set here
	// required: false
	// max length: 64
	DefaultUserName string `json:"defaultUserName" gorm:"default:default_username"`

	// User.Type
	// Defines the type of User this in the scope of the Conduit Application as a whole
	// required:true
	Type string `json:"type"`

	// User.Servers
	// Define the servers that the user is in
	Servers []*server.Server `json:"servers" gorm:"many2many:user_servers"`

	internal.CustomGromModel `json:"-"`
}

type Users []*User

// GetUsers returns all Users from the database
func GetUsers() (Users, error) {
	db := postgres.GetPostgresDB()
	var users []*User

	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByID returns a single user which matches the id from the
// database.
func GetUserByID(id uint) (*User, error) {

	db := postgres.GetPostgresDB()
	var user *User
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// AddUser to add a given user to the database. Takes in a User and the Server it's attached as params.
func AddUser(user User) error {
	db := postgres.GetPostgresDB()
	if err := db.Create(&User{
		ID:              user.ID,
		ABUserID:        user.ABUserID,
		DefaultUserName: user.DefaultUserName,
		Type:            user.Type,
		Servers:         user.Servers,
	}).Error; err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes a user from the database
func DeleteUser(id uint) error {
	db := postgres.GetPostgresDB()
	var user *User

	if err := db.Delete(&user, id).Error; err != nil {
		return err
	}
	return nil
}

// UpdateUser replaces a user in the database with the given item.
func UpdateUser(user User) error {
	db := postgres.GetPostgresDB()

	if err := db.Model(&user).Updates(&User{
		ID:              user.ID,
		ABUserID:        user.ABUserID,
		DefaultUserName: user.DefaultUserName,
		Type:            user.Type,
		Servers:         user.Servers,
	}).Error; err != nil {
		return err
	}
	return nil
}
