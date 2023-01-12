package handlers

import (
	"backend/internal"
	"backend/internal/authentication/handlers"
	"fmt"
)

var ErrUserServiceBaseError = fmt.Errorf("[ERROR] [USER-SERVICE] ")
var WarnUserServiceBaseWarning = fmt.Errorf("[WARNING] [USER-SERVICE] ")

// UserHandlerKey is a key used for the  object in the context
type UserHandlerKey struct{}

type UserHandler struct {
	StandardHandler *handlers.ServiceHandler
	validator       *internal.Validation
}

func NewUserHandler(h *handlers.ServiceHandler, validator *internal.Validation) *UserHandler {
	return &UserHandler{h, validator}
}

// GenericError is a generic error message returned by a servers
type GenericError struct {
	Message string `json:"error-message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"error-messages"`
}
