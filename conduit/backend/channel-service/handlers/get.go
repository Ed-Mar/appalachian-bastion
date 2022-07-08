package handlers

import (
	"backend/channel-service/models"
	"backend/internal"
	"backend/internal/helper"
	"net/http"
)

//TODO Need to either Restrict this or just remove it due to No User ever be able to do this later

// ListEveryChannel Returns all the channel in the database
func (channel *Channels) ListEveryChannel(rw http.ResponseWriter, r *http.Request) {
	channel.APILogger.Println("[DEBUG] get every dam channel")

	leChannel, err := models.GetEveryChannel()
	if err != nil {
		channel.APILogger.Println("[ERROR]: ", err)
	}

	err = internal.ToJSON(leChannel, rw)
	if err != nil {
		// we should never be here but log the error just in-case
		channel.APILogger.Println("[ERROR] serializing servers", err)
	}

}

// swagger:route GET servers{id}/channels channels listchannels
// Return a list of channels from the database
// responses:
//	200: channelsResponse

// ListAllChannelsWithMatchingServerID handles GET requests and returns all current channels
func (channel *Channels) ListAllChannelsWithMatchingServerID(rw http.ResponseWriter, r *http.Request) {
	channel.APILogger.Println("[DEBUG] get all channels for given matching server id")
	rw.Header().Add("Content-Type", "application/json")

	// Grabbing Server UUID from the URI
	serverID, err := helper.GetUUIDFromReqParm(r, "serverID")
	if err != nil {
		// Bad Server UUID passed or can't convert it for some reason
		rw.WriteHeader(http.StatusBadRequest)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			// if encode to JSON fails just logged from the JSON side
			return
		}
		return
	}

	channel.APILogger.Println("[DEBUG] Getting All Channels for Server ID: ", serverID)
	channels, err := models.GetChannelsViaServerID(serverID)
	switch err {
	case nil:
	case models.ErrChannelNotFound:
		rw.WriteHeader(http.StatusNotFound)
		err := internal.ToJSON(GenericError{Message: err.Error()}, rw)
		if err != nil {
			channel.APILogger.Println("[ERROR] in JSON encoding: ", err)
		}
		return
	default:
		channel.APILogger.Println("[ERROR] fetching channels", err)

		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(GenericError{Message: err.Error()}, rw)
		if err != nil {
			channel.APILogger.Println("[ERROR] in JSON encoding: ", err)
		}
		return
	}

	err = internal.ToJSON(channels, rw)
	if err != nil {
		// we should never be here but log the error just encase
		channel.APILogger.Println("[ERROR] serializing channel(s)", err)
	}
}

// Return a list of channels from the database
// responses:
//	200: channelResponse
//	404: errorResponse

// ListSingle handles GET requests
func (channel *Channels) ListSingle(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	serverID, err := helper.GetUUIDFromReqParm(r, "channelID")
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

	channel.APILogger.Println("[DEBUG] get record id", serverID)

	serv, err := models.GetChannelViaChannelID(serverID)

	switch err {
	case nil:
	case models.ErrChannelNotFound:
		channel.APILogger.Println("[ERROR] fetching channel", err)

		rw.WriteHeader(http.StatusNotFound)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			channel.APILogger.Println("[ERROR] in JSON encoding: ", err)
		}
		return
	default:
		channel.APILogger.Println("[ERROR] fetching channel", err)

		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			channel.APILogger.Println("[ERROR] in JSON encoding: ", err)
		}
		return
	}

	err = internal.ToJSON(serv, rw)
	if err != nil {
		// we should never be here but log the error just incase
		channel.APILogger.Println("[ERROR] serializing channel", err)
	}
}
