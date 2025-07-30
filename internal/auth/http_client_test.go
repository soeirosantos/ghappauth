package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewHTTPClient(t *testing.T) {
	client := NewHTTPClient(nil)
	if client == nil {
		t.Fatal("NewHTTPClient() returned nil")
	}
	if client.config.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", client.config.Timeout)
	}
	if client.config.MaxRetries != 3 {
		t.Errorf("Expected max retries 3, got %d", client.config.MaxRetries)
	}

	customConfig := &HTTPClientConfig{
		Timeout:    10 * time.Second,
		MaxRetries: 5,
		UserAgent:  "test-agent",
	}
	client = NewHTTPClient(customConfig)
	if client.config.Timeout != 10*time.Second {
		t.Errorf("Expected timeout 10s, got %v", client.config.Timeout)
	}
	if client.config.MaxRetries != 5 {
		t.Errorf("Expected max retries 5, got %d", client.config.MaxRetries)
	}
	if client.config.UserAgent != "test-agent" {
		t.Errorf("Expected user agent 'test-agent', got %s", client.config.UserAgent)
	}
}

func TestHTTPClient_doRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Accept") != "application/vnd.github.v3+json" {
			t.Errorf("Expected Accept header 'application/vnd.github.v3+json', got %s", r.Header.Get("Accept"))
		}
		if r.Header.Get("User-Agent") != "ghappauth/1.0" {
			t.Errorf("Expected User-Agent header 'ghappauth/1.0', got %s", r.Header.Get("User-Agent"))
		}
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := NewHTTPClient(nil)
	resp, err := client.doRequest(context.Background(), &RequestConfig{
		Method:    "GET",
		URL:       server.URL,
		AuthToken: "test-token",
	})

	if err != nil {
		t.Fatalf("doRequest() error = %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHTTPClient_doRequest_RetryOn5xx(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "success"}`))
		}
	}))
	defer server.Close()

	client := NewHTTPClient(&HTTPClientConfig{
		MaxRetries: 3,
		RetryDelay: 10 * time.Millisecond,
	})

	resp, err := client.doRequest(context.Background(), &RequestConfig{
		Method: "GET",
		URL:    server.URL,
	})

	if err != nil {
		t.Fatalf("doRequest() error = %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestHTTPClient_doRequest_RetryOn429(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.WriteHeader(http.StatusTooManyRequests)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "success"}`))
		}
	}))
	defer server.Close()

	client := NewHTTPClient(&HTTPClientConfig{
		MaxRetries: 3,
		RetryDelay: 10 * time.Millisecond,
	})

	resp, err := client.doRequest(context.Background(), &RequestConfig{
		Method: "GET",
		URL:    server.URL,
	})

	if err != nil {
		t.Fatalf("doRequest() error = %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if attempts != 2 {
		t.Errorf("Expected 2 attempts, got %d", attempts)
	}
}

func TestHTTPClient_doRequest_NoRetryOn4xx(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "bad request"}`))
	}))
	defer server.Close()

	client := NewHTTPClient(&HTTPClientConfig{
		MaxRetries: 3,
		RetryDelay: 10 * time.Millisecond,
	})

	resp, err := client.doRequest(context.Background(), &RequestConfig{
		Method: "GET",
		URL:    server.URL,
	})

	if err != nil {
		t.Fatalf("doRequest() error = %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

func TestHTTPClient_DoRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "test-app", "id": 123}`))
	}))
	defer server.Close()

	client := NewHTTPClient(nil)
	
	var result struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	}
	
	err := client.DoRequest(context.Background(), &RequestConfig{
		Method:        "GET",
		URL:           server.URL,
		ExpectedStatus: http.StatusOK,
	}, &result)

	if err != nil {
		t.Fatalf("DoRequest() error = %v", err)
	}
	if result.Name != "test-app" {
		t.Errorf("Expected name 'test-app', got %s", result.Name)
	}
	if result.ID != 123 {
		t.Errorf("Expected ID 123, got %d", result.ID)
	}
}

func TestHTTPClient_DoRequest_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Not Found", "documentation_url": "https://docs.github.com"}`))
	}))
	defer server.Close()

	client := NewHTTPClient(nil)
	
	var result struct {
		Name string `json:"name"`
	}
	
	err := client.DoRequest(context.Background(), &RequestConfig{
		Method:        "GET",
		URL:           server.URL,
		ExpectedStatus: http.StatusOK,
	}, &result)

	if err == nil {
		t.Fatal("DoRequest() should have returned an error")
	}
	if err.Error() != "GitHub API error: Not Found (status: 404)" {
		t.Errorf("Expected error 'GitHub API error: Not Found (status: 404)', got %v", err)
	}
}

func TestShouldRetry(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expected   bool
	}{
		{"200 OK", 200, false},
		{"400 Bad Request", 400, false},
		{"404 Not Found", 404, false},
		{"408 Request Timeout", 408, true},
		{"429 Too Many Requests", 429, true},
		{"500 Internal Server Error", 500, true},
		{"502 Bad Gateway", 502, true},
		{"503 Service Unavailable", 503, true},
		{"504 Gateway Timeout", 504, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldRetry(tt.statusCode)
			if result != tt.expected {
				t.Errorf("shouldRetry(%d) = %v, want %v", tt.statusCode, result, tt.expected)
			}
		})
	}
}

 