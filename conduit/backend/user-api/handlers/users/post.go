package users

import (
	"backend/user-api/data"
	"net/http"
)

// swagger:route POST /users users createUser
// Create a new servers
//
// responses:
//	200: ServerResponse
//  422: errorValidation
//  501: errorResponse

// Create handles POST requests to add new users
func (user *Users) Create(rw http.ResponseWriter, r *http.Request) {
	// fetch the servers from the context
	rw.Header().Add("Content-Type", "application/json")

	laUser := r.Context().Value(KeyUser{}).(*data.User)
	user.userAPILogger.Printf("[DEBUG] Inserting servers: %#validator\n", laUser)

	// Not Sure If I need to check if the Server Exist before I try to create the User
	err := data.AddUser(*laUser)
	if err != nil {
		user.userAPILogger.Println("[ERROR]: ", err)
	}
}
