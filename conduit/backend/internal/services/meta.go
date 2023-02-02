package services

import (
	"backend/internal/http"
	"net/url"
)

type ServiceMeta interface {
	GetName() string
	GetHost() string
	GetPort() string
	GetHostAndPort() string
	GetFullURL() url.URL
	// GetServiceHttpRoutes Do not know if the map is the best with this, but meh
	GetServiceHttpRoutes() map[string]Route
}

// Route I still think this is little silly, and that there is a better way of doing this. But I need to move on.
type Route struct {
	name string
	//The http Method UserType of this request in the route
	method http.Method
	// should like the last part of th API example:"/servers/{serverID}/channels"
	uRL string
}

// Should I add this to the ServiceMeta interfacemaybe...
func (h Route) NewHttpRoute(name string, method http.Method, uRL string) *Route {
	return &Route{
		name:   name,
		method: method,
		uRL:    uRL,
	}
}
func (h Route) GetHttpRouteName() string {
	return h.name
}
func (h Route) GetHttpRouteMethod() http.Method {
	return h.method
}
func (h Route) GetHttpRouteURL() string {
	return h.uRL
}
