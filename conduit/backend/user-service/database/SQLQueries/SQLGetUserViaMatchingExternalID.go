package SQLQueries

func SQLGetUserViaMatchingExternalID() string {
	return sqlGetUserViaMatchingExternalID
}

const sqlGetUserViaMatchingExternalID = `
	SELECT *
	FROM users
	WHERE
	deleted_at IS NULL AND external_id = $1; 
`
