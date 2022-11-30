package handlers

import (
	"backend/internal/authentication/handlers"
	"log"
	"net/http"
)

func (DevObj *DevTest) PostTest(rw http.ResponseWriter, r *http.Request) {
	// fetch the servers from the context
	rw.Header().Add("Content-Type", "application/json")

	log.Println(handlers.GetAuthFromContext(r.Context()))

}
