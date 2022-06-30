package database

import (
	commonDB "backend/internal/database"
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

//TODO Make this less complex to get to the file
//const channelDBConfigPath = "channel/config/"
const channelDBConfigPath = "database/config/"
const channelDBConfigFileName = "localChannelServiceDBConfig"
const DBConfigFileType = "env"

// GetChannelsDBConnPool  get the connection pool for the Channel Service
func GetChannelsDBConnPool() (*pgxpool.Pool, error) {
	connString, err := commonDB.GetDBPostgresDSN(channelDBConfigPath, channelDBConfigFileName, DBConfigFileType)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.Connect(context.Background(), connString)
	// TODO Figure out if I defer close here or do I have do where I use outside this
	if err != nil {

		var pgErr *pgconn.PgError
		// Checks if the error is PG Error
		if errors.As(err, &pgErr) {
			// Break out into a switch statement
			switch pgErr.Code {
			case pgerrcode.InvalidPassword:
				log.Println("[SQL ERROR] | Incorrect password: ", pgErr)
				// Need to check if the password is incorrect or if the user has not been created.
				return nil, pgErr
			default:
				log.Println("[SQL ERROR]: ", pgErr)
				return nil, pgErr
			}
		} else {
			log.Println("[ERROR]: UnExpected Non-SQL Error: ", err)
			return nil, err
		}
	}
	return db, nil
}
