package SQLQueries

func SQLInsetUserProfile() string {
	return sqlinsertUserProfile
}

// sqlinsertUserProfile
const sqlinsertUserProfile = `
INSERT INTO user_profiles(
	external_auth_id,
	conduit_display_name,
    user_type,                  
	db_status,
	created_at
	)values ($1, $2, $3, $4, $5);
`
