package handlers

import (
	"backend/internal"
	"backend/server-service/models"
	"context"
	"net/http"
)

// MiddlewareValidateServer validates the servers in the request and calls next if ok
func (server *ServerHandler) MiddlewareValidateServer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		serv := &models.Server{}
		err := internal.FromJSON(serv, r.Body)
		if err != nil {
			server.APILogger.Println("[ERROR] JSON deserializing servers", err)

			rw.WriteHeader(http.StatusBadRequest)
			err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
			if err != nil {
				server.APILogger.Println("[ERROR] [JSON] encoding JSON: ", err)
			}
			return
		}

		// validate the servers
		errs := server.validator.Validate(serv)
		if len(errs) != 0 {
			server.APILogger.Println("[ERROR] validating servers", errs)

			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			err := internal.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			if err != nil {
				server.APILogger.Println("[ERROR] encoding JSON: ", err)
			}
			return
		}

		// add the servers to the context
		ctx := context.WithValue(r.Context(), KeyServer{}, serv)
		r = r.WithContext(ctx)

		// Call the next gerneric-handlers, which can be another middleware in the chain, or the final gerneric-handlers.
		next.ServeHTTP(rw, r)
	})
}
