package SQLQueries

func SQLGetChannelViaID() string {
	// not positive if I need to send the error, but hey why not
	return sqlGetChannelViaID
}

const sqlGetChannelViaID = `
		SELECT *
		FROM channels
		WHERE
		deleted_at IS NULL AND channel_id = $1
		LIMIT 1;
`
