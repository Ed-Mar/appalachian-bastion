package SQLQueries

// SQLUpdateServerWithMatchingID updates server_name, server_description, status, and updated_at with given matching ID
// UPDATE servers
// SET server_name = 'Updated Server Name1', server_description= 'Updated Server Description1', status = 'UPDATED1', updated_at = now()
// WHERE server_id = '1aef01af-5f1d-4c4d-8747-ea957a3c8944';
func SQLUpdateServerWithMatchingID() string {
	return sqlUpdateServerWithMatchingID
}

const sqlUpdateServerWithMatchingID = `
	UPDATE servers
	SET server_name = ($2),
	server_description = ($3),
	status = ($4),
	updated_at = ($5)
	WHERE server_id::text = ($1);
`
