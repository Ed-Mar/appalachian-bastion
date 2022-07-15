package event_handler

import "log"

type KeySeverEvent struct{}

type ServerEvents struct {
	EventLogger *log.Logger
}

func NewServerEvent(serverEventLogger *log.Logger) *ServerEvents {
	return &ServerEvents{EventLogger: serverEventLogger}
}
