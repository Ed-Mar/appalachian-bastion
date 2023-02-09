package models

import (
	"backend/user-service/database"
	"backend/user-service/database/SQLQueries"
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"github.com/pkg/errors"
	"log"
	"time"
)

type UserExternalAuthProfile struct {

	//UserExternalAuthProfile.ExternalAuthID
	//The UUID given from the external authentication. (sud)
	// The 'subject' that is returned from the External Auth System
	// required:true
	ExternalAuthID uuid.UUID `json:"external-auth-id,UUID" db:"external_auth_id"`

	//UserExternalAuthProfile.ExternalAuthIssuer
	// Is the URL of the SSO provider used for authentication. (iss)
	// required:true
	ExternalAuthIssuer string `json:"external-auth-issuer" db:"external_auth_issuer"`

	//UserExternalAuthProfile.ExternalAuthParty
	// is the client name used to interface with the external SSO. (auth_party)
	// the id of the party or client ID used to authorize the user
	// required:true
	ExternalAuthParty string `json:"external-auth-party" db:"external_auth_party"`

	// UserExternalAuthProfile.ExternalUserName
	// Application Bastion | Gatehouse UserName
	// the username on the external site they used to log in with
	// this will be in the(JWT AccessToken)(OAuth Identity). (upn) or (preferred_username)
	// required:true
	// max length: 128
	ExternalUserName string `json:"external-user-name" db:"external_user_name"`

	//// Database Control meta information
	// DBStatus to indicate locked or not
	DBStatus string `json:"-" db:"db_status"`

	CreatedAt *time.Time `json:"-" db:"created_at"`
	UpdatedAt *time.Time `json:"-" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}
type UserExternalAuthProfiles []*UserExternalAuthProfile

func AddNewUserExternalAuthProfile(newUserExternalProfile UserExternalAuthProfile) error {
	pool, err := database.GetUsersDBConnPool()
	defer pool.Close()
	if err != nil {
		err = errors.Wrap(errDatabaseConnectionError, " | ")
		return err
	}
	_, err = pool.Exec(context.Background(), SQLQueries.SQLInsertUserExternalAuthProfile(),
		newUserExternalProfile.ExternalAuthID,
		newUserExternalProfile.ExternalAuthIssuer,
		newUserExternalProfile.ExternalAuthParty,
		newUserExternalProfile.ExternalUserName,
		newUserExternalProfile.DBStatus,
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
