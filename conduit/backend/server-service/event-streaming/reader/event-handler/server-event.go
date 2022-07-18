package event_handler

import "log"

type KeySeverEvent struct{}

type ServerEvents struct {
	Logger *log.Logger
}

func NewServerEvent(serverEventLogger *log.Logger) *ServerEvents {
	return &ServerEvents{Logger: serverEventLogger}
}
