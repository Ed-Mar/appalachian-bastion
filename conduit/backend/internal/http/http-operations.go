package http

//Method Its justa  Collection of the HTTP methods, but why does when you plan on using the http package
// Its cause its not its only a String and I want to be able have a collection of just methods while the normal http package
// cause http.MethodGet is just a const not an interface/obj that is. Hope that makes sense for this silly boilerplate
type Method interface {
	GetHTTPMethod() string
}

///---------

type Post struct {
}

func (c Post) GetHTTPMethod() string {
	return "POST"
}

type Get struct {
}

func (r Get) GetHTTPMethod() string {
	return "GET"
}

type Put struct {
}

func (r Put) GetHTTPMethod() string {
	return "PUT"
}

type Delete struct {
}

func (r Delete) GetHTTPMethod() string {
	return "DELETE"
}

type Patch struct {
}

func (p Patch) GetHTTPMethod() string {
	return "PATCH"
}

type Head struct {
}

func (h Head) GetHTTPMethod() string {
	return "HEAD"
}

// TODO the rest of the weird ones
