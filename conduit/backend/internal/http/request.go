package http

import (
	"net/http"
)

func GetInternalDomainHTTPRequest(i interface{}) (request *http.Request) {

	r, _ := http.NewRequest("GET", "http://localhost:9292/servers/", nil)

	return r
}
