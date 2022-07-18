package event_handler

import (
	"backend/event-streaming/kafka/messages/model"
	"io"
	"io/ioutil"
	"net/http"
)

const restServerURL = "http://localhost:9292/servers"

var ErrREST = "[ERROR] [EVENT] [SERVER] [HTTP]:"
var ErrHTTPResponse = "[ERROR] [EVENT] [SERVER] [HTTP] [RESPONSE]:"
var ErrEventMessageContent = "[ERROR] [EVENT] [SERVER] [MESSAGE-CONTENT]:"

func (serverEvent *ServerEvents) EventMux(message model.EventMessage) string {
	switch message.ServiceTargetName {
	case "Get - Singleton":
		serverEvent.Logger.Println("Get - Singleton | FOUND ")
		serverID := message.SagaTransactionData["serverid"]
		var getSingletonURL = restServerURL + "/" + serverID
		serverEvent.Logger.Println("Send HTTP Request to: ", getSingletonURL)
		response, httpErr := http.Get(getSingletonURL)
		if httpErr != nil {
			serverEvent.Logger.Println(ErrREST, httpErr)
			return "shit"
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				serverEvent.Logger.Println(ErrHTTPResponse, "Error Reading Body of respond | ", err)
			}
		}(response.Body)

		body, bodyERR := ioutil.ReadAll(response.Body)
		if bodyERR != nil {
			serverEvent.Logger.Println(ErrREST, bodyERR)
			return "shit"
		}
		serverEvent.Logger.Println(string(body))
		return "LOOK ABOVE"
	case "Get - Collection":
		serverEvent.Logger.Println("Get - Collection | FOUND ")

		serverEvent.Logger.Println("Send HTTP Request to: ", restServerURL)
		response, httpErr := http.Get(restServerURL)
		if httpErr != nil {
			serverEvent.Logger.Println(ErrREST, httpErr)
			return "logs"
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				serverEvent.Logger.Println(ErrHTTPResponse, "Error Reading Body of respond | ", err)
			}
		}(response.Body)

		body, bodyERR := ioutil.ReadAll(response.Body)
		if bodyERR != nil {
			serverEvent.Logger.Println(ErrREST, bodyERR)
			return "shit"
		}
		serverEvent.Logger.Println(string(body))
		return "Look in logs"
	case "":
		serverEvent.Logger.Println(ErrEventMessageContent, "No Target Operation passed |  ", message.ServiceTargetName)
		return "shit"
	default:
		serverEvent.Logger.Println(ErrEventMessageContent, "Non-supported Operation Target Send | ", message.ServiceTargetName)
		return "shit"

	}

}
