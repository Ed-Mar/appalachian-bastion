package handlers

import (
	"backend/channel-service/models"
	"fmt"
	"net/http"
)

// swagger:route POST /servers/{id}/channels channels createChannel
// Create a new servers
//
// responses:
//	200: ChannelResponse
//  422: errorValidation
//  501: errorResponse

// CreateChannelWithoutParms handles POST requests to add new channel
func (channel *Channels) CreateChannelWithoutParms(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	leChannel := r.Context().Value(KeyChannel{}).(*models.Channel)
	channel.APILogger.Println("[DEBUG] Inserting Server")
	fmt.Println(leChannel.Status)
	if leChannel.Status == "" {
		leChannel.Status = "PENDING"
	}
	err := models.AddChannel(*leChannel)
	if err != nil {
		return
	}

}
