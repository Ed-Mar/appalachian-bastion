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

// TODO Make this less complex to get to the file
// const userDBConfigPath = "channel/config/"
const userDBConfigPath = "channel-service/database/config/"
const userDBConfigFileName = "localChannelServiceDBConfig"
const DBConfigFileType = "env"

// GetUsersDBConnPool  get the connection pool for the User Service
func GetUsersDBConnPool() (*pgxpool.Pool, error) {
	connString, err := commonDB.GetDBPostgresDSN(userDBConfigPath, userDBConfigFileName, DBConfigFileType)
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
