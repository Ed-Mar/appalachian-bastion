package handlers

import (
	"backend/internal/authentication/handlers"
)

// KeyTestDevObj is a key used for the  object in the context
type KeyTestDevObj struct{}

type DevTest struct {
	GenericHandler *handlers.ServiceHandler
}

func NewDevTestHandler(h *handlers.ServiceHandler) *DevTest {
	return &DevTest{h}
}
