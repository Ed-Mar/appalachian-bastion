package handlers

import (
	"backend/channel-service/models"
	"backend/internal"
	"net/http"
)

// swagger:route PUT /servers/{id}/channels channels updateChannel
// UpdateWithoutIDParm a channels details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// UpdateWithoutIDParm handles PUT requests to update channels
func (channel *Channels) UpdateWithoutParms(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	// fetch the channel from the context
	leChannel := r.Context().Value(KeyChannel{}).(*models.Channel)
	channel.APILogger.Println("[DEBUG] updating record id", leChannel.ID)
	// Sending the Channel info to try and be updated
	err := models.UpdateChannel(*leChannel)
	switch err {
	// It is updated has happened and all gone as epected
	case nil:
		rw.WriteHeader(http.StatusOK)
		return
	// The updated was tried, but was not found either to never existing or has been soft-deleted
	case models.ErrChannelNotFound:

		rw.WriteHeader(http.StatusNotFound)
		err := internal.ToJSON(GenericError{Message: err.Error()}, rw)
		if err != nil {
			channel.APILogger.Println("[ERROR] in JSON encoding: ", err)
			return
		}
		return
	//some fucky-wucky has occurred, granted I did not leave a lot of options to be caught
	default:
		channel.APILogger.Println("[ERROR Unexpected Error when trying to update this channel: " + leChannel.ID.String())
		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(GenericError{Message: err.Error()}, rw)
		if err != nil {
			channel.APILogger.Println("[ERROR] in JSON encoding: ", err)
			return
		}
		return
	}
}
