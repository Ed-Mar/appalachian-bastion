package SQLQueries

func SQLGetEveryChannel() string {
	// not positive if I need to send the error, but hey why not
	return sqlGetAllChannels
}

const sqlGetAllChannels = `
		SELECT *
		FROM channels
		WHERE deleted_at IS NULL; 
`
