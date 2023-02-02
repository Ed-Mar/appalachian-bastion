package models

import (
	"backend/user-service/database"
	"backend/user-service/database/SQLQueries"

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

var ErrUserNotFound = fmt.Errorf("user not found")
var errDatabaseConnectionError = fmt.Errorf("[ERROR] [DATABASE]: Database Connection Error: ")
var errGenericSQLERROR = fmt.Errorf("[ERROR] [SQL]: ")

// User defines the structure for the User API
type User struct {
	///  "iss": "http://keycloak.test/realms/gatehouse",
	//  "aud": [
	//    "dev-conduit-rust"
	//  ],
	//  "exp": 1668503640,
	//  "iat": 1668503340,
	//  "auth_time": 1668503340,
	//  "nonce": "0impYhgbROniQbqeZblMgw",
	//  "acr": "1",
	//  "azp": "dev-conduit-rust",
	//  "at_hash": "kLw-rQym60VHKV6HSTSJqQ",
	//  "sub": "8b603c06-ee77-4f60-a0bd-72652808b861",
	//  "name": "Test User 0000",
	//  "given_name": "Test User",
	//  "family_name": "0000",
	//  "preferred_username": "test-user-0000"
	//}

	//User.ID
	//id used to id in the conduit environment
	//required:true
	UserID uuid.UUID `json:"conduit-id,UUID" db:"user_id"`

	//User.ExternalID
	//The UUID given from the external authentication. (sud)
	// required:true
	ExternalID uuid.UUID `json:"external-auth-id,UUID" db:"external_id"`

	//User.ExternalAuthProvider
	// Is the URL of the SSO provider used for authentication. (iss)
	// required:true
	ExternalAuthProvider string `json:"external-auth-provider" db:"external_auth_provider"`

	//User.ExternalAuthClientID
	// is the client name used to interface with the external SSO. (auth_party)
	// required:true
	ExternalAuthClientID string `json:"external-auth-client-id" db:"external_auth_client_id"`

	// User.ExternalUserName
	// Application Bastion | Gatehouse UserName
	// this will be in the(JWT AccessToken) or/and preferred_username(OAuth Identity). (upn)
	// required:true
	// max length: 128
	ExternalUserName string `json:"external-user-name" db:"external_user_name"`

	// User.DisplayUserName
	// This a User settable field that has a default of the GateHouse User.GateHouseUserName
	// required: true
	// max length: 128
	DisplayUserName string `json:"display-user-name" db:"default_username"`

	// User.UserType
	// Defines the type of User this in the scope of the Conduit Application as a whole
	// required:true
	UserType string `json:"user-type" db:"user_type"`

	// User.Servers
	// Define the servers that the user is in
	Servers []uuid.UUID `json:"servers" db:"servers"`

	Status    string     `json:"status" db:"status"`
	CreatedAt *time.Time `json:"-" db:"created_at"`
	UpdatedAt *time.Time `json:"-" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

type Users []*User

// GetAllUsers returns all users from the database
func GetAllUsers() (Users, error) {
	pool, err := database.GetUsersDBConnPool()
	defer pool.Close()
	if err != nil {
		//If error occurs just send the simple error to the user, and pop it in the logs for us
		log.Panic(err)
		return nil, errDatabaseConnectionError
	}
	var users []*User
	err = pgxscan.Select(context.Background(), pool, &users, SQLQueries.SQLGetEveryUser())
	var pgErr *pgconn.PgError
	if err != nil {
		// Checks if the error is PG Error
		if errors.As(err, &pgErr) {
			// Break out into a switch statement
			switch pgErr.Code {
			case pgerrcode.CaseNotFound:
				return nil, ErrUserNotFound
			default:
				log.Println(errGenericSQLERROR, pgErr)
				return nil, pgErr
			}
		} else {
			log.Panic("[ERROR]: Expected SQL Error got something else:  ", err)
			return nil, err
		}
	}
	return users, nil
}

// GetUserViaUserID TODO I am not positive I am going to use UserId over the incoming externalid(sid). The only reason I can think not is if I open up to other log authentication parties
// GetUserViaUserID returns user with matching user id
func GetUserViaUserID(userid uuid.UUID) (Users, error) {
	pool, err := database.GetUsersDBConnPool()
	defer pool.Close()
	if err != nil {
		//If error occurs just send the simple error to the user, and pop it in the logs for us
		log.Panic(err)
		return nil, errDatabaseConnectionError
	}
	var users []*User
	err = pgxscan.Select(context.Background(), pool, &users, SQLQueries.SQLGetUserViaMatchingUserID(), userid)
	var pgErr *pgconn.PgError
	if err != nil {
		// Checks if the error is PG Error
		if errors.As(err, &pgErr) {
			// Break out into a switch statement
			switch pgErr.Code {
			case pgerrcode.CaseNotFound:
				return nil, ErrUserNotFound
			default:
				log.Println(errGenericSQLERROR, pgErr)
				return nil, pgErr
			}
		} else {
			log.Panic("[ERROR]: Expected SQL Error got something else:  ", err)
			return nil, err
		}
	}
	return users, nil
}
func GetUserViaExternalID(externalUUID uuid.UUID) (Users, error) {
	pool, err := database.GetUsersDBConnPool()
	defer pool.Close()
	if err != nil {
		//If error occurs just send the simple error to the user, and pop it in the logs for us(who is us it just you dumb fuck)
		log.Panic(err)
		return nil, errDatabaseConnectionError
	}
	var users []*User
	err = pgxscan.Select(context.Background(), pool, &users, SQLQueries.SQLGetUserViaMatchingExternalID(), externalUUID)
	var pgErr *pgconn.PgError
	if err != nil {
		// Checks if the error is PG Error
		if errors.As(err, &pgErr) {
			// Break out into a switch statement
			switch pgErr.Code {
			case pgerrcode.CaseNotFound:
				return nil, ErrUserNotFound
			default:
				log.Println(errGenericSQLERROR, pgErr)
				return nil, pgErr
			}
		} else {
			log.Panic("[ERROR]: Expected SQL Error got something else:  ", err)
			return nil, err
		}
	}
	if len(users) < 1 {
		return nil, ErrUserNotFound
	}
	return users, nil
}

func AddUser(user User) error {
	pool, err := database.GetUsersDBConnPool()
	defer pool.Close()
	if err != nil {
		log.Println(errDatabaseConnectionError, err)
		return err
	}
	//Expected Order of insert
	//	1, external_id,
	//	2 external_auth_provider,
	//	3 external_auth_client_id,
	//	4 external_user_name,
	//  5 default_username,
	//	6 user_type,
	//	7 status,
	//	8 created_at
	_, err = pool.Exec(context.Background(), SQLQueries.SQLInsertUser(),
		user.ExternalID,
		user.ExternalAuthProvider,
		user.ExternalAuthClientID,
		user.ExternalUserName,
		user.ExternalUserName, // For User Initialization This will the UPN copied then can and should be updated by the user at a later date.
		user.UserType,
		user.Status,
		time.Now())
	var pgErr *pgconn.PgError
	if err != nil {
		// Checks if the error is Postgres Error
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			default:
				log.Println(errGenericSQLERROR, pgErr, pgErr.Code)
				return pgErr
			}
		} else {
			log.Panic("[ERROR]: Expected SQL Error got something else:  ", err)
			return err
		}
	}
	return nil

}
