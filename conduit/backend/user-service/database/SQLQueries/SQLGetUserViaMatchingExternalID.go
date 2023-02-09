package SQLQueries

func SQLGetUserViaMatchingExternalID() string {
	return sqlGetUserViaMatchingExternalID
}

const sqlGetUserViaMatchingExternalID = `
	SELECT *
	FROM userAuthenticationProfiles
	WHERE
	deleted_at IS NULL AND external_id = $1; 
`
