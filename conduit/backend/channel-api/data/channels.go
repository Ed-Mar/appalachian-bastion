package data

import (
	"backend/database/postgres"
	"backend/internal"
	"fmt"
)

// ErrChannelNotFound is an error raised when a channel can not be found in the database
var ErrChannelNotFound = fmt.Errorf("channel not found")

// Channel defines the structure for an API channel
// swagger:model
type Channel struct {
	// the id for the channel in relation to  servers
	//
	// required: true
	// min: 1
	ID uint `json:"id" gorm:"primaryKey,autoIncrement,unique,not null"`

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
	Type   string `json:"type" gorm:"default:default_channel_type"`
	Server uint   `gorm:""`

	//Messages []uint64 `json:"messages,omitempty" gorm:"ForeignKey:ChannelID"`

	// This for database use not to be returned
	internal.CustomGromModel `json:"-"`
}
type Channels []*Channel

// GetChannels  returns all channels from the database
func GetChannels() (Channels, error) {
	db := postgres.GetPostgresDB()
	var channels []*Channel

	if err := db.Find(&channels).Error; err != nil {
		return nil, err
	}

	return channels, nil
}

// GetChannelByID returns a single channel which matches the id from the
// database.
func GetChannelByID(id uint) (*Channel, error) {

	db := postgres.GetPostgresDB()
	var channel *Channel
	if err := db.First(&channel, id).Error; err != nil {
		return nil, err
	}

	return channel, nil
}

// AddChannel to add a given channel to the database. Takes in a Channel and the Server it's attached as params.
func AddChannel(channel Channel) error {
	db := postgres.GetPostgresDB()
	if err := db.Create(&Channel{
		ID:          channel.ID,
		Name:        channel.Name,
		Description: channel.Description,
		Type:        channel.Type,
		//	Server: channel.Server
	}).Error; err != nil {
		return err
	}
	return nil
}

// DeleteChannel delete a channel with the given matching uint
func DeleteChannel(id uint) error {
	db := postgres.GetPostgresDB()
	var channel *Channel
	if err := db.Delete(&channel, id).Error; err != nil {
		return err
	}
	return nil
}

// UpdateChannel replaces with a given channel with the provided channel in the database.
// Note: that a Channel.ServerID cannot updated and the any provided ServerID will be ignored
func UpdateChannel(channel Channel) error {

	db := postgres.GetPostgresDB()
	if err := db.Model(&channel).Updates(&Channel{
		ID:          channel.ID,
		Name:        channel.Name,
		Description: channel.Description,
		Type:        channel.Type,
	}).Error; err != nil {
		return err
	}
	return nil
}

func findIndexByChannelID(id uint) int {
	// Copy pasta from the List all ID, will leave this for now
	// This is now dead code maybe will be used for future use if needed
	db := postgres.GetPostgresDB()
	var channels []*Channel
	db.Find(&channels)

	for i, channel := range channels {
		if channel.ID == id {
			return i
		}
	}

	return -1
}
