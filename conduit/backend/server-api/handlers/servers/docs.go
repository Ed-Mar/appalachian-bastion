// Package handlers classification of Server API
//
// Documentation for Server API
//
//	Schemes: http
//	BasePath: /
//	Version: 0.0.1
//
//	Consumes:
//	- application/json
//
// swagger:meta
package handlers

import "backend/server-api/data"

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body GenericError
}

// Validation errors defined as an array of strings
// swagger:response errorValidation
type errorValidationWrapper struct {
	// Collection of the errors
	// in: body
	Body ValidationError
}

// A list of servers
// swagger:response ServersResponse
type serversResponseWrapper struct {
	// All current servers
	// in: body
	Body []data.Server
}

// Data structure representing a single Server
// swagger:response ServerResponse
type serverResponseWrapper struct {
	// Newly created server
	// in: body
	Body data.Server
}

// No content is returned by this API endpoint
// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

// swagger:parameters updateServer createServer
type serverParamsWrapper struct {
	// server data structure to Update or Create.
	// Note: the id field is ignored by update and create operations
	// in: body
	// required: true
	Body data.Server
}

// swagger:parameters listSingleServer deleteServer
type serverIDParamsWrapper struct {
	// The id of the server for which the operation relates
	// in: path
	// required: true
	ID int `json:"id"`
}
