package http_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"order-system/pkg/infra/config"
	client "order-system/pkg/infra/http"
)

func TestNewClient(t *testing.T) {
	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second

	c := client.NewClient(cfg, "http://example.com")
	if c == nil {
		t.Error("expected non-nil client")
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		headers    map[string]string
		wantErr    bool
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			body:       "success",
			headers:    map[string]string{"X-Test": "test"},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET method, got %s", r.Method)
				}
				for k, v := range tt.headers {
					if r.Header.Get(k) != v {
						t.Errorf("expected header %s=%s, got %s", k, v, r.Header.Get(k))
					}
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.body))
			}))
			defer server.Close()

			cfg := &config.Config{}
			cfg.HTTP.RequestTimeout = 30 * time.Second
			cfg.HTTP.MaxRequestSize = 1024

			c := client.NewClient(cfg, server.URL)

			ctx := context.Background()
			resp, err := c.Get(ctx, "/test", &client.RequestOption{
				Timeout:       5 * time.Second,
				RetryCount:    1,
				RetryInterval: time.Millisecond,
				MaxBodySize:   1024,
				Headers:       tt.headers,
			})

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if !tt.wantErr {
				if resp.StatusCode != tt.statusCode {
					t.Errorf("expected status code %d, got %d", tt.statusCode, resp.StatusCode)
				}
				if string(resp.Body) != tt.body {
					t.Errorf("expected body %s, got %s", tt.body, string(resp.Body))
				}
			}
		})
	}
}

func TestPost(t *testing.T) {
	tests := []struct {
		name       string
		reqBody    string
		statusCode int
		respBody   string
		wantErr    bool
	}{
		{
			name:       "success",
			reqBody:    "test body",
			statusCode: http.StatusCreated,
			respBody:   "created",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST method, got %s", r.Method)
				}
				body, _ := io.ReadAll(r.Body)
				if string(body) != tt.reqBody {
					t.Errorf("expected request body %s, got %s", tt.reqBody, string(body))
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.respBody))
			}))
			defer server.Close()

			cfg := &config.Config{}
			cfg.HTTP.RequestTimeout = 30 * time.Second
			cfg.HTTP.MaxRequestSize = 1024

			c := client.NewClient(cfg, server.URL)

			ctx := context.Background()
			resp, err := c.Post(ctx, "/test", []byte(tt.reqBody), nil)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if !tt.wantErr {
				if resp.StatusCode != tt.statusCode {
					t.Errorf("expected status code %d, got %d", tt.statusCode, resp.StatusCode)
				}
				if string(resp.Body) != tt.respBody {
					t.Errorf("expected body %s, got %s", tt.respBody, string(resp.Body))
				}
			}
		})
	}
}

func TestRequestTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 1 * time.Second
	cfg.HTTP.MaxRequestSize = 1024

	c := client.NewClient(cfg, server.URL)

	ctx := context.Background()
	_, err := c.Get(ctx, "/test", nil)
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}
