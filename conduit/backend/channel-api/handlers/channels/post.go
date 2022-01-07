package channels

import (
	"backend/channel-api/data"
	servers "backend/server-service/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

	chanl := r.Context().Value(KeyChannel{}).(*data.Channel)
	channel.APILogger.Printf("[DEBUG] Inserting channel: %#validator\n", chanl)

	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		// I don't think this is even possible but just to make sure
		channel.APILogger.Printf("[ERROR] id is missing in parameters for servers")
	} else {
		// this checks if there is Server ID being passed in the Channel being passed
		if chanl.Server != 0x0 {

			// This is to conversion form the String param
			idConversion, err := strconv.ParseUint(id, 10, 32)
			if err != nil {
				channel.APILogger.Printf("[ERROR] with id param String conversion into uint", err)
			}
			idConversionUint := uint(idConversion)
			// Checks if the param given for the servers id matches the servers id given with in the JSON Object
			if chanl.Server != idConversionUint {
				channel.APILogger.Printf("[ERROR] passed URI param id and the JSON ServerID do not match")

			} else {
				_, err := servers.GetServerByID(idConversionUint)
				if err != nil {
					channel.APILogger.Printf("[ERROR] fetching servers by ID returned and error: ", err)
				}
				err = data.AddChannel(*chanl)
				if err != nil {
					channel.APILogger.Printf("[ERROR] adding Channel to db: ", err)
				}
			}
		} else {
			idConversion, err := strconv.ParseUint(id, 10, 32)
			if err != nil {
				channel.APILogger.Printf("[ERROR] with id param String conversion into uint", err)
			}
			idConversionUint := uint(idConversion)
			_, err = servers.GetServerByID(idConversionUint)
			if err != nil {
				channel.APILogger.Printf("[ERROR] fetching servers by ID returned and error: ", err)
			}
			err = data.AddChannel(*chanl)
			if err != nil {
				channel.APILogger.Printf("[ERROR] adding Channel to db: ", err)
			}

		}

	}

}
