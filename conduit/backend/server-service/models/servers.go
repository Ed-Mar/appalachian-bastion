package models

import (
	serverDatabase "backend/server-service/database"
	"backend/server-service/database/SQLQueries"
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

// ErrServerNotFound is an error raised when a servers can not be found in the database
var ErrServerNotFound = fmt.Errorf("server not found")
var errDatabaseConnectionError = fmt.Errorf("[ERROR] [DATABASE]: Database Connection Error: ")
var errGenericSQLERROR = fmt.Errorf("[ERROR] [SQL]: ")

// Server defines the structure for an API servers
// swagger:models
type Server struct {
	// the id for the servers
	//
	// required: false
	// min: 1public
	ID uuid.UUID `json:"id,UUID" db:"server_id"`
	// the name for this servers
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required" db:"server_name"`

	// the description for the server
	//
	// required: false
	Description string `json:"description" db:"server_description"`

	Status    string     `json:"status" db:"status"`
	CreatedAt *time.Time `json:"-" db:"created_at"`
	UpdatedAt *time.Time `json:"-" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

// Servers defines a slice of Server
type Servers []*Server

// GetServers returns all servers from the database
func GetServers() (Servers, error) {
	pool, err := serverDatabase.GetServersDBConnPool()
	if err != nil {
		return nil, err
	}
	var servers []*Server
	err = pgxscan.Select(context.Background(), pool, &servers, SQLQueries.SQLGetAllServers())
	var pgErr *pgconn.PgError
	if err != nil {
		// Checks if the error is PG Error
		if errors.As(err, &pgErr) {
			// Break out into a switch statement
			switch pgErr.Code {
			case pgerrcode.CaseNotFound:
				return nil, ErrServerNotFound
			default:
				log.Println(errGenericSQLERROR, pgErr)
				return nil, pgErr
			}
		} else {
			log.Panic("[ERROR]: Expected SQL Error got something else:  ", err)
			return nil, err
		}
	}
	return servers, nil
}

// GetServerByID returns a single servers which matches the id from the
// database.
// If a servers is not found this function returns a ServerNotFound error
func GetServerByID(id uuid.UUID) (*Server, error) {

	pool, err := serverDatabase.GetServersDBConnPool()
	if err != nil {
		log.Println(errDatabaseConnectionError, err)
		return nil, err
	}
	// Ok I had the Select working with ease while trying to use the pgxscan.Get
	//which is for only for one row was giving me lip,
	//so I just used the working select bit and grabbed the first item
	var servers []*Server
	err = pgxscan.Select(context.Background(), pool, &servers, SQLQueries.SQLGetServerWithMatchingID(), id)

	//TODO work out if this PgError Block is working, I do not think it is working. I knew it did one time
	var pgErr *pgconn.PgError
	if err != nil {
		// Checks if the error is PG Error
		if errors.As(err, &pgErr) {
			// Break out into a switch statement
			switch pgErr.Code {
			case pgerrcode.CaseNotFound:
				return nil, ErrServerNotFound
			default:
				log.Println(errGenericSQLERROR, pgErr)
				return nil, pgErr
			}
		} else {
			log.Panic("[ERROR]: Expected SQL Error got something else:  ", err)
			return nil, err
		}
	}
	if len(servers) < 1 {
		return nil, ErrServerNotFound
	}
	return servers[0], nil

}

// UpdateServer replaces a servers in the database with the given item.
// If a servers with the given id does not exist in the database
// this function returns a ServerNotFound error
func UpdateServer(server Server) error {
	matchingID, err := doesServerExistWithMatchingID(server.ID)
	if err != nil {
		return err
	}
	if matchingID == false {
		return ErrServerNotFound
	}

	pool, err := serverDatabase.GetServersDBConnPool()
	if err != nil {
		log.Println(errDatabaseConnectionError, err)
		return err
	}

	//should return "UPDATE 1"
	_, err = pool.Exec(context.Background(), SQLQueries.SQLUpdateServerWithMatchingID(), server.ID, server.Name, server.Description, server.Status, time.Now())
	if err != nil {
		log.Println(errGenericSQLERROR, err)
		return err
	}
	return nil
}

// AddServer adds a new servers to the database
// TODO interface with keycloak to remove that permissions in the role listing
func AddServer(server Server) error {
	pool, err := serverDatabase.GetServersDBConnPool()
	if err != nil {
		log.Println(errDatabaseConnectionError, err)
		return err
	}
	_, err = pool.Exec(context.Background(), SQLQueries.SQLInsertServer(), server.Name, server.Description, server.Status, time.Now())
	if err != nil {
		log.Println(errGenericSQLERROR, err)
		return err
	}
	//SAMPLE output "INSERT 0 1"
	return nil
}

// DeleteServer deletes a servers from the database
func DeleteServer(id uuid.UUID) error {
	matchingID, err := doesServerExistWithMatchingID(id)
	if err != nil {
		return err
	}
	if matchingID == false {
		return ErrServerNotFound
	}
	pool, err := serverDatabase.GetServersDBConnPool()
	if err != nil {
		log.Println(errDatabaseConnectionError, err)
		return err
	}
	_, err = pool.Exec(context.Background(), SQLQueries.SQLSoftDeleteServerWithMatchingID(), id, "DELETED", time.Now())
	if err != nil {
		log.Println(errGenericSQLERROR, err)
		return err
	}
	return nil
}

func doesServerExistWithMatchingID(id uuid.UUID) (bool, error) {
	pool, err := serverDatabase.GetServersDBConnPool()
	if err != nil {
		log.Println(errDatabaseConnectionError, err)
		return false, err
	}
	var doesServerExist bool
	err = pool.QueryRow(context.Background(), SQLQueries.SQLUpdateServerWithMatchingID(), id).Scan(&doesServerExist)

	if err != nil {
		log.Println(errGenericSQLERROR, err)
		return false, err
	}
	return doesServerExist, nil
}
