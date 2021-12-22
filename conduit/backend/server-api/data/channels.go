package data

import (
	"backend/internal"
	"fmt"
)

// ErrChannelNotFound is an error raised when a server can not be found in the database
var ErrChannelNotFound = fmt.Errorf("channel not found")

// Channel defines the structure for an API server
// swagger:model
type Channel struct {
	// the id for the channel in relation to  server
	//
	// required: false
	// min: 1
	ID int `json:"id" gorm:"primaryKey,autoIncrement"`

	// the name for this channel
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this channel
	//
	// required: false
	// max length: 512
	Description string `json:"description"`

	// the type of channel this is
	//
	// required: true
	// default: default_channel_type
	// max length: 64
	Type string `json:"type" gorm:"default:default_channel_type"`

	// This for database use not to be returned
	internal.CustomGromModel `json:"-"`
}
type Channels []*Channel
