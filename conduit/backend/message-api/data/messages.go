package data

import (
	data2 "backend/channel-api/data"
	"backend/database/postgres"
	"backend/internal"
	"backend/user-api/data"
)

type Message struct {
	ID uint64 `json:"id" gorm:"primaryKey,autoIncrement"`

	UserID uint `json:"userID" gorm:"column:user_id"`
	User   data.User

	ChannelID uint `json:"channelID" gorm:"column:channel_id"`
	Channel   data2.Channel

	Type string `json:"type" gorm:"default:standard_text_only_message"`

	MessageReferenceID uint64 `json:"messageReferenceID,omitempty"`
	UserReferenceID    uint   `json:"userReferenceID,omitempty"`

	Status string `json:"status"`

	MessageContent string `json:"messageContent"`
	//Reactions			[]*Reaction

	internal.CustomGromModel `json:"-"`
}

type Messages []*Message

// GetMessages  returns all Messages from the database
func GetMessages() (Messages, error) {
	db := postgres.GetPostgresDB()
	var messages []*Message

	if err := db.Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

// GetMessageByID returns a single message which matches the id from the
// database.
func GetMessageByID(id uint) (*Message, error) {

	db := postgres.GetPostgresDB()
	var message *Message
	if err := db.First(&message, id).Error; err != nil {
		return nil, err
	}

	return message, nil
}

// AddMessage to add a given channel to the database. Takes in a Channel and the Server it's attached as params.
func AddMessage(message Message, userID uint, channelID uint) error {
	db := postgres.GetPostgresDB()
	if err := db.Create(&Message{
		ID:                 message.ID,
		UserID:             userID,
		ChannelID:          channelID,
		Type:               message.Type,
		MessageReferenceID: message.MessageReferenceID,
		UserReferenceID:    message.UserReferenceID,
		Status:             message.Status,
		MessageContent:     message.MessageContent,
	}).Error; err != nil {
		return err
	}
	return nil
}
