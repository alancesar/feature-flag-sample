package middleware

import (
	"context"
	"net/http"
)

type (
	FeatureFlagHandler interface {
		GetFeatureMap() ([]byte, error)
	}
)

func InjectUserData(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID := r.Header.Get("Client-ID")
		ctx := context.WithValue(r.Context(), "client-id", clientID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
