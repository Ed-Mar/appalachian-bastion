package meta

import (
	"backend/internal/http"
	"backend/internal/services"
)

//TODO Comeback and get this from a config-file
const (
	serverServiceName = "server-service"
	port              = "9292"
	host              = "localhost"
)

type serverServiceMeta struct {
}

func (s serverServiceMeta) GetName() string {

	return serverServiceName
}
func (s serverServiceMeta) GetHost() string {
	return host
}
func (s serverServiceMeta) GetPort() string {
	return port
}

// GetHostAndPort TODO Comeback and deal with when we containerize this
func (s serverServiceMeta) GetHostAndPort() string {
	return host + ":" + port
}

// GetFullURL TODO Comeback and deal with when we containerize this
func (s serverServiceMeta) GetFullURL() string {
	return "http//:" + host + ":" + port
}

func (s serverServiceMeta) GetServiceHttpRoutes() (routes map[string]services.Route) {
	routes[getAllServers.GetHttpRouteName()] = *getAllServers
	routes[getServerViaServerID1.GetHttpRouteName()] = *getServerViaServerID1
	routes[updateServer1.GetHttpRouteName()] = *updateServer1
	routes[postServer.GetHttpRouteName()] = *postServer
	routes[deleteServerViaServerID.GetHttpRouteName()] = *deleteServerViaServerID

	return routes
}

var (
	getAllServers           = services.Route{}.NewHttpRoute("getAllServers", http.Get{}, "/servers")
	getServerViaServerID1   = services.Route{}.NewHttpRoute("getServerViaServerID", http.Get{}, "/servers/{serverID}")
	updateServer1           = services.Route{}.NewHttpRoute("updateServer", http.Put{}, "/server")
	postServer              = services.Route{}.NewHttpRoute("postServer", http.Post{}, "/server")
	deleteServerViaServerID = services.Route{}.NewHttpRoute("deleteServerViaServerID", http.Delete{}, "/servers/{serverID}")
)
