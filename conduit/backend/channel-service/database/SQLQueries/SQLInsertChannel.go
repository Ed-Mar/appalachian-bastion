package SQLQueries

func SQLInsertChannel() string {
	return sqlInsertChannel
}

//sqlInsertChannel
//INSERT INTO channels(
//Example SQL:
//  channel_name, channel_description,channel_type ,server_id, status, created_at
//  )values ('PRUE SQL INSERT #2','INSERT VIA QUERY','TESTING_WITH_THA_BESTING','c7390d43-2cdd-42ba-a7c6-97a1aa847160','PRUE_SQL_INSERT',now());
const sqlInsertChannel = `
INSERT INTO channels(
	channel_name,
	channel_description,
	channel_type,
	server_id,
	status,
	created_at
	)values ($1, $2, $3, $4, $5, $6);
`
