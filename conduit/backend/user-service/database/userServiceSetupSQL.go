package database

// Note just putting these here will need them later, but not positive where some of these should go,

// sqlCreateChannelServiceDBUser Creates the Channel Service database user
// Needs to be done by another user
// also this needs to be hardcoded as trying to dynamic query with arguments cannot be done safely
const sqlCreateChannelServiceDBUser = `
CREATE USER channel_service
    ENCRYPTED PASSWORD 'channel_person_password'
    NOCREATEDB
    NOCREATEROLE
    NOSUPERUSER
    LOGIN;
`

// sqlCreateServiceDB Creates the Channel Database
// Needs to done by another user
const sqlCreateServerServiceDB = `
CREATE DATABASE user_service OWNER user_service;
`
