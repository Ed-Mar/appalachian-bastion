package models

import (
	channelDB "backend/channel-service/database"
	channelSQL "backend/channel-service/database/SQLQueries"

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
	Name string `json:"name" validate:"required" db:"channel_name"`

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

func GetEveryChannel() (Channels, error) {
	pool, err := channelDB.GetChannelsDBConnPool()
	defer pool.Close()
	if err != nil {
		//If error occurs just send the simple error to the user, and pop it in the logs for us
		log.Panic(err)
		return nil, errDatabaseConnectionError
	}
	var channels []*Channel
	err = pgxscan.Select(context.Background(), pool, &channels, channelSQL.SQLGetEveryChannel())
	var pgErr *pgconn.PgError
	if err != nil {
		// Checks if the error is PG Error
		if errors.As(err, &pgErr) {
			// Break out into a switch statement
			switch pgErr.Code {
			case pgerrcode.CaseNotFound:
				//I think this will only happen if the database is empty, but mite as well leave it in
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
	//TODO noticed the error we use to get from database is not longer there mite have some to do with trying it with an array idk look into it later
	if len(channels) < 1 {
		return nil, ErrChannelNotFound
	}
	return channels, nil
}

// GetChannelsViaServerID  returns all channels from the database
func GetChannelsViaServerID(serverID uuid.UUID) (Channels, error) {
	pool, err := channelDB.GetChannelsDBConnPool()
	defer pool.Close()
	if err != nil {
		return nil, err
	}
	var channels []*Channel
	err = pgxscan.Select(context.Background(), pool, &channels, channelSQL.SQLGetChannelsMatchingServerID(), serverID)
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
	//TODO noticed the error we use to get from database is not longer there mite have some to do with trying it with an array idk look into it later
	if len(channels) < 1 {
		return nil, ErrChannelNotFound
	}
	return channels, nil
}

// GetChannelViaChannelID returns a single channel which matches the id from the
// database.
func GetChannelViaChannelID(channelID uuid.UUID) (Channels, error) {
	pool, err := channelDB.GetChannelsDBConnPool()
	defer pool.Close()
	if err != nil {
		log.Println(errDatabaseConnectionError, err)
		return nil, err
	}
	var channels []*Channel
	err = pgxscan.Select(context.Background(), pool, &channels, channelSQL.SQLGetChannelViaID(), channelID)
	var pgErr *pgconn.PgError
	if err != nil {
		log.Println(err)
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
	//TODO noticed the error we use to get from database is not longer there mite have some to do with trying it with an array idk look into it later
	if len(channels) < 1 {
		return nil, ErrChannelNotFound
	}

	return channels, nil

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
	_, err = pool.Exec(context.Background(), channelSQL.SQLInsertChannel(), channel.Name, channel.Description, channel.Type, channel.ServerID, channel.Status, time.Now())
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

// UpdateChannel replaces with a given channel with the provided channel in the database.
// Note: that a Channel.ServerID cannot update and any provided ServerID will be ignored
func UpdateChannel(channel Channel) error {
	queryingID, err := doesChannelExistWithMatchingID(channel.ID)
	// Need to send the Error on along
	if err != nil {
		return err
	}
	if !queryingID {
		return ErrChannelNotFound
	}
	pool, err := channelDB.GetChannelsDBConnPool()
	if err != nil {
		log.Println(errDatabaseConnectionError, err)
		return err
	}

	cmdTag, err1 := pool.Exec(context.Background(), channelSQL.SQLChannelUpdateMatchingID(), channel.ID, channel.Name, channel.Description, channel.Status)
	if err1 != nil {
		log.Println(errGenericSQLERROR, err1)
		return err1
	}

	const expectedUpdateReturnValue = "UPDATE 1"
	if cmdTag.String() != expectedUpdateReturnValue {
		return ErrChannelNotFound
	}
	return nil

}

// DeleteChannel delete a channel with the given matching uint
func DeleteChannel(channelID uuid.UUID) error {
	recordExist, err := doesChannelExistWithMatchingID(channelID)
	if err != nil {
		return err
	}
	if !recordExist {
		return ErrChannelNotFound
	}

	pool, err := channelDB.GetChannelsDBConnPool()
	defer pool.Close()
	if err != nil {
		log.Println(errDatabaseConnectionError, err)
		return err
	}
	cmdTag, err1 := pool.Exec(context.Background(), channelSQL.SQLSoftDeleteChannelMatchingID(), channelID)
	if err1 != nil {
		log.Println(errGenericSQLERROR, err1)
		return err1
	}
	const expectedUpdateReturnValue = "UPDATE 1"
	if cmdTag.String() != expectedUpdateReturnValue {
		log.Printf("[ERROR] [SERVER] [MODEL] [DELETE] Some werid happen and did not get UpdateSingleton 1 but instead got: %v\n", cmdTag.String())
		return errGenericSQLERROR
	}

	return nil

}
func doesChannelExistWithMatchingID(channelID uuid.UUID) (bool, error) {
	pool, err := channelDB.GetChannelsDBConnPool()
	defer pool.Close()
	if err != nil {
		log.Println(errDatabaseConnectionError, err)
		return false, err
	}
	var doesChannelExist bool
	err = pool.QueryRow(context.Background(), channelSQL.SQLDoesChannelExistWithMatchingID(), channelID).Scan(&doesChannelExist)

	if err != nil {
		log.Println(errGenericSQLERROR, err)
		return false, err
	}
	return doesChannelExist, nil

}
