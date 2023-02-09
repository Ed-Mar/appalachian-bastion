package SQLQueries

func SQLInsertUser() string {
	return sqlInsertUserAuthenticationProfile
}

// sqlinsertUserProfile
const sqlInsertUserAuthenticationProfile = `
INSERT INTO userAuthenticationProfiles(
	external_id,
	external_auth_provider,
	external_auth_client_id,
	external_user_name,    
	user_type,                  
	status,
	created_at
	)values ($1, $2, $3, $4, $5, $6, $7);
`
