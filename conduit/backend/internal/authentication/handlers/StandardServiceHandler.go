package handlers

import (
	"log"
)

//Note I am not positive this is the right idea to give all the services all the same id. Cause of the serives I have made
// don't differ and are microservices so I should not have to worry about using the same Key in the context

// KeyForHandler Used to Id it in the context
type KeyForHandler struct {
	//ServiceNameKey string
}

type ServiceHandler struct {
	ServiceName   string
	ServiceLogger *log.Logger
}

func NewHandler(serviceName string, serviceLogger *log.Logger) *ServiceHandler {
	return &ServiceHandler{serviceName, serviceLogger}
}

// GenericError is a generic error message returned by a ID
type GenericError struct {
	Message string `json:"error-message"`
}
