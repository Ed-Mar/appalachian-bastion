package database

//sqlCreateUUIDExtension
//Does not seem like '' quotes and needs the double "" quotes to be accecpted by the syntax
//Needs to be done in the desired database ( so check the connection to the correct DB to make sure it needed)
//CREATE EXTENSION loads a new extension into the current database. There must not be an extension of the same name already loaded.
const sqlCreateUUIDExtension = `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp"
`

// sqlCheckIfUUIDExtensionsExist checks if the uuid extension is installed on current connected database
// Needs to be run on database that the user is connected on to be work
// SELECT EXISTS (SELECT FROM pg_extension WHERE extname = ('uuid-ossp'))
const sqlCheckIfUUIDExtensionsExist = `
	SELECT EXISTS
	(SELECT FROM pg_extension
	WHERE extname = ('uuid-ossp'))
`

//sqlCreateChannelsTable sql used to create the table for servers
// I think it the db user needs to be the owner of the uuid extension to use of it
const sqlCreateChannelsTable = `
CREATE TABLE IF NOT EXISTS channels(
    channel_id uuid DEFAULT uuid_generate_v4() NOT NULL UNIQUE,
    PRIMARY KEY (channel_id),
    channel_name text NOT NULL,
    channel_description text,
    channel_type text DEFAULT 'TEXT',
    server_id uuid NOT NULL,
    status text NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp,
    deleted_at timestamp
)
`
