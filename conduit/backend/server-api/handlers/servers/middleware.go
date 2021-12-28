package servers

import (
	"backend/internal"
	"backend/server-api/data"
	"context"
	"net/http"
)

// MiddlewareValidateServer validates the server in the request and calls next if ok
func (server *Servers) MiddlewareValidateServer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		serv := &data.Server{}

		err := internal.FromJSON(serv, r.Body)
		if err != nil {
			server.severAPILogger.Println("[ERROR] JSON deserializing server", err)

			rw.WriteHeader(http.StatusBadRequest)
			err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
			if err != nil {
				server.severAPILogger.Println("[ERROR] encoding JSON: ", err)
			}
			return
		}

		// validate the server
		errs := server.validator.Validate(serv)
		if len(errs) != 0 {
			server.severAPILogger.Println("[ERROR] validating server", errs)

			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			err := internal.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			if err != nil {
				server.severAPILogger.Println("[ERROR] encoding JSON: ", err)
			}
			return
		}

		// add the server to the context
		ctx := context.WithValue(r.Context(), KeyServer{}, serv)
		r = r.WithContext(ctx)

		// Call the next handlers, which can be another middleware in the chain, or the final handlers.
		next.ServeHTTP(rw, r)
	})
}
