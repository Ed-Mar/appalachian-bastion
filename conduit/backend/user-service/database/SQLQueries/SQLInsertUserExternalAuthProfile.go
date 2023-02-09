package SQLQueries

func SQLInsertUserExternalAuthProfile() string {
	return sqlUserExternalAuthProfile
}

// sqlinsertUserProfile
const sqlUserExternalAuthProfile = `
INSERT INTO user_external_auth_profiles(
	external_auth_id,
	external_auth_issuer,
    external_auth_party,
	external_user_name,
	db_status,
	created_at
	)values ($1, $2, $3, $4, $5, $6);
`
