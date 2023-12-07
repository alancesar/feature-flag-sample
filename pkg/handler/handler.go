package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"poc-growthbook/pkg/presenter"
)

type (
	FeatureFlagHandler interface {
		Eval(ctx context.Context, name string) bool
	}

	FeatureFlagRefresher interface {
		Refresh() error
	}
)

func Home(handler FeatureFlagHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := presenter.NewResponse(
			r.Header.Get("User-Agent"),
			r.Header.Get("Client-ID"),
		)

		if handler.Eval(r.Context(), "payload-v2") {
			payloadV2 := presenter.NewResponseV2(payload)
			bytes, _ := json.Marshal(&payloadV2)
			_, _ = w.Write(bytes)
			return
		}

		bytes, _ := json.Marshal(&payload)
		_, _ = w.Write(bytes)
	}
}

func Callback(refresher FeatureFlagRefresher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := refresher.Refresh(); err != nil {
			fmt.Println("failed to refresh features", err.Error())
		}
	}
}
