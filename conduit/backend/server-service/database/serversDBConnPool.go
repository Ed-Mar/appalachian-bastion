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

const serverDBConfigPath = "server-service/config/"
const serverDBConfigFileName = "localServerDB"
const serverDBConfigFileType = "env"

// func GetServersDBConnPool get the connection pool for the Server Service
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
			// Break out into a switch statment
			switch pgErr.Code {
			case pgerrcode.InvalidPassword:
				log.Println("[SQL ERROR] | Incorrect password: ", err)
				return nil, err
			default:
				log.Println("[SQL ERROR]: ", err)
				return nil, err
			}
		}
	}
	return db, nil
}
