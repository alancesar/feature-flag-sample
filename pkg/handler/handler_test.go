package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type (
	mockedFeatureFlagHandler struct {
		isOn bool
	}
)

func (h mockedFeatureFlagHandler) Eval(_ context.Context, _ string) bool {
	return h.isOn
}

func TestHome(t *testing.T) {
	type args struct {
		handler FeatureFlagHandler
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Should return payload v1",
			args: args{
				handler: mockedFeatureFlagHandler{
					isOn: false,
				},
			},
			want: []byte(`{"client_id":"some-client","user_agent":"httptest","version":1}`),
		},
		{
			name: "Should return payload v2",
			args: args{
				handler: mockedFeatureFlagHandler{
					isOn: true,
				},
			},
			want: []byte(`{"client_id":"some-client","user_agent":"httptest","version":2,"message":"looks like the feature flag is on"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("User-Agent", "httptest")
			req.Header.Set("Client-ID", "some-client")
			w := httptest.NewRecorder()

			Home(tt.args.handler).ServeHTTP(w, req)
			res := w.Result()
			defer func() {
				_ = res.Body.Close()
			}()

			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			if !reflect.DeepEqual(data, tt.want) {
				t.Errorf("expected %s got %s", tt.want, data)
			}
		})
	}
}
