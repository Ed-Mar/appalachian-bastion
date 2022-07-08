package users

import (
	"backend/internal"
	"backend/user-api/data"
	"net/http"
)

// swagger:route DELETE /servers/{id} servers deleteServer
// UpdateSingleton a servers details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  501: errorResponse

// Delete handles DELETE requests and removes items from the database
func (user *Users) Delete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	id := getUserID(r)

	user.userAPILogger.Println("[DEBUG] deleting record id", id)

	err := data.DeleteUser(uint(id))

	//NOT success response
	if err != nil {
		user.userAPILogger.Println("[ERROR] deleting record", err)

		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			user.userAPILogger.Println(err)
		}
		return
	}
	//success response
	rw.WriteHeader(http.StatusNoContent)
}
