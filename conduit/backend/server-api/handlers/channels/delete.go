package channels

import (
	"backend/internal"
	"backend/server-api/data"
	"net/http"
)

// swagger:route DELETE /channels/{id} channels deleteChannel
// Update a channels details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  501: errorResponse

// Delete handles DELETE requests and removes items from the database
func (channel *Channels) Delete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	id := getChannelID(r)

	channel.channelAPILogger.Println("[DEBUG] deleting record id", id)

	err := data.DeleteChannel(uint(id))
	if err == data.ErrChannelNotFound {
		channel.channelAPILogger.Println("[ERROR] deleting record id does not exist")

		rw.WriteHeader(http.StatusNotFound)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			channel.channelAPILogger.Println("[ERROR] encoding to JSON: ", err)
		}
		return
	}

	if err != nil {
		channel.channelAPILogger.Println("[ERROR] deleting record", err)

		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			channel.channelAPILogger.Println("[ERROR] encoding to JSON: ", err)
		}
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
