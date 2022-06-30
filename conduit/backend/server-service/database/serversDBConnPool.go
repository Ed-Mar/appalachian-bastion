package database

import (
	commondb "backend/internal/database"
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

//const serverDBConfigPath = "server-service/config/"
const serverDBConfigPath = "database/config/"
const serverDBConfigFileName = "localServerServiceDBConfig"
const serverDBConfigFileType = "env"

// GetServersDBConnPool  get the connection pool for the Server Service
func GetServersDBConnPool() (*pgxpool.Pool, error) {
	connString, err := commondb.GetDBPostgresDSN(serverDBConfigPath, serverDBConfigFileName, serverDBConfigFileType)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.Connect(context.Background(), connString)
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
			log.Panic("[ERROR]: Expected SQL Error got something else:  ", err)
			return nil, err
		}
	}
	return db, nil
}
