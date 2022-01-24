package models

import (
	server_database "backend/server-service/database"
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
	pool, err := server_database.GetServersDBConnPool()
	if err != nil {
		return nil, err
	}
	var servers []*Server
	err = pgxscan.Select(context.Background(), pool, &servers, sqlGetAllServers)
	if err != nil {
		return nil, err
	}

	servers, err = removeSoftDeletedItems(servers)
	//Don't think this should happen, but encase it does
	if err == errNoServersInSlice {
		log.Println(errNoServersInSlice)
		return nil, ErrServerNotFound
	} else if err != nil {
		return nil, err
	}

	return servers, nil

}

// GetServerByID returns a single servers which matches the id from the
// database.
// If a servers is not found this function returns a ServerNotFound error
func GetServerByID(id uuid.UUID) (*Server, error) {

	pool, err := server_database.GetServersDBConnPool()
	if err != nil {
		log.Println(errDatabaseConnectionError, err)
		return nil, err
	}
	// Ok I had the Select working with ease while trying to use the pgxscan.Get
	//which is for only for one row was giving me lip,
	//so I just used the working select bit and grabbed the first item
	var servers []*Server
	err = pgxscan.Select(context.Background(), pool, &servers, sqlGetServerWithMatchingID, id)

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
	// I don't think this should ever happen, but a record of it would be nice it does happen
	if len(servers) > 1 {
		log.Printf("[WARRNING] GetServerByMatchingID is returning more than one item.")
		return servers[0], nil
	}
	// Checks if the server is deleted
	if servers[0].DeletedAt == nil {
		return servers[0], nil
	}
	return nil, ErrServerNotFound

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

	pool, err := server_database.GetServersDBConnPool()
	if err != nil {
		log.Println(errDatabaseConnectionError, err)
		return err
	}

	//should return "UPDATE 1"
	_, err = pool.Exec(context.Background(), sqlUpdateServerWithMatchingId, server.ID, server.Name, server.Description, server.Status, time.Now())
	if err != nil {
		log.Println(errGenericSQLERROR, err)
		return err
	}
	return nil
}

// AddServer adds a new servers to the database
//TODO interface with keycloak to remove that permissions in the role listing
func AddServer(server Server) error {
	pool, err := server_database.GetServersDBConnPool()
	if err != nil {
		log.Println(errDatabaseConnectionError, err)
		return err
	}
	_, err = pool.Exec(context.Background(), sqlInsertServer, server.Name, server.Description, server.Status, time.Now())
	if err != nil {
		log.Println(errGenericSQLERROR, err)
		return err
	}
	//SAMPLE output "INSERT 0 1"
	return nil
}

// DeleteServer deletes a servers from the database
//TODO interface with keycloak to remove that permissions in the role listing
func DeleteServer(id uuid.UUID) error {
	matchingID, err := doesServerExistWithMatchingID(id)
	if err != nil {
		return err
	}
	if matchingID == false {
		return ErrServerNotFound
	}
	pool, err := server_database.GetServersDBConnPool()
	if err != nil {
		log.Println(errDatabaseConnectionError, err)
		return err
	}
	_, err = pool.Exec(context.Background(), sqlSoftDeleteServerWithMatchingId, id, "DELETED", time.Now())
	if err != nil {
		log.Println(errGenericSQLERROR, err)
		return err
	}
	return nil
}

func doesServerExistWithMatchingID(id uuid.UUID) (bool, error) {
	pool, err := server_database.GetServersDBConnPool()
	if err != nil {
		log.Println(errDatabaseConnectionError, err)
		return false, err
	}
	var doesServerExist bool
	err = pool.QueryRow(context.Background(), sqlDoesServerExistWithMatchingID, id).Scan(&doesServerExist)

	if err != nil {
		log.Println(errGenericSQLERROR, err)
		return false, err
	}
	return doesServerExist, nil
}

func removeSoftDeletedItems(servers []*Server) (Servers, error) {
	if len(servers) > 0 {
		var temp = len(servers)
		for index := 0; index < temp; index++ {
			// encase you need this later
			//log.Printf("Index: %d | Length: %d", index, temp)
			//log.Printf("Status: %v | Sever Delete At: %v", servers[index].Status, servers[index].DeletedAt)
			if servers[index].DeletedAt != nil {
				//checks the if the last element needs it
				//bound check to make sure it's not the last element in the slice
				if index >= len(servers) {
					servers = servers[:len(servers)-1]
					index++
				} else {
					servers[index] = servers[len(servers)-1]
					servers[len(servers)-1] = nil
					servers = servers[:len(servers)-1]
					index++
					// So this weird thing I found while making this mess
					// if the second to last element is removed the last element is not checked
					// due to the new size being smaller and meeting the end condition of the loop
					// I did it this way over the append edit, due wanting to go fast
					///-------
					// so what this does if the index equals the new length it does the check again
					// and removes the last element if needed.
					var temp1 = len(servers)
					//log.Printf("Index: %d | Other-Length: %d", index, temp1)

					if index >= temp1-1 {
						if servers[index].DeletedAt != nil {
							servers = servers[:len(servers)-1]
							index++
							break // idk this wasn't working do I added the index++
						}
					}
				}
			}
		}
	} else {

		return servers, errNoServersInSlice
	}
	return servers, nil
}

var errNoServersInSlice = fmt.Errorf("no item in passed slice")

//sqlInsertServer used to insert a server to the servers table
//parms: serverName ServerDescription, Status, & Creation Timestamp
//INSERT INTO servers (server_name,server_description,status,created_at) VALUES('Pure SQL Insert','test','Fake_Status',2022-01-11 00:36:37.783025 )
const sqlInsertServer = "" +
	"INSERT INTO servers" +
	" (server_name," +
	" server_description," +
	" status," +
	" created_at)" +
	" VALUES($1, $2, $3, $4)"

//sqlGetAllServers get all servers
const sqlGetAllServers = `
		SELECT
		server_id,
		server_name,
		server_description,
		status,
		created_at,
		updated_at,
		deleted_at
		FROM servers`

//sqlGetServerWithMatchingID Get server with matching param UUID
const sqlGetServerWithMatchingID = "" +
	"SELECT" +
	" server_id," +
	" server_name," +
	" server_description," +
	" status," +
	" created_at," +
	" updated_at," +
	" deleted_at" +
	" FROM servers" +
	" WHERE server_id = $1"

//sqlDoesServerExistWithMatchingID returns a bool if matching server id is in servers table
//SELECT EXISTS(SELECT FROM servers WHERE server_id::text ='2b698f82-ffef-4faa-aa70-c9bb79073ce9' );
const sqlDoesServerExistWithMatchingID = `
	SELECT EXISTS
		(SELECT FROM servers
		WHERE server_id::text = ($1));
`

//sqlUpdateServerWithMatchingId updates server_name, server_description, status, and updated_at with given matching ID
//UPDATE servers
//SET server_name = 'Updated Server Name1', server_description= 'Updated Server Description1', status = 'UPDATED1', updated_at = now()
//WHERE server_id = '1aef01af-5f1d-4c4d-8747-ea957a3c8944';
const sqlUpdateServerWithMatchingId = `
	UPDATE servers
	SET server_name = ($2),
	server_description = ($3),
	status = ($4),
	updated_at = ($5)
	WHERE server_id::text = ($1);
`

//sqlSoftDeleteServerWithMatchingId Updates the server with given matching id
//UPDATE servers SET status = 'DELETED', deleted_at = now()
//WHERE server_id::text = 'dafbccb4-3e45-4dd4-ab36-9ff82c69cbc4';
const sqlSoftDeleteServerWithMatchingId = `
	UPDATE servers
	SET status = ($2),
	deleted_at = ($3)
	WHERE server_id::text = ($1);
`
