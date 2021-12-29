package handlers

import (
	"backend/internal"
	"backend/user-api/data"
	"context"
	"net/http"
)

// MiddlewareValidateUser validates the user in the request and calls next if ok
func (user *Users) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		laUser := &data.User{}

		err := internal.FromJSON(laUser, r.Body)
		if err != nil {
			user.userAPILogger.Println("[ERROR] JSON deserializing user", err)

			rw.WriteHeader(http.StatusBadRequest)
			err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
			if err != nil {
				user.userAPILogger.Println("[ERROR] encoding JSON: ", err)
			}
			return
		}

		// validate the user
		errs := user.validator.Validate(laUser)
		if len(errs) != 0 {
			user.userAPILogger.Println("[ERROR] validating user", errs)

			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			err := internal.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			if err != nil {
				user.userAPILogger.Println("[ERROR] encoding JSON: ", err)
			}
			return
		}

		// add the user to the context
		ctx := context.WithValue(r.Context(), KeyUser{}, laUser)
		r = r.WithContext(ctx)

		// Call the next handlers, which can be another middleware in the chain, or the final handlers.
		next.ServeHTTP(rw, r)
	})
}
