package handlers

import (
	"backend/internal/authentication/handlers"
	"net/http"
)

// ValidateHeaderTestDev validates the user in the request and calls next if ok
func (DevObj *DevTest) ValidateHeaderTestDev(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		rw.Header().Set("Access-Control-Allow-Methods", "POST")
		//bear := r.Header.Get("Authentication")
		//DevObj.APILogger.Printf("AccessToken: %#v", bear)

		auth, err := handlers.newAuthViaLocalVerification(rw, r)
		if err != nil {
			DevObj.GenericHandler.ServiceLogger.Println(err)
			return
		}

		r = r.WithContext(handlers.SetAuthInContext(r.Context(), auth))

		next.ServeHTTP(rw, r)
	})
}

//
