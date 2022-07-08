package handlers

import (
	"backend/channel-service/models"
	"backend/internal"
	"backend/internal/helper"
	"net/http"
)

// swagger:route DELETE /channels/{channelID} channels deleteChannel
// UpdateSingleton a channels details
//
// responses:
//	204: noContentResponse
//  400: BadRequestResponse
//  404: ResourceNotFoundResponse
//  500: ServerErrorResponse

// Delete handles DELETE requests and removes items from the database
func (channel *Channels) Delete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	channelID, err := helper.GetUUIDFromReqParm(r, "channelID")
	switch err {
	case nil: //Not Error
	case helper.ErrIncorrectUUIDFormat: // Format Error with UUID passed
		rw.WriteHeader(http.StatusBadRequest)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			// if encode to JSON fails just logged from the JSON side
			return
		}
		return
	default: // Catch all error
		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			channel.APILogger.Println("[ERROR] in JSON encoding: ", err)
		}
		return
	}
	channel.APILogger.Println("[DEBUG] deleting record id", channelID)
	err = models.DeleteChannel(channelID)
	switch err {
	case nil: // No Error
		rw.WriteHeader(http.StatusNoContent)
	case models.ErrChannelNotFound:
		channel.APILogger.Println("[ERROR] deleting record id does not exist")

		rw.WriteHeader(http.StatusNotFound)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			channel.APILogger.Println("[ERROR] encoding to JSON: ", err)
		}
		return
	default:
		channel.APILogger.Println("[ERROR] deleting record", err)

		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			channel.APILogger.Println("[ERROR] encoding to JSON: ", err)
		}
		return
	}
}
