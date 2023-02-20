package SQLQueries

// SQLGetServerWithMatchingID returns a bool if matching server id is in servers table
// SELECT EXISTS(SELECT FROM servers WHERE server_id::text ='2b698f82-ffef-4faa-aa70-c9bb79073ce9' );
func SQLGetServerWithMatchingID() string {
	return sqlGetServerWithMatchingID
}

const sqlGetServerWithMatchingID = `
	SELECT EXISTS
		(SELECT FROM servers
		WHERE server_id::text = ($1));

`
