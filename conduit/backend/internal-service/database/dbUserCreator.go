package database

import (
	common "backend/internal/database"
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

const serverDBConfigPath = "internal-service/config/"
const serverDBConfigFileName = "localDBSetupUserDBConfig"
const serverDBConfigFileType = "env"

// Yes I know this is stupid to send passwords to database user but I am planning on doing something with it later and I need the bones in now
func CreateServiceUser(serviceName string, servicePassWord string) (bool, error) {
	pool, err := getSetupUserDBConnPool()
	isConnected, err := checkDBConnection(pool)
	if err != nil {
		return false, err
	}
	if isConnected {
		_, err1 := pool.Exec(context.Background(), sqlCreateDBServiceUser, serviceName, servicePassWord)
		if err1 != nil {
			var pgErr *pgconn.PgError
			// Checks if the error is PG Error
			if errors.As(err1, &pgErr) {
				// Break out into a switch statement
				switch pgErr.Code {
				case pgerrcode.DuplicateObject:
					log.Println("[SQL ERROR] | User Already Exist: ", pgErr)
					return false, pgErr
				default:
					log.Println("[SQL ERROR]: ", pgErr)
					return false, pgErr
				}
			} else {
				log.Println("[ERROR]: Expected SQL Error got something else:  ", err)
				return false, err
			}
		} else {
			return true, nil
		}

	}
	return false, nil

}

func getSetupUserDBConnPool() (*pgxpool.Pool, error) {
	connString, err := common.GetDBPostgresDSN(serverDBConfigPath, serverDBConfigFileName, serverDBConfigFileType)
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
func checkDBConnection(pool *pgxpool.Pool) (bool, error) {
	err := pool.Ping(context.Background())
	if err != nil {
		log.Println("[ERROR]cannot connect to database: ", err)
		return false, err
	} else {
		return true, err
	}
}

// sqlCreateDBServiceUser is a sql query to create a service users needs to be provided the username and password
// "CREATE USER test_db_setup ENCRYPTED PASSWORD 'test-db-setup' CREATEDB  CREATEROLE LOGIN NOSUPERUSER"
const sqlCreateDBServiceUser = "" +
	"CREATE ROLE $1" +
	" ENCRYPTED PASSWORD '$2'" +
	" NOCREATEDB" +
	" NOCREATEUSER" +
	" NOSUPERUSER" +
	" LOGIN"
