package config

// sqlCreateChannelServiceDBUser
// create the channel service DB Users
const sqlCreateChannelServiceDBUser = `CREATE USER
    channel_service WITH
    NOCREATEDB
    NOCREATEROLE
    LOGIN ENCRYPTED PASSWORD 'channel_service_password'
`

// sqlCreateServerServiceDBUser
// create the server service DB Users
const sqlCreateServerServiceDBUser = `CREATE USER
    server_service WITH
    NOCREATEDB
    NOCREATEROLE
    LOGIN ENCRYPTED PASSWORD 'server_service_password'
`

// sqlCreatUserServiceDBUser
// create the user service DB Users
const sqlCreatUserServiceDBUser = `CREATE USER
    user_service WITH
    NOCREATEDB
    NOCREATEROLE
    LOGIN ENCRYPTED PASSWORD 'user_service_password'
`

// sqlCreateChannelServiceDB
// creates the actual database for the channel service
const sqlCreateChannelServiceDB = `
CREATE DATABASE channel_service
OWNER channel_service;`

// sqlCreateServerServiceDB
// creates the actual database for the server service
const sqlCreateServerServiceDB = `CREATE DATABASE server_service
OWNER server_service;`

// sqlCreateUserServiceDB
// creates the actual database for the user service
const sqlCreateUserServiceDB = `CREATE DATABASE user_service
OWNER user_service;`
