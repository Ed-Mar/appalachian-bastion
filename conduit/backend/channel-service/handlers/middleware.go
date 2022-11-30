package handlers

import (
	"backend/channel-service/models"
	"backend/internal"
	"context"
	"net/http"
)

// MiddlewareValidateChannel validates the channel in the request and calls next if ok
func (channel *Channels) MiddlewareValidateChannel(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		chanl := &models.Channel{}

		err := internal.FromJSON(chanl, r.Body)
		if err != nil {
			channel.APILogger.Println("[ERROR] JSON deserializing channel", err)

			rw.WriteHeader(http.StatusBadRequest)
			err := internal.ToJSON(&GenericError{Message: err.Error()}, rw)
			if err != nil {
				channel.APILogger.Println("[ERROR] encoding JSON: ", err)
			}
			return
		}

		// validate the channel
		errs := channel.validator.Validate(chanl)
		if len(errs) != 0 {
			channel.APILogger.Println("[ERROR] validating channel", errs)

			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			err := internal.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			if err != nil {
				channel.APILogger.Println("[ERROR] encoding JSON: ", err)
			}
			return
		}

		// add the channel to the context
		ctx := context.WithValue(r.Context(), KeyChannel{}, chanl)
		r = r.WithContext(ctx)

		// Call the next gerneric-handlers, which can be another middleware in the chain, or the final gerneric-handlers.
		next.ServeHTTP(rw, r)
	})
}
