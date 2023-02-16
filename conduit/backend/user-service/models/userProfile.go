package models

import (
	"backend/user-service/database"
	"backend/user-service/database/SQLQueries"
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/pkg/errors"
	"log"
	"time"
)

type UserProfile struct {

	//UserExternalAuthProfile.ExternalAuthID
	//The UUID given from the external authentication. (sud)
	// The 'subject' that is returned from the External Auth System
	// required:true
	ExternalAuthID uuid.UUID `json:"external-auth-id,UUID" db:"external_auth_id"`

	//UserProfile.ConduitUserID
	//id used to id in the conduit environment
	//required:true
	UserID uuid.UUID `json:"conduit-id,UUID" db:"conduit_user_id"`

	DisplayName string `json:"conduit-display-name" db:"conduit_display_name"`

	UserType string `json:"user-type" db:"user_type"`

	//TODO This is a many-to-many database relationship. And will require its own saga on this sides and the server service.
	// So I am going to have comeback to this.
	//UserServers []uuid.UUID `json:"user-servers" db:"user-servers"`

	//// Database Control meta information
	// DBStatus to indicate locked or not
	DBStatus string `json:"-" db:"db_status"`

	CreatedAt *time.Time `json:"-" db:"created_at"`
	UpdatedAt *time.Time `json:"-" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

type UserProfiles []*UserProfile

func AddNewUserProfile(userprofile UserProfile) error {
	pool, err := database.GetUsersDBConnPool()
	defer pool.Close()
	if err != nil {
		err = errors.Wrap(errDatabaseConnectionError, " | ")
		return err
	}
	_, err = pool.Exec(context.Background(), SQLQueries.SQLInsetUserProfile(),
		userprofile.ExternalAuthID,
		userprofile.DisplayName,
		userprofile.UserType,
		userprofile.DBStatus,
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

func GetUserProfileViaExternalID(externalUUID uuid.UUID) (UserProfiles, error) {
	var errTooManyUsers = fmt.Errorf("[ERROR] [SQL]: Pulled to many Matching Users... Somehow")

	pool, err := database.GetUsersDBConnPool()
	defer pool.Close()
	if err != nil {
		//If error occurs just send the simple error to the user, and pop it in the logs for us(who is us it just you dumb fuck)
		log.Panic(err)
		return nil, errDatabaseConnectionError
	}
	log.Println("Pulling User.")
	var userProfiles []*UserProfile
	err = pgxscan.Select(context.Background(), pool, &userProfiles, SQLQueries.SQLGetUserProfileViaMatchingExternalID(), externalUUID)
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
	fmt.Printf("User Pulled %v\n", userProfiles)
	if len(userProfiles) < 1 {
		return nil, errTooManyUsers
	}
	return userProfiles, nil
}

// HardDelete
func HardDeleteUserProfile(externalUUID uuid.UUID) error {

	pool, err := database.GetUsersDBConnPool()
	defer pool.Close()
	if err != nil {
		//If error occurs just send the simple error to the user, and pop it in the logs for us(who is us it just you dumb fuck)
		log.Panic(err)
		return errDatabaseConnectionError
	}
	response, err := pool.Exec(context.Background(), SQLQueries.SQLHardDeleteUserProfile(), externalUUID)
	var pgErr *pgconn.PgError
	if err != nil {
		// Checks if the error is PG Error
		if errors.As(err, &pgErr) {
			// Break out into a switch statement
			switch pgErr.Code {
			case pgerrcode.CaseNotFound:
				return ErrUserNotFound
			default:
				log.Println(errGenericSQLERROR, pgErr)
				return pgErr
			}
		} else {
			log.Panic("[ERROR] [USER HARD DELETE]: Expected SQL Error got something else:  ", err)
			return err
		}
	}
	rowsDelete := response.RowsAffected()
	switch {
	case rowsDelete <= 0:
		{
			return ErrUserNotFound
		}
	case rowsDelete == 1:
		{
			return nil
		}
	case rowsDelete <= 1:
		{
			log.Default().Printf("[ERROR] [ERROR] [USERPROFILE] [HARD_DELETE_USER_PROFILE]: Deleted more than one user when trying to delete user: ", externalUUID.String())
			return ErrUserNotFound
		}
	}
	return nil
}
