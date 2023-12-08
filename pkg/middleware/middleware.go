package middleware

import (
	"net/http"
	"poc-growthbook/pkg/tracing"
)

type (
	FeatureFlagHandler interface {
		GetFeatureMap() ([]byte, error)
	}
)

func InjectUserData(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID := r.Header.Get("Client-ID")
		r = r.WithContext(tracing.SetClientIDInContext(r.Context(), clientID))
		next.ServeHTTP(w, r)
	})
}
