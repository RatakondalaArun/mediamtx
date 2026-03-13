package extras

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bluenviron/mediamtx/internal/conf"
)

func TestLoadConfigFromAPI(t *testing.T) {
	tests := []struct {
		name    string
		apiURL  string
		handler func(w http.ResponseWriter, r *http.Request)
		want    *conf.Conf
		wantErr bool
	}{
		{
			name:    "empty URL",
			apiURL:  "",
			want:    nil,
			wantErr: true,
		},
		{
			name:   "successful config load",
			apiURL: "http://example.com/config",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(&conf.Conf{})
			},
			want:    &conf.Conf{},
			wantErr: false,
		},
		{
			name:   "API error status",
			apiURL: "http://example.com/config",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "invalid JSON response",
			apiURL: "http://example.com/config",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("invalid json"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "load from localhost:7777/mediamtx.config.json",
			apiURL: "http://localhost:7777/mediamtx.config.json",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				config := &conf.Conf{
					Paths: make(map[string]*conf.Path),
				}
				json.NewEncoder(w).Encode(config)
			},
			want:    &conf.Conf{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var server *httptest.Server
			if tt.handler != nil {
				server = httptest.NewServer(http.HandlerFunc(tt.handler))
				defer server.Close()
				tt.apiURL = server.URL
			}

			got, err := LoadConfigFromAPI(tt.apiURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfigFromAPI() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Errorf("LoadConfigFromAPI() got nil, want non-nil")
			}
		})
	}
}
