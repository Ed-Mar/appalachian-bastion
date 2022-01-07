package channels

import (
	"backend/channel-api/data"
	"backend/internal"
	"net/http"
)

// swagger:route GET /channels channels listchannels
// Return a list of channels from the database
// responses:
//	200: channelsResponse

// ListAll handles GET requests and returns all current channels
func (channel *Channels) ListAll(rw http.ResponseWriter, r *http.Request) {
	channel.APILogger.Println("[DEBUG] get all records")
	rw.Header().Add("Content-Type", "application/json")

	servs, _ := data.GetChannels()

	err := internal.ToJSON(servs, rw)
	if err != nil {
		// we should never be here but log the error just incase
		channel.APILogger.Println("[ERROR] serializing channel", err)
	}
}

// swagger:route GET /channels/{id} channels listSinglechannel
// Return a list of channels from the database
// responses:
//	200: channelResponse
//	404: errorResponse

// ListSingle handles GET requests
func (channel *Channels) ListSingle(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	id := getChannelID(r)

	channel.APILogger.Println("[DEBUG] get record id", id)

	serv, err := data.GetChannelByID(uint(id))

	switch err {
	case nil:
	case data.ErrChannelNotFound:
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
