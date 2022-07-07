package SQLQueries

func SQLGetChannelsMatchingServerID() string {
	// not positive if I need to send the error, but hey why not
	return sqlGetChannelsMatchingServerID
}

const sqlGetChannelsMatchingServerID = `
		SELECT *
		FROM channels
		WHERE 
		deleted_at IS NULL AND server_id = $1; 
`
