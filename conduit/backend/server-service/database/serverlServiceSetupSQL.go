package database

// sqlCreateServerServiceDBUser Creates the Server Service database user
// Needs to be done by another user
// also this needs to be hardcoded as trying to dynamic query with arguments cannot be done safely
const sqlCreateServerServiceDBUser = "" +
	"CREATE User server_service" +
	" ENCRYPTED PASSWORD 'server_service_password'" +
	" NOCREATEDB" +
	" NOCREATEROLE" +
	" NOSUPERUSER" +
	" LOGIN"

// sqlCreateServerServiceDB
const sqlCreateServerServiceDB = "" +
	"CREATE DATABASE server_service" +
	" OWNER server_service"

// sqlCheckIfUUIDExtensionsExist checks if the uuid extension is installed on current connected database
// Needs to be run on database that the user is connected on to be work
// SELECT EXISTS (SELECT FROM pg_extension WHERE extname = ('uuid-ossp'))
const sqlCheckIfUUIDExtensionsExist = "" +
	"SELECT EXISTS" +
	" (SELECT FROM pg_extension" +
	" WHERE extname = ('uuid-ossp'))"

//sqlDoesExtensionExistWithExtName checks aginist
const sqlDoesExtensionExistWithExtName = "" +
	"SELECT EXISTS" +
	" (SELECT FROM pg_extension" +
	" WHERE extname = ($1))"

//sqlCreateUUIDExtension
//Does not seem like '' quotes and needs the double "" quotes to be accecpted by the syntax
//Needs to be done in the desired
//CREATE EXTENSION loads a new extension into the current database. There must not be an extension of the same name already loaded.
const sqlCreateUUIDExtension = "" +
	"CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\""
