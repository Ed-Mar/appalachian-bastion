package handlers

import (
	"backend/internal"
	"fmt"
	"log"
)

type KeyChannel struct{}

// Channels handlers for getting and updating Channels
type Channels struct {
	APILogger *log.Logger
	validator *internal.Validation
}

// NewChannels returns a new Channels handlers with the given APILogger
func NewChannels(channelAPILogger *log.Logger, validator *internal.Validation) *Channels {
	return &Channels{channelAPILogger, validator}
}

// ErrInvalidChannelPath is an error message when the ID path is not valid
var ErrInvalidChannelPath = fmt.Errorf("invalid Path, path should be /servers/[id]/channels/[id]")

// GenericError is a generic error message returned by a ID
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}
