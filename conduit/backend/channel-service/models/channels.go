package models

import (
	channelDB "backend/channel-service/database"
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"log"
	"time"
)

// ErrChannelNotFound is an error raised when a channel can not be found in the database
var ErrChannelNotFound = fmt.Errorf("channel not found")
var errDatabaseConnectionError = fmt.Errorf("[ERROR] [DATABASE]: Database Connection Error: ")
var errGenericSQLERROR = fmt.Errorf("[ERROR] [SQL]: ")

// Channel defines the structure for an API channel
// swagger:models
type Channel struct {
	// the id for the channel in relation to  servers
	//
	// required: true
	ID uuid.UUID `json:"id,UUID" db:"channel_id"`

	// the name of the channel
	//
	// required: true
	Name string `json:"name" validate:"required" db:"channel_id"`

	// the description of the channel
	//
	// required: false
	Description string `json:"description" db:"channel_description"`

	// the type of channel
	//
	// required: true
	// default: default_channel_type
	// max length: 64
	Type     string    `json:"type" db:"channel_type"`
	ServerID uuid.UUID `json:"serverID" db:"server_id"`

	Status    string     `json:"status" db:"status"`
	CreatedAt *time.Time `json:"-" db:"created_at"`
	UpdatedAt *time.Time `json:"-" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}
type Channels []*Channel

// GetChannels  returns all channels from the database
func GetChannels(serverID uuid.UUID) (Channels, error) {
	pool, err := channelDB.GetChannelsDBConnPool()
	defer pool.Close()
	if err != nil {
		return nil, err
	}
	var channels []*Channel
	//TODO Finish this with the SQL
	err = pgxscan.Select(context.Background(), pool, &channels, sqlGetAllChannels)
	var pgErr *pgconn.PgError
	if err != nil {
		// Checks if the error is PG Error
		if errors.As(err, &pgErr) {
			// Break out into a switch statement
			switch pgErr.Code {
			case pgerrcode.CaseNotFound:
				return nil, ErrChannelNotFound
			default:
				log.Println(errGenericSQLERROR, pgErr)
				return nil, pgErr
			}
		} else {
			log.Panic("[ERROR]: Expected SQL Error got something else:  ", err)
			return nil, err
		}
	}
	return channels, nil
}

// GetChannelByID returns a single channel which matches the id from the
// database.
func GetChannelByID(serverID uuid.UUID) (*Channel, error) {

	return nil, nil
}

// AddChannel to add a given channel to the database. Takes in a Channel and the Server it's attached as params.
func AddChannel(channel Channel) error {
	pool, err := channelDB.GetChannelsDBConnPool()
	defer pool.Close()
	if err != nil {
		return err
	}

	//	channel_name = $1
	//	channel_description = $2
	//	channel_type = $3
	//	server_id = $4
	//	status = $5
	//	created_at = $6
	_, err = pool.Exec(context.Background(), sqlInsertChannel, channel.Name, channel.Description, channel.Type, channel.ServerID, channel.Status, time.Now())
	var pgErr *pgconn.PgError
	if err != nil {
		// Checks if the error is Postgres Error
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			default:
				log.Println(errGenericSQLERROR, pgErr)
				return pgErr
			}
		} else {
			log.Panic("[ERROR]: Expected SQL Error got something else:  ", err)
			return err
		}
	}
	return nil
}

// DeleteChannel delete a channel with the given matching uint
func DeleteChannel(id uuid.UUID) error {

	return nil
}

// UpdateChannel replaces with a given channel with the provided channel in the database.
// Note: that a Channel.ServerID cannot updated and the any provided ServerID will be ignored
func UpdateChannel(channel Channel) error {

	return nil
}

const sqlGetAllChannels = `
	

`

//sqlInsertChannel
//INSERT INTO channels(
//Example SQL:
//  channel_name, channel_description,channel_type ,server_id, status, created_at
//  )values ('PRUE SQL INSERT #2','INSERT VIA QUERY','TESTING_WITH_THA_BESTING','c7390d43-2cdd-42ba-a7c6-97a1aa847160','PRUE_SQL_INSERT',now());
const sqlInsertChannel = `
INSERT INTO channels(
	channel_name,
	channel_description,
	channel_type,
	server_id,
	status,
	created_at
	)values ($1, $2, $3, $4, $5, $6);
`
