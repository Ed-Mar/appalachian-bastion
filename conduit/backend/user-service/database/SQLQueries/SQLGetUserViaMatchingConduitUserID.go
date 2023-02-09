package SQLQueries

func SQLGetUserViaMatchingUserID() string {
	return sqlGetUserViaMatchingUserID
}

const sqlGetUserViaMatchingUserID = `
	SELECT *
	FROM userAuthenticationProfiles
	WHERE
	deleted_at IS NULL AND conduit_user_id = $1; 
`
