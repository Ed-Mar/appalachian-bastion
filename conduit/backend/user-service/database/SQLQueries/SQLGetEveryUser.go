package SQLQueries

func SQLGetEveryUser() string {
	// not positive if I need to send the error, but hey why not
	return sqlGetAllUsers
}

const sqlGetAllUsers = `
		SELECT *
		FROM users
		WHERE deleted_at IS NULL; 
`
