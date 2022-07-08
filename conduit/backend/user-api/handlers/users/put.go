package users

import (
	"backend/internal"
	"backend/user-api/data"
	"net/http"
)

// swagger:route PUT /servers servers updateServer
// UpdateSingleton a servers details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// Update handles PUT requests to update servers
func (user *Users) Update(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	// fetch the user from the context
	laUser := r.Context().Value(KeyUser{}).(*data.User)
	user.userAPILogger.Println("[DEBUG] updating record id", laUser.ID)

	err := data.UpdateUser(*laUser)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			user.userAPILogger.Println(err)
		}
		return
	}

	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)
}
