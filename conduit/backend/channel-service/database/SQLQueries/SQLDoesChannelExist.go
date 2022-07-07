package SQLQueries

func SQLDoesChannelExistWithMatchingID() string {
	// not positive if I need to send the error, but hey why not
	return sqlDoesChannelExistWithMatchingID
}

const sqlDoesChannelExistWithMatchingID = `
	SELECT EXISTS
		(SELECT FROM channels
		WHERE channel_id::text = ($1)
		AND deleted_at IS NULL)LIMIT 1;
`
