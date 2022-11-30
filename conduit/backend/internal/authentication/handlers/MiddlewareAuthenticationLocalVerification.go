package handlers

import (
	"net/http"
)

func (ha *ServiceHandler) MiddlewareAuthenticationLocalVerification(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		r, rw, ok := NewAuthViaLocalVerificationInContext(r, &rw, ha.ServiceLogger)
		if !ok {
			return
		}

		next.ServeHTTP(rw, r)
	})
}
