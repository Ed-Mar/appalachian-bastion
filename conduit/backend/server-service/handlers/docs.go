// Package servers Package gerneric-handlers' classification of Server API
//
// Documentation for Server API
//
//	Schemes: http
//	BasePath: /servers
//	Version: 0.0.1
//
//	Consumes:
//	- application/json
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import (
	"backend/server-service/models"
)

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
	Body []models.Server
}

// Data structure representing a single Server
// swagger:response ServerResponse
type serverResponseWrapper struct {
	// Newly created servers
	// in: body
	Body models.Server
}

// No content is returned by this API endpoint
// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

// swagger:parameters updateServer createServer
type serverParamsWrapper struct {
	// servers models structure to UpdateSingleton or CreateSingleton.
	// Note: the id field is ignored by update and create operations
	// in: body
	// required: true
	Body models.Server
}

// swagger:parameters listSingleServer deleteServer
type serverIDParamsWrapper struct {
	// The id of the servers for which the operation relates
	// in: path
	// required: true
	ID int `json:"id"`
}
