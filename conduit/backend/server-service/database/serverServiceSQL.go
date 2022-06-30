package database

// Note just putting these here will need them later, but not positive where some of these should go,

//sqlCreateServersTable sql used to create the table for servers
//CREATE TABLE IF NOT EXISTS servers(
//server_id uuid DEFAULT uuid_generate_v4 () NOT NULL UNIQUE,
//PRIMARY KEY (server_id),
//server_name text NOT NULL,
//server_description text,
//status text NOT NULL,
//created_at timestamp NOT NULL,
//updated_at timestamp ,
//deleted_at timestamp
//);
const sqlCreateServersTable = "" +
	"CREATE TABLE IF NOT EXISTS servers(" +
	" server_id uuid DEFAULT uuid_generate_v4 () NOT NULL UNIQUE," +
	" PRIMARY KEY (server_id)," +
	" server_name text NOT NULL," +
	" server_description text," +
	" status text NOT NULL," +
	" created_at timestamp NOT NULL," +
	" updated_at timestamp, " +
	" deleted_at timestamp" +
	");"

//sqlInsertServer used to insert a server to the servers table
//parms: serverName ServerDescription, Status, & Creation Timestamp
//INSERT INTO servers (server_name,server_description,status,created_at) VALUES('Pure SQL Insert','test','Fake_Status',2022-01-11 00:36:37.783025 )
const sqlInsertServer = "" +
	"INSERT INTO servers" +
	" (server_name," +
	" server_description," +
	" status," +
	" created_at" +
	" VALUES(($1)," +
	" ($2)," +
	" ($3)," +
	" ($4)"
