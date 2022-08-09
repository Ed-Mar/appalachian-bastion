package handlers

import (
	"backend/channel-service/models"
	"backend/channel-service/saga"
	saga_pattern "backend/internal/saga-pattern"
	"net/http"
)

// swagger:route POST /servers/{id}/channels channels createChannel
// CreateSingleton a new servers
//
// responses:
//	200: ChannelResponse
//  422: errorValidation
//  501: errorResponse

// CreateChannelWithoutParms handles POST requests to add new channel
func (channel *Channels) CreateChannelWithoutParms(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	leChannel := r.Context().Value(KeyChannel{}).(*models.Channel)
	channel.APILogger.Println("[DEBUG] Inserting Channel")

	createChannelSaga, err := saga.CreateChannelSaga(*leChannel)
	if err != nil {
		channel.APILogger.Printf("SAGA was setup Incorrect", err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
	err = saga_pattern.ExecuteSaga(*createChannelSaga)
	if err != nil {
		channel.APILogger.Printf("[ERROR]: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
