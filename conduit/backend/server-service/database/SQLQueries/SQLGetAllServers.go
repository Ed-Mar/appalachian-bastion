package SQLQueries

// SQLGetAllServers get all servers
func SQLGetAllServers() string {
	return sqlGetAllServers
}

const sqlGetAllServers = `
		SELECT *
		FROM servers
		WHERE deleted_at IS NULL; 
`
