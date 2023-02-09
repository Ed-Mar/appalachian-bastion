package SQLQueries

func SQLGetUserProfileViaMatchingExternalID() string {
	return sqlGetUserProfileViaMatchingExternalID
}

const sqlGetUserProfileViaMatchingExternalID = `
	SELECT *
	FROM user_profiles
	WHERE
	deleted_at IS NULL AND external_auth_id = $1; 
`
