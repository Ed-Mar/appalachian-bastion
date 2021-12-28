package channels

import (
	"backend/internal"
	"net/http"

	"backend/server-api/data"
)

// swagger:route PUT /servers/{id}/channels channels updateChannel
// Update a channels details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// Update handles PUT requests to update channels
func (channel *Channels) Update(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	// fetch the channel from the context
	serv := r.Context().Value(KeyChannel{}).(*data.Channel)
	channel.channelAPILogger.Println("[DEBUG] updating record id", serv.ID)

	err := data.UpdateChannel(*serv)
	if err == data.ErrChannelNotFound {
		channel.channelAPILogger.Println("[ERROR] channel not found", err)

		rw.WriteHeader(http.StatusNotFound)
		err := internal.ToJSON(&GenericError{Message: "channel not found in database"}, rw)
		if err != nil {
			channel.channelAPILogger.Println(err)
		}
		return
	}

	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)
}
