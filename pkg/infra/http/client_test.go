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
		t.Error("Expected non-nil client")
	}
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
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
		Headers: map[string]string{
			"X-Test": "test",
		},
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if string(resp.Body) != "success" {
		t.Errorf("Expected body 'success', got '%s'", string(resp.Body))
	}
}

func TestPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != "test body" {
			t.Errorf("Expected body 'test body', got '%s'", string(body))
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second
	cfg.HTTP.MaxRequestSize = 1024

	c := client.NewClient(cfg, server.URL)

	ctx := context.Background()
	resp, err := c.Post(ctx, "/test", []byte("test body"), nil)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}
}

func TestPut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != "update" {
			t.Errorf("Expected body 'update', got '%s'", string(body))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second
	cfg.HTTP.MaxRequestSize = 1024

	c := client.NewClient(cfg, server.URL)

	ctx := context.Background()
	resp, err := c.Put(ctx, "/test", []byte("update"), nil)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second
	cfg.HTTP.MaxRequestSize = 1024

	c := client.NewClient(cfg, server.URL)

	ctx := context.Background()
	resp, err := c.Delete(ctx, "/test", nil)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}
