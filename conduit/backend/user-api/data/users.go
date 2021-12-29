package data

import "backend/internal"

//User defines the structure for the User API
//swagger:model
type User struct {

	//User.ID
	//conduit id for the
	//required:true
	//min:1
	ID uint `json:"id" gorm:"primaryKey,unique"`

	// User.ABUserID
	// Application Bastion User ID
	// required:true
	// max length: 254
	ABUserID uint `json:"-"`

	// User.DefaultUserName
	// When a user joins a Server the default username for that server can be set here
	// required: false
	// max length: 64
	DefaultUserName string `json:"defaultUserName"`

	// User.Type
	// Defines the type of User this in the scope of the Conduit Application as a whole
	// required:true
	Type string `json:"type"`

	internal.CustomGromModel `json:"-"`
}
