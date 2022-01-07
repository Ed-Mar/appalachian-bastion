package users

import (
	"backend/internal"
	"backend/user-api/data"
	"net/http"
)

// swagger:route GET /users user listUsers
// Return a list of users from the database
// responses:
//	200: ServersResponse

// ListAll handles GET requests and returns all current servers
func (user *Users) ListAll(rw http.ResponseWriter, r *http.Request) {
	user.userAPILogger.Println("[DEBUG] get all records")
	rw.Header().Add("Content-Type", "application/json")

	laUser, err := data.GetUsers()
	if err != nil {
		user.userAPILogger.Println("[ERROR]: ", err)
	}

	err = internal.ToJSON(laUser, rw)
	if err != nil {
		// we should never be here but log the error just encase
		user.userAPILogger.Println("[ERROR] serializing servers", err)
	}
}

// swagger:route GET /users/{id} servers listSingleUser
// Return a user from the database with the matching provided id
// responses:
//	200: ServerResponse
//	404: errorResponse

// ListSingle handles GET requests
func (user *Users) ListSingle(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	id := getUserID(r)

	user.userAPILogger.Println("[DEBUG] get record id", id)

	serv, err := data.GetUserByID(uint(id))

	switch err {
	case nil:
	default:
		user.userAPILogger.Println("[ERROR] fetching servers", err)

		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			user.userAPILogger.Println("[ERROR] in JSON encoding: ", err)
		}
		return
	}

	err = internal.ToJSON(serv, rw)
	if err != nil {
		// we should never be here but log the error just encase
		user.userAPILogger.Println("[ERROR] serializing servers", err)
	}
}
