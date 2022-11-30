package SQLQueries

func SQLGetChannelsMatchingServerID() string {
	return sqlGetChannelsMatchingServerID
}

const sqlGetChannelsMatchingServerID = `
		SELECT *
		FROM channels
		WHERE 
		deleted_at IS NULL AND server_id = $1; 
`
