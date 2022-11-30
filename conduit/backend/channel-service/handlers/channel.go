package handlers

import (
	"backend/internal"
	"log"
)

type KeyChannel struct{}

// Channels gerneric-handlers for getting and updating Channels
type Channels struct {
	APILogger *log.Logger
	validator *internal.Validation
}

// NewChannels returns a new Channels gerneric-handlers with the given APILogger
func NewChannels(channelAPILogger *log.Logger, validator *internal.Validation) *Channels {
	return &Channels{channelAPILogger, validator}
}

// GenericError is a generic error message returned by a ID
type GenericError struct {
	Message string `json:"error-message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"error-messages"`
}
