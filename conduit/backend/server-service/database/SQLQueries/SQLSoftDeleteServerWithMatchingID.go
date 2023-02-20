package SQLQueries

// SQLSoftDeleteServerWithMatchingID Updates the server with given matching id
// UPDATE servers SET status = 'DELETED', deleted_at = now()
// WHERE server_id::text = 'dafbccb4-3e45-4dd4-ab36-9ff82c69cbc4';
func SQLSoftDeleteServerWithMatchingID() string {
	return sqlSoftDeleteServerWithMatchingID
}

const sqlSoftDeleteServerWithMatchingID = `
	UPDATE servers
	SET status = ($2),
	deleted_at = ($3)
	WHERE server_id::text = ($1);
`
