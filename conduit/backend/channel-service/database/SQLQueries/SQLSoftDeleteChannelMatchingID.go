package SQLQueries

func SQLSoftDeleteChannelMatchingID() string {
	return sqlSoftDeleteServerWithMatchingId
}

const sqlSoftDeleteServerWithMatchingId = `
	UPDATE channels SET
	status = 'DELETED',
	deleted_at = now()
	WHERE channel_id::text = ($1)
	AND
    deleted_at IS NULL;
`
