package handlers

import (
	"backend/channel-service/models"
	"backend/internal"
	"backend/internal/helper"
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
	id, _ := helper.GetURIParmWithMatchingName(r, "serverID")

	channel.APILogger.Println("[DEBUG] deleting record id", id)

	err := models.DeleteChannel(id)
	if err == models.ErrChannelNotFound {
		channel.APILogger.Println("[ERROR] deleting record id does not exist")

		rw.WriteHeader(http.StatusNotFound)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			channel.APILogger.Println("[ERROR] encoding to JSON: ", err)
		}
		return
	}

	if err != nil {
		channel.APILogger.Println("[ERROR] deleting record", err)

		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			channel.APILogger.Println("[ERROR] encoding to JSON: ", err)
		}
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
