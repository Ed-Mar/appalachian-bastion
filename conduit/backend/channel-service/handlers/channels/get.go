package channels

import (
	"backend/channel-service/models"
	"backend/internal"
	"net/http"
)

// swagger:route GET servers{id}/channels channels listchannels
// Return a list of channels from the database
// responses:
//	200: channelsResponse

// ListAllWithMatchingID handles GET requests and returns all current channels
func (channel *Channels) ListAllWithMatchingID(rw http.ResponseWriter, r *http.Request) {
	channel.APILogger.Println("[DEBUG] get all records")
	rw.Header().Add("Content-Type", "application/json")

	serverID, err := getURIParmWithMatchingName(r, "serverID")
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
	channels, err := models.GetChannels(serverID)
	switch err {
	case nil:
	case models.ErrChannelNotFound:
		rw.WriteHeader(http.StatusNotFound)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			channel.APILogger.Println("[ERROR] in JSON encoding: ", err)
		}
		return
	default:
		channel.APILogger.Println("[ERROR] fetching channels", err)

		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
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

	serverID, err := getURIParmWithMatchingName(r, "serverName")
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

	channel.APILogger.Println("[DEBUG] get record id", serverID)

	serv, err := models.GetChannelByID(serverID)

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
