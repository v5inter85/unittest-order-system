package http_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	httpclient "order-system/pkg/infra/http"
	"order-system/pkg/infra/config"
)

func TestNewClient(t *testing.T) {
	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second

	client := httpclient.NewClient(cfg, "http://example.com")
	require.NotNil(t, client)
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "test", r.Header.Get("X-Test"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second
	cfg.HTTP.MaxRequestSize = 1024

	client := httpclient.NewClient(cfg, server.URL)

	ctx := context.Background()
	resp, err := client.Get(ctx, "/test", &httpclient.RequestOption{
		Timeout:     5 * time.Second,
		RetryCount:  1,
		MaxBodySize: 1024,
		Headers: map[string]string{
			"X-Test": "test",
		},
	})

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, []byte("success"), resp.Body)
}

func TestPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, []byte("test body"), body)
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second
	cfg.HTTP.MaxRequestSize = 1024

	client := httpclient.NewClient(cfg, server.URL)

	ctx := context.Background()
	resp, err := client.Post(ctx, "/test", []byte("test body"), nil)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestPut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, []byte("update"), body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second
	cfg.HTTP.MaxRequestSize = 1024

	client := httpclient.NewClient(cfg, server.URL)

	ctx := context.Background()
	resp, err := client.Put(ctx, "/test", []byte("update"), nil)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second
	cfg.HTTP.MaxRequestSize = 1024

	client := httpclient.NewClient(cfg, server.URL)

	ctx := context.Background()
	resp, err := client.Delete(ctx, "/test", nil)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestMaxBodySize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("large response body exceeding limit"))
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second
	cfg.HTTP.MaxRequestSize = 1024

	client := httpclient.NewClient(cfg, server.URL)

	ctx := context.Background()
	_, err := client.Get(ctx, "/test", &httpclient.RequestOption{
		MaxBodySize: 5,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "response body too large")
}

func TestContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 30 * time.Second
	cfg.HTTP.MaxRequestSize = 1024

	client := httpclient.NewClient(cfg, server.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.Get(ctx, "/test", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), context.DeadlineExceeded.Error())
}

func TestDoRequestErrorHandling(t *testing.T) {
	cfg := &config.Config{}
	cfg.HTTP.RequestTimeout = 1 * time.Second
	cfg.HTTP.MaxRequestSize = 1024

	client := httpclient.NewClient(cfg, "http://127.0.0.1:12345") // Invalid URL

	ctx := context.Background()
	_, err := client.Get(ctx, "/test", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "request failed")
}
