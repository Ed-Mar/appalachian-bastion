package channels

import (
	"backend/channel-service/models"
	"backend/internal"
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

// Create handles POST requests to add new channels
func (channel *Channels) Create(rw http.ResponseWriter, r *http.Request) {
	// fetch the servers from the context
	rw.Header().Add("Content-Type", "application/json")

	leChannel := r.Context().Value(KeyChannel{}).(*models.Channel)
	switch getNumOfURIParms(r) {
	case 1:
		{
			serverID, err := getURIParmWithMatchingName(r, "serverID")
			channel.APILogger.Printf("[DEBUG] this is the server id from the parm: ", serverID)
			if serverID != leChannel.ServerID {
				rw.WriteHeader(http.StatusBadRequest)
				ErrMismatchServerIDParmBody := fmt.Errorf("URI ServerID parm doesn not equal ServerID of Channel in JSON Body ")
				err := internal.ToJSON(&GenericError{Message: ErrMismatchServerIDParmBody.Error()}, rw)
				if err != nil {
					channel.APILogger.Printf("[ERROR] [JSON] Mismatch in ServerID in pass URI parameters and JSON Body URI: %d | BODY: %v ", serverID, leChannel.ServerID)
				}
				return
			}
			channel.APILogger.Printf("[DEBUG] Inserting channel: %#validator\n", leChannel)
			err = models.AddChannel(*leChannel)
			if err != nil {
				channel.APILogger.Printf("[ERROR] adding Channel to db: ", err)

			}
		}

	}

}
func isServerIDPassViaURI(r *http.Request) {

}
