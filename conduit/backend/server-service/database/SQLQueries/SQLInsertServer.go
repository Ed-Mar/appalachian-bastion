package SQLQueries

// SQLInsertServer used to insert a server to the servers table
// parms: serverName ServerDescription, Status, & Creation Timestamp
// INSERT INTO servers (server_name,server_description,status,created_at) VALUES('Pure SQL Insert','test','Fake_Status',2022-01-11 00:36:37.783025 )
func SQLInsertServer() string {
	return sqlInsertServer
}

const sqlInsertServer = `
	INSERT INTO servers
	 (server_name,
	 server_description,
	 status,
	 created_at)
	 VALUES($1, $2, $3, $4);
`
