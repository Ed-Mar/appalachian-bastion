package meta

import (
	"backend/server-service/models"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func GetGETAllServersHTTPRequest() *http.Request {
	r, _ := http.NewRequest("GET", "http://localhost:9292/servers", nil)
	return r
}
func GetGETServerViaHTTPServerIDRequest(serverUUID string) *http.Request {

	r, _ := http.NewRequest("GET", "http://localhost:9292/servers/"+serverUUID, nil)
	return r
}

// GetPUTServerHTTPRequest returns the http.Request for Update Server when given the model.Server
func GetPUTServerHTTPRequest(server models.Server) *http.Request {
	la, err := json.Marshal(server)
	if err != nil {
		log.Println("updateServer error during marshaling of server obj.")
	}
	temp := strings.NewReader(string(la))
	r, _ := http.NewRequest("PUT", "http://localhost:9292/server/", temp)
	return r
}
func GetPOSTServerHTTPRequest(server models.Server) *http.Request {
	la, err := json.Marshal(server)
	if err != nil {
		log.Println("updateServer error during marshaling of server obj.")
	}
	temp := strings.NewReader(string(la))
	r, _ := http.NewRequest("POST", "http://localhost:9292/server/", temp)
	return r

}
func GetDeleteServerViaID(serverUUID string) *http.Request {
	r, _ := http.NewRequest("DELETE", "http://localhost:9292/servers/"+serverUUID, nil)
	return r

}
