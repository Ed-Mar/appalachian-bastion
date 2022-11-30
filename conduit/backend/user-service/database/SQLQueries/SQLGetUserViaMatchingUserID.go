package SQLQueries

func SQLGetUserViaMatchingUserID() string {
	return sqlGetUserViaMatchingUserID
}

const sqlGetUserViaMatchingUserID = `
	SELECT *
	FROM users
	WHERE
	deleted_at IS NULL AND user_id = $1; 
`
